package executor

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/letscool/lc-go/dij"
	. "github.com/letscool/lc-go/lg"
	"golang.org/x/net/netutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"
)

const (
	DefaultWebServerPort = 8000
	HttpTagName          = "http"
	WebConfigKey         = "webserver.config"
)

type WebContext struct {
	*gin.Context
}

type HandlerFunc func(*WebContext)

type WebServerSpec interface {
	WebControllerSpec
	iAmAWebServer()
}

type WebServer struct {
	WebController
}

func (w *WebServer) iAmAWebServer() {

}

type WebControllerSpec interface {
	iAmAWebController()
	SetupRouter(router gin.IRouter)
}

type WebController struct {
}

func (w *WebController) iAmAWebController() {

}

func (w *WebController) SetupRouter(router gin.IRouter) {
	fmt.Printf("Should setup router?\n")
}

type WebMiddlewareSpec interface {
	iAmAWebMiddleware()
	GetHandlers() []HandlerFunc
}

type WebMiddleware struct {
}

func (m *WebMiddleware) iAmAWebMiddleware() {

}

func (m *WebMiddleware) GetHandlers() []HandlerFunc {
	fmt.Printf("Should provide handlers?\n")
	return nil
}

//func IsTypeOfWebServer(typ reflect.Type) bool {
//	wssTyp := reflect.TypeOf((*WebServerSpec)(nil)).Elem()
//	fmt.Printf("wssTyp: %v\n", wssTyp)
//	return typ.Implements(wssTyp)
//}

func IsTypeOfWebServer(typ reflect.Type) bool {
	if typ.Kind() == reflect.Pointer {
		if elemTyp := typ.Elem(); elemTyp.Kind() == reflect.Struct {
			return IsTypeOfWebServer(elemTyp)
		}
		return false
	}
	instPtrValue := reflect.New(typ)
	instIf := instPtrValue.Interface()
	_, ok := instIf.(WebServerSpec)
	return ok
}

func IsTypeOfWebController(typ reflect.Type) bool {
	if typ.Kind() == reflect.Pointer {
		if elemTyp := typ.Elem(); elemTyp.Kind() == reflect.Struct {
			return IsTypeOfWebController(elemTyp)
		}
		return false
	}
	instPtrValue := reflect.New(typ)
	instIf := instPtrValue.Interface()
	_, ok := instIf.(WebControllerSpec)
	return ok
}

func IsTypeOfWebMiddleware(typ reflect.Type) bool {
	if typ.Kind() == reflect.Pointer {
		if elemTyp := typ.Elem(); elemTyp.Kind() == reflect.Struct {
			return IsTypeOfWebMiddleware(elemTyp)
		}
		return false
	}
	instPtrValue := reflect.New(typ)
	instIf := instPtrValue.Interface()
	_, ok := instIf.(WebMiddlewareSpec)
	return ok
}

type WebServerConfig struct {
	address string
	port    int
	maxConn int
}

func CreateServeMux(servInstPtr any) *http.ServeMux {
	servPtrValue := reflect.ValueOf(servInstPtr)
	servPtrType := servPtrValue.Type()
	servValue := servPtrValue.Elem()
	servType := servValue.Type()

	log.Printf("Server type: %v, ptr type: %v", servType, servPtrType)
	mux := http.NewServeMux()
	//mux.HandleFunc("/", ws.HandleRequest)

	for j := 0; j < servType.NumField(); j++ {
		fieldSpec := servType.Field(j)
		_, existingHttpTag := fieldSpec.Tag.Lookup(HttpTagName)
		if existingHttpTag {
			// f := func(writer http.ResponseWriter, r *http.Request) {}
		}
	}

	return mux
}

func LaunchWebServer(webServerType reflect.Type, others ...any) {
	ref := dij.DependencyReference{}
	config := WebServerConfig{}
	for _, other := range others {
		otherTyp := reflect.TypeOf(other)
		switch v := other.(type) {
		case WebServerConfig:
			config = v
			ref[WebConfigKey] = v
		default:
			log.Println("No ideal about type:", otherTyp)
		}
	}
	instPtr, err := dij.CreateInstance(webServerType, &ref, "")
	if err != nil {
		log.Panic(err)
	}

	addr := fmt.Sprintf("%v:%d", config.address, Ife(config.port <= 0, DefaultWebServerPort, config.port))
	router := CreateServeMux(instPtr)

	srv := http.Server{
		ReadHeaderTimeout: time.Second * 5,
		ReadTimeout:       time.Second * 10,
		Handler:           router,
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	if config.maxConn > 0 {
		listener = netutil.LimitListener(listener, config.maxConn)
		//log.Printf("max connections set to %d\n", ws.maxConn)
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			//
		}
	}(listener)

	log.Printf("listening on %s\n", listener.Addr().String())

	go func() {
		if err := srv.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	<-signalChannel

	log.Printf("interrupted, shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v\n", err)
	}
}

func PrepareGin(webServerType reflect.Type, others ...any) (*gin.Engine, dij.DependencyReferencePtr, error) {
	ref := dij.DependencyReference{}
	config := WebServerConfig{}
	for _, other := range others {
		otherTyp := reflect.TypeOf(other)
		switch v := other.(type) {
		case WebServerConfig:
			config = v
		default:
			log.Println("No ideal about type:", otherTyp)
		}
	}
	ref[WebConfigKey] = config

	if !IsTypeOfWebServer(webServerType) {
		return nil, nil, fmt.Errorf("the type(%v) is not a web server", webServerType)
	}

	instPtr, err := dij.CreateInstance(webServerType, &ref, "^")
	if err != nil {
		log.Panic(err)
	}

	router := gin.Default()

	if err := setupRouter(instPtr, webServerType, router, &ref); err != nil {
		return nil, nil, err
	}

	return router, &ref, nil
}

func LaunchGin(webServerType reflect.Type, others ...any) error {
	engine, refPtr, err := PrepareGin(webServerType, others)
	if err != nil {
		return err
	}
	v, _ := refPtr.Get(WebConfigKey)
	config := v.(WebServerConfig)

	addr := fmt.Sprintf("%v:%d", config.address, Ife(config.port <= 0, DefaultWebServerPort, config.port))
	return engine.Run(addr)
}

func setupRouter(instPtr any, instType reflect.Type, router gin.IRouter, refPtr dij.DependencyReferencePtr) error {
	predecessor := make([]int, 0)
	plugins := make([]int, 0)
	extenders := make([]int, 0)
	for i := 0; i < instType.NumField(); i++ {
		field := instType.Field(i)
		fieldTyp := field.Type
		if IsTypeOfWebController(fieldTyp) {
			if field.Anonymous {
				predecessor = append(predecessor, i)
			} else {
				extenders = append(extenders, i)
			}
		} else if IsTypeOfWebMiddleware(fieldTyp) {
			plugins = append(plugins, i)
		}
	}

	// setup current router
	if len(predecessor) != 1 {
		return fmt.Errorf("struct '%s' should embeded web controller or web server.(%d)", instType.Name(), len(predecessor))
	} else {
		field := instType.Field(predecessor[0])
		tag, exists := field.Tag.Lookup(HttpTagName)
		if exists {
			attrs := ParseStructTag(tag)
			nameAttr, existingName := attrs.FirstAttrWithValOnly()
			if existingName {
				router = router.Group(nameAttr.Val)
			}
		}
		ctrl := instPtr.(WebControllerSpec)
		ctrl.SetupRouter(router)
	}

	// setup middlewares
	if len(plugins) > 0 {
		for _, idx := range plugins {
			field := instType.Field(idx)
			fieldTyp := field.Type
			fmt.Printf("Middleware: %v from %v\n", fieldTyp, instType)
			if fieldTyp.Kind() != reflect.Pointer || fieldTyp.Elem().Kind() != reflect.Struct {
				return fmt.Errorf("middleware's type(%v) should be a kind of struct point", fieldTyp)
			}
			instValue := reflect.ValueOf(instPtr).Elem()
			var fieldIf any
			var exists bool
			fieldIf, exists = refPtr.GetForDiField(instType, idx)
			if !exists {
				fieldValue := instValue.Field(idx)
				if fieldValue.IsZero() {
					return fmt.Errorf("field(%s) should not be zero", fieldTyp.Name())
				}
				if field.IsExported() {
					fieldIf = fieldValue.Interface()
				} else {
					fieldIf = reflect.NewAt(fieldTyp, fieldValue.Addr().UnsafePointer()).Elem().Interface()
				}
			} else {
				//fmt.Printf("middleware load from dij: %v\n", fieldTyp)
			}
			handlers := fieldIf.(WebMiddlewareSpec).GetHandlers()
			if len(handlers) > 0 {
				for _, handler := range handlers {
					router.Use(func(c *gin.Context) {
						// TODO: refine this
						handler(&WebContext{c})
					})
				}
			}
		}
	}

	// setup extenders
	if len(extenders) > 0 {
		for _, idx := range extenders {
			field := instType.Field(idx)
			fieldTyp := field.Type
			fmt.Printf("Extender: %v from %v\n", fieldTyp, instType)
			if fieldTyp.Kind() != reflect.Pointer || fieldTyp.Elem().Kind() != reflect.Struct {
				return fmt.Errorf("appending controller's type(%v) should be a kind of struct point", fieldTyp)
			}
			instValue := reflect.ValueOf(instPtr).Elem()
			var fieldIf any
			var exists bool
			fieldIf, exists = refPtr.GetForDiField(instType, idx)
			if !exists {
				fieldValue := instValue.Field(idx)
				if fieldValue.IsZero() {
					return fmt.Errorf("field(%s) should not be zero", fieldTyp.Name())
				}
				if field.IsExported() {
					fieldIf = fieldValue.Interface()
				} else {
					fieldIf = reflect.NewAt(fieldTyp, fieldValue.Addr().UnsafePointer()).Elem().Interface()
				}
			} else {
				//fmt.Printf("extenders load from dij: %v\n", fieldTyp)
			}
			if err := setupRouter(fieldIf, fieldTyp.Elem(), router, refPtr); err != nil {
				return err
			}
		}
	}

	// all done
	return nil
}

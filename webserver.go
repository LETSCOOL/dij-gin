package dij_gin

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/letscool/lc-go/dij"
	. "github.com/letscool/lc-go/lg"
	"log"
	"reflect"
	"regexp"
	"strings"
)

const (
	DefaultWebServerPort = 8000
	HttpTagName          = "http"
	WebConfigKey         = "webserver.config"
)

type WebContextSpec interface {
	iAmAWebContext()
}

type WebContext struct {
	*gin.Context
}

func (c *WebContext) iAmAWebContext() {

}

//type HandlerFunc func(*WebContext)

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
	//SetupRouter(router gin.IRouter, others ...any)
}

type WebController struct {
}

func (w *WebController) iAmAWebController() {

}

type WebCtlParamDef struct {
	index         int
	fieldSpec     reflect.StructField
	existsTag     bool           // exists http tag
	attrs         StructTagAttrs // come from http tag
	preferredName string
}

func (c *WebCtlParamDef) preferredText(key string, allowedValOnly bool, allowedFieldName bool) string {
	if c.existsTag {
		if allowedValOnly {
			if attr, ok := c.attrs.FirstAttrWithValOnly(); ok {
				if len(attr.Val) > 0 {
					return attr.Val
				}
			}
		}
		if attr, ok := c.attrs.FirstAttrsWithKey(key); ok {
			if len(attr.Val) > 0 {
				return attr.Val
			}
		}
	}
	if allowedFieldName {
		return c.fieldSpec.Name
	}
	return ""
}

//func (w *WebController) SetupRouter(router gin.IRouter, others ...any) {
//fmt.Printf("Should setup router?\n")
//setupHandlers(router, others[0])
//}

type handlerWrapper struct {
	method  string
	path    string
	handler func(*gin.Context)
}

func setupHandlers(routes gin.IRoutes, instPtr any) {
	wrappers := createHandlers(instPtr, true)
	for _, w := range wrappers {
		routes.Handle(w.method, w.path, w.handler)
	}
}

func createHandlers(instPtr any, forReq bool) []handlerWrapper {
	wrappers := make([]handlerWrapper, 0)
	webCtxType := reflect.TypeOf(WebContext{})
	instPtrType := reflect.TypeOf(instPtr)
	handleMethodRegex := regexp.MustCompile(Ife(forReq, `^(get|post|put|patch|delete|head|connect|options|trace)`, `^(handle)`))
	// TODO: how to deal routing for static pages
	for i := 0; i < instPtrType.NumMethod(); i++ {
		method := instPtrType.Method(i)
		if method.IsExported() {
			methodType := method.Type
			if methodType.NumIn() == 2 && methodType.NumOut() <= 1 {
				param1Typ := methodType.In(1)
				if IsTypeOfWebContext(param1Typ) && param1Typ.Kind() == reflect.Struct {
					methodName := method.Name
					lowerMethodName := strings.ToLower(methodName)
					reqMethod := string(handleMethodRegex.Find([]byte(lowerMethodName)))
					reqPath := lowerMethodName[len(reqMethod):]
					param1Defs := make([]WebCtlParamDef, 0)
					if param1Typ != webCtxType {
						fmt.Printf("[*%v]'s method %d: func %v(%s)\n", instPtrType.Elem().Name(), i, methodName, param1Typ)
						for f := 0; f < param1Typ.NumField(); f++ {
							field := param1Typ.Field(f)
							//fieldType := fieldSpec.Type
							tag, existsTag := field.Tag.Lookup(HttpTagName)
							diTag := ParseStructTag(tag)
							def := WebCtlParamDef{
								index:     f,
								fieldSpec: field,
								existsTag: existsTag,
								attrs:     diTag,
							}
							if field.Anonymous && field.Type == webCtxType {
								// extended/embedded struct, retrieve request name and method from http tag
								if existsTag {
									if path := def.preferredText(Ife(forReq, "path", "name"), true, false); len(path) > 0 {
										reqPath = string(handleMethodRegex.Find([]byte(path)))
									}
									if attr, b := diTag.FirstAttrsWithKey("method"); b {
										if len(attr.Val) > 0 {
											reqMethod = strings.ToUpper(attr.Val)
										}
									}
								}
							} else {
								def.preferredName = def.preferredText("name", true, true)
							}
							param1Defs = append(param1Defs, def)
						}
						wrappers = append(wrappers, handlerWrapper{
							strings.ToUpper(reqMethod),
							reqPath,
							func(c *gin.Context) {
								param1InstPtrVal := reflect.New(param1Typ)
								param1InstVal := param1InstPtrVal.Elem()
								get := func(key string, typ reflect.Type) (data any, exists bool) {
									var text string
									if text, exists = c.GetQuery(key); !exists {
										if text, exists = c.GetPostForm(key); !exists {
											if text = c.Param(key); len(text) == 0 {
												return nil, false
											}
										}
									}
									instPtrVal := reflect.New(typ)
									if err := json.Unmarshal([]byte(text), instPtrVal.Interface()); err != nil {
										fmt.Printf("parse key:'%s' with value:'%s' incorrectly, %v\n", key, text, err)
									}
									return instPtrVal.Elem().Interface(), true
								}

								for _, def := range param1Defs {
									fieldSpec := def.fieldSpec
									field := param1InstVal.Field(def.index)
									if fieldSpec.Anonymous && fieldSpec.Type == webCtxType {
										ctx := WebContext{c}
										field.Set(reflect.ValueOf(ctx))
									} else {
										if val, ok := get(def.preferredName, def.fieldSpec.Type); ok {
											fieldName := fieldSpec.Name
											if reflect.TypeOf(val) == def.fieldSpec.Type {
												if len(fieldName) == 0 || fieldName[0] == '_' {
													// ignore
												} else if fieldName[0] >= 'A' && fieldName[0] <= 'Z' {
													field.Set(reflect.ValueOf(val))
												} else {
													dij.SetUnexportedField(field, val)
												}
											}
										}
									}
								}
								//fmt.Printf("I'm in")
								outData := reflect.ValueOf(instPtr).MethodByName(methodName).Call([]reflect.Value{param1InstVal})
								generateOutputData(c, methodName, outData)

							},
						})
					} else {
						if len(reqMethod) == 0 {
							continue
						}
						fmt.Printf("[*%v]'s method %d: func %v(%s)\n", instPtrType.Elem().Name(), i, methodName, param1Typ.Name())
						wrappers = append(wrappers, handlerWrapper{
							strings.ToUpper(reqMethod),
							reqPath,
							func(c *gin.Context) {
								ctx := WebContext{c}
								//fmt.Printf("I'm in")
								outData := reflect.ValueOf(instPtr).MethodByName(methodName).Call([]reflect.Value{reflect.ValueOf(ctx)})
								generateOutputData(c, methodName, outData)
							},
						})
					}
				}
			}
		}
	}
	return wrappers
}

func generateOutputData(c *gin.Context, method string, output []reflect.Value) {
	/*for _, d := range output {
		c.JSON()
	}*/
}

type WebMiddlewareSpec interface {
	iAmAWebMiddleware()
}

type WebMiddleware struct {
}

func (m *WebMiddleware) iAmAWebMiddleware() {

}

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

func IsTypeOfWebContext(typ reflect.Type) bool {
	if typ.Kind() == reflect.Pointer {
		if elemTyp := typ.Elem(); elemTyp.Kind() == reflect.Struct {
			return IsTypeOfWebContext(elemTyp)
		}
		return false
	}
	instPtrValue := reflect.New(typ)
	instIf := instPtrValue.Interface()
	_, ok := instIf.(WebContextSpec)
	return ok
}

type WebServerConfig struct {
	address string
	port    int
	maxConn int
}

//func CreateServeMux(servInstPtr any) *http.ServeMux {
//	servPtrValue := reflect.ValueOf(servInstPtr)
//	servPtrType := servPtrValue.Type()
//	servValue := servPtrValue.Elem()
//	servType := servValue.Type()
//
//	log.Printf("Server type: %v, ptr type: %v", servType, servPtrType)
//	mux := http.NewServeMux()
//	//mux.HandleFunc("/", ws.HandleRequest)
//
//	for j := 0; j < servType.NumField(); j++ {
//		fieldSpec := servType.Field(j)
//		_, existingHttpTag := fieldSpec.Tag.Lookup(HttpTagName)
//		if existingHttpTag {
//			// f := func(writer http.ResponseWriter, r *http.Request) {}
//		}
//	}
//
//	return mux
//}
//
//func LaunchWebServer(webServerType reflect.Type, others ...any) {
//	ref := dij.DependencyReference{}
//	config := WebServerConfig{}
//	for _, other := range others {
//		otherTyp := reflect.TypeOf(other)
//		switch v := other.(type) {
//		case WebServerConfig:
//			config = v
//			ref[WebConfigKey] = v
//		default:
//			log.Println("No ideal about type:", otherTyp)
//		}
//	}
//	instPtr, err := dij.CreateInstance(webServerType, &ref, "")
//	if err != nil {
//		log.Panic(err)
//	}
//
//	addr := fmt.Sprintf("%v:%d", config.address, Ife(config.port <= 0, DefaultWebServerPort, config.port))
//	router := CreateServeMux(instPtr)
//
//	srv := http.Server{
//		ReadHeaderTimeout: time.Second * 5,
//		ReadTimeout:       time.Second * 10,
//		Handler:           router,
//	}
//
//	listener, err := net.Listen("tcp", addr)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	if config.maxConn > 0 {
//		listener = netutil.LimitListener(listener, config.maxConn)
//		//log.Printf("max connections set to %d\n", ws.maxConn)
//	}
//	defer func(listener net.Listener) {
//		err := listener.Close()
//		if err != nil {
//			//
//		}
//	}(listener)
//
//	log.Printf("listening on %s\n", listener.Addr().String())
//
//	go func() {
//		if err := srv.Serve(listener); err != nil && err != http.ErrServerClosed {
//			log.Fatal(err)
//		}
//	}()
//
//	signalChannel := make(chan os.Signal, 1)
//	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
//	<-signalChannel
//
//	log.Printf("interrupted, shutting down")
//
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
//	defer cancel()
//
//	if err := srv.Shutdown(ctx); err != nil {
//		log.Printf("graceful shutdown failed: %v\n", err)
//	}
//}

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

	if err := setupRouterHandlers(instPtr, webServerType, router, &ref); err != nil {
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

func setupRouterHandlers(instPtr any, instType reflect.Type, router gin.IRouter, refPtr dij.DependencyReferencePtr) error {
	routers := router.(gin.IRoutes)
	//
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

	// prepare middlewares
	mwHdlWrappers := map[string]handlerWrapper{}
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
					return fmt.Errorf("fieldSpec(%s) should not be zero", fieldTyp.Name())
				}
				if field.IsExported() {
					fieldIf = fieldValue.Interface()
				} else {
					fieldIf = reflect.NewAt(fieldTyp, fieldValue.Addr().UnsafePointer()).Elem().Interface()
				}
			} else {
				//fmt.Printf("middleware load from dij: %v\n", fieldTyp)
			}
			wrappers := createHandlers(fieldIf, false)
			for _, w := range wrappers {
				fmt.Printf("%v %v\n", w.method, w.path)
				_, exists := mwHdlWrappers[w.path]
				if exists {
					return fmt.Errorf("middleware's handler '%s' is duplicated", w.path)
				}
				mwHdlWrappers[w.path] = w
			}
		}
	}

	// setup current router
	if len(predecessor) != 1 {
		return fmt.Errorf("struct '%s' should embeded web controller or web server.(%d)", instType.Name(), len(predecessor))
	} else {
		field := instType.Field(predecessor[0])
		if tag, exists := field.Tag.Lookup(HttpTagName); exists {
			attrs := ParseStructTag(tag)
			if attr, existingName := attrs.FirstAttrWithValOnly(); existingName {
				router = router.Group(attr.Val)
			} else if attr, exists := attrs.FirstAttrsWithKey("path"); exists {
				router = router.Group(attr.Val)
			}
			routers = router.(gin.IRoutes)
			if attr, exists := attrs.FirstAttrsWithKey("middleware"); exists {
				middlewares := strings.Split(attr.Val, ",")
				for _, m := range middlewares {
					name := strings.TrimSpace(m)
					if w, b := mwHdlWrappers[name]; !b {
						return fmt.Errorf("middleware's handler '%s' doesn't exist", name)
					} else {
						routers = routers.Use(w.handler)
					}
				}
			}
		}
		fmt.Printf("Set router for %v\n", instType)
		setupHandlers(routers, instPtr)
		//ctrl := instPtr.(WebControllerSpec)
		//ctrl.SetupRouter(router, instPtr)
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
					return fmt.Errorf("fieldSpec(%s) should not be zero", fieldTyp.Name())
				}
				if field.IsExported() {
					fieldIf = fieldValue.Interface()
				} else {
					fieldIf = reflect.NewAt(fieldTyp, fieldValue.Addr().UnsafePointer()).Elem().Interface()
				}
			} else {
				//fmt.Printf("extenders load from dij: %v\n", fieldTyp)
			}
			if err := setupRouterHandlers(fieldIf, fieldTyp.Elem(), router, refPtr); err != nil {
				return err
			}
		}
	}

	// all done
	return nil
}

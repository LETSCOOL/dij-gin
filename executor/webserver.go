package executor

import (
	"context"
	"fmt"
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
	TagName              = "http"
	WebConfigKey         = "webserver.config"
)

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
		_, existingHttpTag := fieldSpec.Tag.Lookup(TagName)
		if existingHttpTag {
			// f := func(writer http.ResponseWriter, r *http.Request) {}
		}
	}

	return mux
}

func LaunchWebServer(webServerType reflect.Type, others ...any) {
	ref := map[dij.DependencyKey]any{}
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

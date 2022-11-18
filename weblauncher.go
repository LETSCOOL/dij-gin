// Copyright 2022 Yuchi Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dij_gin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/letscool/lc-go/dij"
	. "github.com/letscool/lc-go/lg"
	"log"
	"reflect"
	"strings"
)

const (
	DefaultWebServerPort = 8000
	HttpTagName          = "http"
	WebConfigKey         = "webserver.config"
	WebSpecRecord        = "webserver.record"
)

func setupHandlers(routes gin.IRoutes, instPtr any) {
	wrappers := GenerateHandlerWrappers(instPtr, HandlerForReq)
	for _, w := range wrappers {
		routes.Handle(w.method, w.path, w.handler)
	}
}

func generateOutputData(c *gin.Context, method string, output []reflect.Value) {
	/*for _, d := range output {
		c.JSON()
	}*/
}

type WebServerConfig struct {
	address           string
	port              int
	maxConn           int
	enabledSpecRecord bool
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
	mwHdlWrappers := map[string]HandlerWrapper{}
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
			wrappers := GenerateHandlerWrappers(fieldIf, HandlerForMid)
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
		ctrl := instPtr.(WebControllerSpec)
		ctrl.SetupRouter(router, instPtr)
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

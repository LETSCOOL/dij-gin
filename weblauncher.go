// Copyright 2022 Yuchi Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dij_gin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/letscool/dij-gin/spec"
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
	WebSpecRecord        = "webserver.spec.record"
)

func PrepareGin(webServerType reflect.Type, others ...any) (*gin.Engine, dij.DependencyReferencePtr, error) {
	ref := dij.DependencyReference{}
	config := NewWebConfig()
	for _, other := range others {
		otherTyp := reflect.TypeOf(other)
		switch v := other.(type) {
		case WebConfig:
			config = &v
		case *WebConfig:
			config = v
		default:
			log.Println("No ideal about type:", otherTyp, " value:", other)
		}
	}
	config.ApplyDefaultValues()
	ref[WebConfigKey] = config
	website := spec.WebSiteSpec{
		Swagger: "2.0",
		Info: &spec.Info{
			License:        nil,
			Contact:        nil,
			Description:    "This site is still under construction.",
			TermsOfService: "",
			Title:          "A dij-gin base API",
			Version:        "0.0.1",
		},
		Host:     Ife(config.Address == "", "localhost", config.Address),
		BasePath: config.BasePath,
		Tags:     nil,
		Schemes:  config.Schemes,
		Paths:    nil,
	}
	ref[WebSpecRecord] = &website

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
	engine, refPtr, err := PrepareGin(webServerType, others...)
	if err != nil {
		return err
	}
	v, _ := refPtr.Get(WebConfigKey)
	config := v.(*WebConfig)

	addr := fmt.Sprintf("%v:%d", config.Address, Ife(config.Port <= 0, DefaultWebServerPort, config.Port))
	return engine.Run(addr)
}

func setupRouterHandlers(instPtr any, instType reflect.Type, router WebRouter, refPtr dij.DependencyReferencePtr) error {
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
					if name := strings.TrimSpace(m); len(name) > 0 {
						if w, b := mwHdlWrappers[name]; !b {
							return fmt.Errorf("middleware's handler '%s' doesn't exist", name)
						} else {
							routers = routers.Use(w.handler)
						}
					}
				}
			}
		}
		fmt.Printf("Set router for %v\n", instType)
		if webRoutes, ok := routers.(WebRoutes); ok {
			setupRoutesHandlers(webRoutes, instPtr, mwHdlWrappers)
			ctrl := instPtr.(WebControllerSpec)
			ctrl.SetupRouter(router, instPtr)
		} else {
			log.Fatalln("IRoutes doesn't have BasePath??? Fix it.")
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

// setupRoutesHandlers set routing path for controller
func setupRoutesHandlers(routes WebRoutes, instPtr any, mwHdlWrappers map[string]HandlerWrapper) {
	wrappers := GenerateHandlerWrappers(instPtr, HandlerForReq)
	for _, w := range wrappers {
		var handlers []gin.HandlerFunc
		for _, name := range w.middlewareNames {
			if name = strings.TrimSpace(name); len(name) > 0 {
				if h, b := mwHdlWrappers[name]; b {
					handlers = append(handlers, h.handler)
				} else {
					log.Fatalf("middleware's handler '%s' doesn't exist", name)
				}
			}
		}
		handlers = append(handlers, w.handler)
		routes.Handle(w.method, w.path, handlers...)
	}
}

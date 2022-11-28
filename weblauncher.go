// Copyright 2022 Yuchi Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dij_gin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
	WebValidator         = "webserver.validator"
	WebDijRef            = "webserver.dij.ref"
)

func PrepareGin(webServerType reflect.Type, others ...any) (*gin.Engine, dij.DependencyReferencePtr, error) {
	if !IsTypeOfWebServer(webServerType) {
		return nil, nil, fmt.Errorf("the type(%v) is not a web server", webServerType)
	}
	ref := dij.DependencyReference{}
	// setup web config
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
	port := Ife(config.Port <= 0, DefaultWebServerPort, config.Port)
	url := fmt.Sprintf("%v:%d/%s", config.Address, port, config.BasePath)
	// setup web spec record, aka. swagger
	website := spec.Openapi{
		Openapi: "3.0.3",
		Info: &spec.Info{
			License:        nil,
			Contact:        nil,
			Description:    "This site is still under construction.",
			TermsOfService: "",
			Title:          "A dij-gin base API",
			Version:        "0.0.1",
		},
		Servers: []spec.Server{
			{
				//Url:         "{schemes}://{addr}:{port}/{basePath}",
				Url: "{schemes}://" + url,
				//Description: "API",
				Variables: map[string]spec.ServerVariable{
					"schemes": {
						Enum:    config.Schemes,
						Default: "http",
					},
					//"addr": {
					//	Default: config.Address,
					//},
					//"basePath": {
					//	Default: config.BasePath,
					//},
					//"port": {
					//	Default: strconv.Itoa(port),
					//},
				},
			},
		},
		Tags:  nil,
		Paths: nil,
	}
	ref[WebSpecRecord] = &website
	// setup validator
	v := validator.New()
	v.SetTagName("binding")
	ref[WebValidator] = v
	// save ref self
	ref[WebDijRef] = &ref
	// create instance
	//dij.EnableLog()
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
				fmt.Printf("%v %v\n", w.ReqMethod(), w.ReqPath())
				_, exists := mwHdlWrappers[w.ReqPath()]
				if exists {
					return fmt.Errorf("middleware's handler '%s' is duplicated", w.ReqPath())
				}
				mwHdlWrappers[w.ReqPath()] = w
			}
		}
	}

	// setup current router
	if len(predecessor) != 1 {
		return fmt.Errorf("struct '%s' should embeded web controller or web server.(%d)", instType.Name(), len(predecessor))
	} else {
		// TODO: base controller has some middlewares installed, and extended controllers also inherited those middlewares, why? fix it ????
		routers := router.(gin.IRoutes)
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
				middlewares := strings.Split(attr.Val, "&")
				for _, m := range middlewares {
					if name := strings.TrimSpace(m); len(name) > 0 {
						if w, b := mwHdlWrappers[name]; !b {
							return fmt.Errorf("middleware's handler '%s' doesn't exist", name)
						} else {
							routers = routers.Use(w.Handler)
						}
					}
				}
			}
		}
		fmt.Printf("Set router for %v\n", instType)
		if webRoutes, ok := routers.(WebRoutes); ok {
			siteSpec := (*refPtr)[WebSpecRecord].(*spec.Openapi)
			setupRoutesHandlers(webRoutes, instPtr, mwHdlWrappers, siteSpec)
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
func setupRoutesHandlers(routes WebRoutes, instPtr any, mwHdlWrappers map[string]HandlerWrapper, siteSpec *spec.Openapi) {
	basePath := routes.BasePath()

	wrappers := GenerateHandlerWrappers(instPtr, HandlerForReq)
	for _, w := range wrappers {
		// process gin structure
		var handlers []gin.HandlerFunc
		for _, name := range w.Spec.MiddlewareNames {
			if name = strings.TrimSpace(name); len(name) > 0 {
				if h, b := mwHdlWrappers[name]; b {
					handlers = append(handlers, h.Handler)
				} else {
					log.Fatalf("middleware's handler '%s' doesn't exist", name)
				}
			}
		}
		handlers = append(handlers, w.Handler)
		routes.Handle(w.UpperReqMethod(), w.ReqPath(), handlers...)
		// process swagger structure
		fullPath, pathParamNames := w.ConcatOpenapiPath(basePath)
		method := w.ReqMethod()
		var parameters spec.ParameterList
		shouldBodyCoding := method == "post" || method == "put" || method == "patch"
		consumeCoding := make([]spec.MediaTypeCoding, 0) // "application/x-www-form-urlencoded", "multipart/form-data", "application/json"
		var objCoding, formCoding int
		for _, fieldDef := range w.Spec.FieldsOfBaseParam {
			fieldSpec := fieldDef.FieldSpec
			fieldSpecType := fieldSpec.Type
			attrs := fieldDef.Attrs
			if fieldSpec.Anonymous && fieldSpecType == WebCtxType {
				vals := attrs.AttrsWithValOnly()
				for _, v := range vals {
					if c, isObj, isCodingAttr := CodingFormAttr(v.Val); isCodingAttr {
						if isObj {
							objCoding++
						} else {
							formCoding++
						}
						consumeCoding = append(consumeCoding, c)
					}
				}
				break
			} else {
				// later
			}
		}
		if objCoding > 0 && formCoding > 0 {
			log.Fatalf("Obj-coding and form-coding should not set at same time.")
		}
		if len(consumeCoding) > 0 && !shouldBodyCoding {
			log.Fatalf("Only post or put method support body coding")
		}
		var formWayCnt, bodyWayCnt int
		for _, fieldDef := range w.Spec.FieldsOfBaseParam {
			fieldSpec := fieldDef.FieldSpec
			fieldSpecType := fieldSpec.Type
			attrs := fieldDef.Attrs
			if fieldSpec.Anonymous && fieldSpecType == WebCtxType {
				// ignore
			} else {
				paramSpec := spec.Parameter{
					Name: fieldDef.PreferredName,
				}
				varKind := spec.GetVariableKind(fieldSpecType)
				if varKind == spec.VarKindUnsupported {
					log.Fatalf("unsupport variable type: %v", fieldSpecType)
				}
				if Contains(pathParamNames, fieldDef.PreferredName) {
					paramSpec.In = InPathWay
				} else {
					if attr, b := attrs.FirstAttrsWithKey("in"); b {
						in := attr.Val
						if !IsCorrectInWay(in) {
							log.Fatalf("unsupported in way: %s", in)
						}
						paramSpec.In = in
					} else {
						switch method {
						case "post", "put", "patch":
							switch varKind {
							case spec.VarKindArray, spec.VarKindObject:
								paramSpec.In = InBodyWay
							default:
								paramSpec.In = InFormWay
							}
						default:
							paramSpec.In = InQueryWay
						}
					}
				}
				paramSpec.ApplyType(fieldSpecType)
				if attrs.ContainsAttrWithValOnly("required") {
					paramSpec.Required = true
				}
				if paramSpec.In == InFormWay {
					formWayCnt++
				} else if paramSpec.In == InBodyWay {
					bodyWayCnt++
				}
				parameters = parameters.AppendParam(&paramSpec)
			}
		}
		if shouldBodyCoding && len(consumeCoding) == 0 {
			if bodyWayCnt > 0 {
				consumeCoding = append(consumeCoding, spec.JsonObject)
			} else {
				consumeCoding = append(consumeCoding, spec.UrlEncoded)
			}
		}
		if formWayCnt > 0 && bodyWayCnt > 0 {
			log.Fatalf("Form way variable and body way variable should not come together")
		}
		if bodyWayCnt > 1 {
			log.Fatalf("Only support one body way variable")
		}

		resp := spec.Response{
			Content: spec.Content{
				"application/json": {
					Schema: &spec.SchemaR{Schema: &spec.Schema{
						Type: "string",
					},
					},
				},
			},
			Description: "ok 200",
		}
		operation := spec.Operation{
			//Consumes:   consumeCoding,
			//Produces:   []spec.MediaTypeCoding{"application/json"},
			Parameters: parameters,
			Responses:  spec.Responses{"200": spec.ResponseR{Response: &resp}},
		}
		siteSpec.AddPathOperation(fullPath, method, operation)
	}
}

func CodingFormAttr(v string) (coding spec.MediaTypeCoding, isObjective bool, isCodingAttr bool) {
	switch v {
	case "form", "multipart":
		return spec.MultipartForm, false, true
	case "urlenc", "urlencoded":
		return spec.UrlEncoded, false, true
	case "json":
		return spec.JsonObject, true, true
	case "xml":
		return spec.XmlObject, true, true
	}
	return "", false, false
}

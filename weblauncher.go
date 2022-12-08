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
	HttpTagName            = "http"
	DescriptionTagName     = "description"
	RefKeyForWebConfig     = "_.webserver.config"
	RefKeyForWebSpecRecord = "_.webserver.spec.record"
	RefKeyForWebValidator  = "_.webserver.validator"
	RefKeyForWebDijRef     = "_.webserver.dij.ref"
)

func PrepareGin(webServerTypeOrInst any, others ...any) (*gin.Engine, dij.DependencyReferencePtr, error) {
	var webServerType reflect.Type
	var webServerInst any
	if typ, ok := webServerTypeOrInst.(reflect.Type); ok {
		webServerType = typ
	} else {
		webServerInst = webServerTypeOrInst
		webServerType = reflect.TypeOf(webServerTypeOrInst).Elem()
	}
	webServerTypeOrInst = nil

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
	ref[RefKeyForWebConfig] = config
	gin.DefaultWriter = config.DefaultWriter
	//
	for k, v := range config.DependentRefs {
		ref[k] = v
	}
	//
	port := Ife(config.Port <= 0, DefaultWebServerPort, config.Port)
	url := fmt.Sprintf("%v:%d/%s", config.Address, port, config.BasePath)
	// setup web spec record, aka. swagger
	if config.OpenApi.Enabled {
		website := spec.Openapi{
			Openapi: "3.0.3",
			Info: &spec.Info{
				License:        nil,
				Contact:        nil,
				Description:    Ife(config.OpenApi.Description == "", "This site is still under construction.", config.OpenApi.Description),
				TermsOfService: "",
				Title:          Ife(config.OpenApi.Title == "", "A dij-gin base API", config.OpenApi.Title),
				Version:        Ife(config.OpenApi.Version == "", "0.0.1", config.OpenApi.Version),
			},
			Servers: []spec.Server{
				{
					//Url:         "{schemes}://{addr}:{port}/{basePath}",
					Url: "{schemes}://" + url,
					//Description: "API",
					Variables: map[string]spec.ServerVariable{
						"schemes": {
							Enum:    config.OpenApi.Schemes,
							Default: config.OpenApi.Schemes[0],
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
		ref[RefKeyForWebSpecRecord] = &website
	}
	// setup validator
	v := validator.New()
	v.SetTagName(config.ValidatorTagName)
	ref[RefKeyForWebValidator] = v
	// save ref self
	ref[RefKeyForWebDijRef] = &ref
	// create instance
	//dij.EnableLog()
	var err error
	if webServerInst != nil {
		var inst any
		inst, err = dij.BuildAnyInstance(webServerInst, &ref, "^")
		if err == nil && inst != webServerInst {
			log.Fatalf("instance is not original instance, it should be a bug.\n%p != %p\n%v",
				inst, webServerInst, err)
		}
	} else {
		webServerInst, err = dij.CreateInstance(webServerType, &ref, "^")
	}
	if err != nil {
		log.Panic(err)
	}

	router := gin.Default()

	if err := setupRouterHandlers(webServerInst, webServerType, router, &ref); err != nil {
		return nil, nil, err
	}

	return router, &ref, nil
}

// LaunchGin launches a web server.
// The webServerTypeOrInst should be a struct type which is embedded WebServer or
// an instance (pointer) of a struct type which is embedded WebServer.
//
// A new web server definition:
//
//	type WebSer struct {
//	  WebServer
//	}
//
//	func main() {
//	  webTyp = reflect.TypeOf(WebSer{})
//	  LaunchGin(webTyp) // launch by type
//
//	  webInst = &WebSer{}
//	  LaunchGin(webInst) // launch by instance
//	}
func LaunchGin(webServerTypeOrInst any, others ...any) error {
	engine, refPtr, err := PrepareGin(webServerTypeOrInst, others...)
	if err != nil {
		return err
	}
	v, _ := refPtr.Get(RefKeyForWebConfig)
	config := v.(*WebConfig)

	addr := fmt.Sprintf("%v:%d", config.Address, Ife(config.Port <= 0, DefaultWebServerPort, config.Port))
	return engine.Run(addr)
}

func setupRouterHandlers(instPtr any, instType reflect.Type, router WebRouter, refPtr dij.DependencyReferencePtr) error {
	rtEnv := ((*refPtr)[RefKeyForWebConfig].(*WebConfig)).RtEnv
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
			wrappers := GenerateHandlerWrappers(fieldIf, HandlerForMid, refPtr)
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
		var apiTag string
		if tag, exists := field.Tag.Lookup(HttpTagName); exists {
			attrs := ParseStructTag(tag)
			if envOnly, ok := attrs.FirstAttrsWithKey("env"); ok {
				if !rtEnv.IsInOnlyEnv(envOnly.Val) {
					return nil
				}
			}
			if apiTagAttr, ok := attrs.FirstAttrsWithKey("tag"); ok {
				apiTag = strings.TrimSpace(apiTagAttr.Val)
			}
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
			setupRoutesHandlers(webRoutes, instPtr, mwHdlWrappers, refPtr, apiTag)
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
func setupRoutesHandlers(routes WebRoutes, instPtr any, mwHdlWrappers map[string]HandlerWrapper, refPtr dij.DependencyReferencePtr, apiTag string) {
	basePath := routes.BasePath()
	wrappers := GenerateHandlerWrappers(instPtr, HandlerForReq, refPtr)
	var openapiSpec *spec.Openapi
	if _, ok := (*refPtr)[RefKeyForWebSpecRecord]; ok {
		openapiSpec = (*refPtr)[RefKeyForWebSpecRecord].(*spec.Openapi)
	}
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

		// check openapi is enabled
		if openapiSpec == nil {
			continue
		}
		// process openapi structure
		fullPath, pathParamNames := w.ConcatOpenapiPath(basePath)
		method := w.ReqMethod()
		var parameters spec.ParameterList
		var bodySchemas []spec.SchemaR
		var reqBody *spec.RequestBodyR
		shouldBodyCoding := method == "post" || method == "put" || method == "patch"
		reqMime := make([]spec.MediaTypeTitle, 0) // "application/x-www-form-urlencoded", "multipart/form-data", "application/json"
		var objCoding, formCoding int
		for _, fieldDef := range w.Spec.InFields {
			fieldSpec := fieldDef.FieldSpec
			fieldSpecType := fieldSpec.Type
			if fieldSpec.Anonymous && fieldSpecType == WebCtxType {
				for _, mt := range fieldDef.SupportedMediaTypesForRequest() {
					if mt.Kind == spec.ObjectiveMediaType {
						objCoding++
					} else {
						formCoding++
					}
					reqMime = append(reqMime, mt.Title)
				}
				if apiTagAttr, ok := fieldDef.Attrs.FirstAttrsWithKey("tag"); ok {
					apiTag = strings.TrimSpace(apiTagAttr.Val)
				}
				break
			} else {
				// later
			}
		}
		if objCoding > 0 && formCoding > 0 {
			log.Fatalf("Obj-coding and form-coding should not set at same time.")
		}
		if len(reqMime) > 0 && !shouldBodyCoding {
			log.Fatalf("Only post or put method support body coding")
		}
		var preferPlainCoding, preferObjCoding int
		for _, fieldDef := range w.Spec.InFields {
			fieldSpec := fieldDef.FieldSpec
			fieldSpecType := fieldSpec.Type
			attrs := fieldDef.Attrs
			if fieldSpec.Anonymous && fieldSpecType == WebCtxType {
				// ignore
			} else {
				var inWay InWay
				varKind := spec.GetVariableKind(fieldSpecType)
				if varKind == spec.VarKindUnsupported {
					log.Fatalf("unsupport variable type: %v", fieldSpecType)
				}
				if Contains(pathParamNames, fieldDef.PreferredName) {
					inWay = InPathWay
				} else {
					if attr, b := attrs.FirstAttrsWithKey("in"); b {
						in := attr.Val
						if !IsCorrectInWay(in) {
							log.Fatalf("unsupported in way: %s", in)
						}
						inWay = in
					} else {
						switch varKind {
						case spec.VarKindArray, spec.VarKindObject:
							preferObjCoding++
						default:
							preferPlainCoding++
						}
						// use default way
						if shouldBodyCoding {
							inWay = InBodyWay
						} else {
							inWay = InQueryWay
						}
					}
				}
				//
				if inWay == InBodyWay {
					// request body
					schema := spec.SchemaR{}
					schema.ApplyType(fieldSpecType)
					bodySchemas = append(bodySchemas, schema)
				} else {
					// parameters
					paramSpec := spec.Parameter{
						Name:        fieldDef.PreferredName,
						In:          inWay,
						Description: fieldDef.Description,
					}
					paramSpec.ApplyType(fieldSpecType)
					if attrs.ContainsAttrWithValOnly("required") {
						paramSpec.Required = true
					}

					parameters = parameters.AppendParam(&paramSpec)
				}
			}
		}
		if shouldBodyCoding && len(reqMime) == 0 {
			if preferObjCoding > 0 {
				reqMime = append(reqMime, spec.JsonObject)
			} else {
				reqMime = append(reqMime, spec.UrlEncoded)
			}
		}
		//if preferPlainCoding > 0 && preferObjCoding > 0 {
		//	log.Fatalf("Form way variable and body way variable should not come together")
		//}
		//if preferObjCoding > 1 {
		//	log.Fatalf("Only support one body way variable")
		//}
		responses := spec.Responses{}
		for _, fieldDef := range w.Spec.OutFields {
			fieldSpec := fieldDef.FieldSpec
			fieldSpecType := fieldSpec.Type
			format := fieldDef.PreferredMediaTypeTitleForResponse()
			schema := spec.SchemaR{}
			if IsError(fieldSpecType) {
				schema.ApplyType(TypeOfWebError)
			} else {
				schema.ApplyType(fieldSpecType)
			}
			content := spec.Content{}
			content[format] = spec.MediaType{Schema: &schema}
			resp := spec.Response{
				Content:     content,
				Description: fieldDef.Description,
			}
			code := fieldDef.PreferredName
			responses[code] = spec.ResponseR{Response: &resp}
		}

		if shouldBodyCoding {
			// At this moment, doesn't support ref RequestBody
			reqBody = &spec.RequestBodyR{
				RequestBody: &spec.RequestBody{},
			}
			var mainSchema spec.SchemaR
			switch len(bodySchemas) {
			case 0:
				mainSchema = spec.SchemaR{}
				mainSchema.ApplyType(reflect.TypeOf(""))
				reqBody.Required = false
			case 1:
				mainSchema = bodySchemas[0]
				reqBody.Required = true
			default:
				mainSchema.ApplyAllOf(bodySchemas...)
				reqBody.Required = true
			}

			for _, coding := range reqMime {
				reqBody.SetMediaType(coding, spec.MediaType{Schema: &mainSchema})
			}
		}

		var tags []string
		if apiTag = strings.Trim(apiTag, "&"); len(apiTag) > 0 {
			tags = strings.Split(apiTag, "&")
		}

		operation := spec.Operation{
			Parameters:  parameters,
			RequestBody: reqBody,
			Responses:   responses,
			Description: w.Spec.Description,
			Tags:        tags,
		}
		openapiSpec.AddPathOperation(fullPath, method, operation)
	}
}

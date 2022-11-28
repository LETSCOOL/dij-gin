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
	"regexp"
	"strings"
)

type HandlerWrapper struct {
	Spec    HandlerSpec
	Handler gin.HandlerFunc
}

func (w *HandlerWrapper) ReqMethod() string {
	return w.Spec.Method
}

func (w *HandlerWrapper) UpperReqMethod() string {
	return w.Spec.UpperMethod()
}

func (w *HandlerWrapper) ReqPath() string {
	return w.Spec.Path
}

func (w *HandlerWrapper) ConcatOpenapiPath(basePath string) (fullPath string, params []string) {
	path := strings.TrimRight(basePath, "/") + "/" + w.ReqPath()
	var pathComps []string
	for _, comp := range strings.Split(path, "/") {
		if len(comp) <= 0 {
			continue
		}
		if comp[0] == ':' || comp[0] == '*' {
			params = append(params, comp[1:])
			pathComps = append(pathComps, "{"+comp[1:]+"}")
		} else {
			pathComps = append(pathComps, comp)
		}
	}
	fullPath = "/" + strings.Join(pathComps, "/")
	return
}

type HandlerSpec struct {
	Purpose           HandlerWrapperPurpose
	MethodType        reflect.Method
	BaseParamType     reflect.Type
	Method            string // lower case, ex: get, post, etc.
	Path              string // lower cast path, does it need to support case-sensitive?
	FieldsOfBaseParam []FieldDefOfHandlerBaseParam
	MiddlewareNames   []string
}

func (s *HandlerSpec) UpperMethod() string {
	return strings.ToUpper(s.Method)
}

type HandlerWrapperPurpose int

const (
	HandlerForReq HandlerWrapperPurpose = iota // for http request
	HandlerForMid                              // for middleware
)

func (h HandlerWrapperPurpose) RegexpText() string {
	switch h {
	case HandlerForReq:
		return `^(get|post|put|patch|delete|head|connect|options|trace)`
	case HandlerForMid:
		return `^(handle)`
	default:
		return ""
	}
}

func (h HandlerWrapperPurpose) Regexp() *regexp.Regexp {
	if text := h.RegexpText(); text != "" {
		return regexp.MustCompile(text)
	}
	return nil
}

func (h HandlerWrapperPurpose) BaseKey() string {
	switch h {
	case HandlerForReq:
		return "path"
	case HandlerForMid:
		return "name"
	default:
		return ""
	}
}

type FieldDefOfHandlerBaseParam struct {
	Index         int
	FieldSpec     reflect.StructField
	ExistsTag     bool           // exists http tag
	Attrs         StructTagAttrs // come from http tag
	PreferredName string
}

func (c *FieldDefOfHandlerBaseParam) preferredText(key string, allowedValOnly bool, allowedFieldName bool) string {
	if c.ExistsTag {
		if allowedValOnly {
			if attr, ok := c.Attrs.FirstAttrWithValOnly(); ok {
				if len(attr.Val) > 0 {
					return attr.Val
				}
			}
		}
		if attr, ok := c.Attrs.FirstAttrsWithKey(key); ok {
			if len(attr.Val) > 0 {
				return attr.Val
			}
		}
	}
	if allowedFieldName {
		return c.FieldSpec.Name
	}
	return ""
}

func GenerateHandlerWrappers(instPtr any, purpose HandlerWrapperPurpose) []HandlerWrapper {
	wrappers := make([]HandlerWrapper, 0)
	instPtrType := reflect.TypeOf(instPtr)
	handleMethodRegex := purpose.Regexp()
	// TODO: how to deal routing for static pages
	for i := 0; i < instPtrType.NumMethod(); i++ {
		method := instPtrType.Method(i)
		if method.IsExported() {
			methodType := method.Type
			if methodType.NumIn() == 2 && methodType.NumOut() <= 1 {
				baseParamType := methodType.In(1)
				if IsTypeOfWebContext(baseParamType) && baseParamType.Kind() == reflect.Struct {
					hdlSpec := HandlerSpec{
						Purpose:       purpose,
						BaseParamType: baseParamType,
						MethodType:    method,
					}
					// Only fit function with one parameter and the parameter extends WebContext.
					methodName := method.Name
					lowerMethodName := strings.ToLower(methodName)
					hdlSpec.Method = string(handleMethodRegex.Find([]byte(lowerMethodName)))
					hdlSpec.Path = lowerMethodName[len(hdlSpec.Method):]
					if baseParamType != WebCtxType {
						// this part extends WebContext, so it also should process tag information
						fmt.Printf("[*%v]'s method %d: func %v(%s)\n", instPtrType.Elem().Name(), i, methodName, baseParamType)
						//fmt.Printf("\t%s\n", baseParamType.Name())
						analyzeBaseParam(baseParamType, purpose, &hdlSpec)

						wrappers = append(wrappers, HandlerWrapper{
							hdlSpec,
							func(c *gin.Context) {
								baseParamInstPtrVal := reflect.New(baseParamType)
								baseParamInstVal := baseParamInstPtrVal.Elem()
								ctx := WebContext{c}
								for _, def := range hdlSpec.FieldsOfBaseParam {
									fieldSpec := def.FieldSpec
									fieldSpecType := fieldSpec.Type
									field := baseParamInstVal.Field(def.Index)
									if fieldSpec.Anonymous && fieldSpecType == WebCtxType {
										field.Set(reflect.ValueOf(ctx))
									} else {
										var val any = nil
										ok := false
										if fieldSpecType.Kind() == reflect.Struct {
											value := reflect.New(fieldSpecType)
											if err := c.ShouldBind(value.Interface()); err == nil {
												val = value.Elem().Interface()
												ok = true
											} else {
												log.Printf("bind type(%v) with json error: %v\n", fieldSpecType, err)
											}
										} else if fieldSpecType.Kind() == reflect.Pointer && fieldSpecType.Elem().Kind() == reflect.Struct {
											value := reflect.New(fieldSpecType.Elem())
											if err := c.ShouldBind(value.Interface()); err == nil {
												val = value.Interface()
												ok = true
											} else {
												log.Printf("bind type(%v) with json error: %v\n", fieldSpecType, err)
											}
										} else {
											in, b := def.Attrs.FirstAttrsWithKey("in")
											val, ok = ctx.GetRequestValueForType(def.PreferredName, fieldSpecType, Ife(b, in.Val, ""))
										}
										if ok {
											fieldName := fieldSpec.Name
											if reflect.TypeOf(val) == fieldSpecType {
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
								outData := reflect.ValueOf(instPtr).MethodByName(methodName).Call([]reflect.Value{baseParamInstVal})
								generateOutputData(c, methodName, outData)
							},
						})
					} else {
						if len(hdlSpec.Method) == 0 {
							continue
						}
						fmt.Printf("[*%v]'s method %d: func %v(%s)\n", instPtrType.Elem().Name(), i, methodName, baseParamType.Name())
						//fmt.Printf("\t%s\n", baseParamType.Name())
						wrappers = append(wrappers, HandlerWrapper{
							hdlSpec,
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

func analyzeBaseParam(baseParamType reflect.Type, purpose HandlerWrapperPurpose, spec *HandlerSpec) {
	baseKey := purpose.BaseKey()
	handleMethodRegex := purpose.Regexp()
	spec.FieldsOfBaseParam = make([]FieldDefOfHandlerBaseParam, 0)
	for f := 0; f < baseParamType.NumField(); f++ {
		field := baseParamType.Field(f)
		tag, existsTag := field.Tag.Lookup(HttpTagName)
		diTag := ParseStructTag(tag)
		def := FieldDefOfHandlerBaseParam{
			Index:     f,
			FieldSpec: field,
			ExistsTag: existsTag,
			Attrs:     diTag,
		}
		if field.Anonymous && field.Type == WebCtxType {
			// extended/embedded struct, retrieve request name and method from http tag
			if existsTag {
				if path := def.preferredText(baseKey, true, false); len(path) > 0 {
					spec.Path = path
				}
				if attr, b := diTag.FirstAttrsWithKey("method"); b {
					if len(attr.Val) > 0 {
						spec.Method = strings.ToLower(string(handleMethodRegex.Find([]byte(attr.Val))))
					}
				}
				if attr, exists := diTag.FirstAttrsWithKey("middleware"); exists {
					middlewares := strings.Split(attr.Val, "&")
					spec.MiddlewareNames = append(spec.MiddlewareNames, middlewares...)
					//fmt.Printf("middlewares: %v, diTag: %v\n", middlewares, diTag)
				}
			}
		} else {
			def.PreferredName = def.preferredText("name", true, true)
			//fmt.Printf("\t%d[%s][%s] %v\n", def.Index, def.PreferredName, def.FieldSpec.Name, def.FieldSpec.Type)
		}
		spec.FieldsOfBaseParam = append(spec.FieldsOfBaseParam, def)
	}
	return
}

func generateOutputData(c *gin.Context, method string, output []reflect.Value) {
	/*for _, d := range output {
		c.JSON()
	}*/
}

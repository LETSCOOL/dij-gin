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
	method  string
	path    string
	handler func(*gin.Context)
}

type HandlerWrapperPurpose int

const (
	HandlerForReq HandlerWrapperPurpose = iota // for http request
	HandlerForMid                              // for middleware
)

func GenerateHandlerWrappers(instPtr any, purpose HandlerWrapperPurpose) []HandlerWrapper {
	var regexText, baseKey string
	switch purpose {
	case HandlerForReq:
		regexText = `^(get|post|put|patch|delete|head|connect|options|trace)`
		baseKey = "path"
	case HandlerForMid:
		regexText = `^(handle)`
		baseKey = "name"
	default:
		return nil
	}
	wrappers := make([]HandlerWrapper, 0)
	webCtxType := reflect.TypeOf(WebContext{})
	instPtrType := reflect.TypeOf(instPtr)
	handleMethodRegex := regexp.MustCompile(regexText)
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
					if param1Typ != webCtxType {
						fmt.Printf("[*%v]'s method %d: func %v(%s)\n", instPtrType.Elem().Name(), i, methodName, param1Typ)
						param1Defs := make([]FieldDefOfHandlerBaseParam, 0)
						for f := 0; f < param1Typ.NumField(); f++ {
							field := param1Typ.Field(f)
							tag, existsTag := field.Tag.Lookup(HttpTagName)
							diTag := ParseStructTag(tag)
							def := FieldDefOfHandlerBaseParam{
								index:     f,
								fieldSpec: field,
								existsTag: existsTag,
								attrs:     diTag,
							}
							if field.Anonymous && field.Type == webCtxType {
								// extended/embedded struct, retrieve request name and method from http tag
								if existsTag {
									if path := def.preferredText(baseKey, true, false); len(path) > 0 {
										reqPath = path
									}
									if attr, b := diTag.FirstAttrsWithKey("method"); b {
										if len(attr.Val) > 0 {
											reqMethod = strings.ToUpper(string(handleMethodRegex.Find([]byte(attr.Val))))
										}
									}
								}
							} else {
								def.preferredName = def.preferredText("name", true, true)
							}
							param1Defs = append(param1Defs, def)
						}
						wrappers = append(wrappers, HandlerWrapper{
							strings.ToUpper(reqMethod),
							reqPath,
							func(c *gin.Context) {
								param1InstPtrVal := reflect.New(param1Typ)
								param1InstVal := param1InstPtrVal.Elem()
								ctx := WebContext{c}
								for _, def := range param1Defs {
									fieldSpec := def.fieldSpec
									fieldSpecType := fieldSpec.Type
									field := param1InstVal.Field(def.index)
									if fieldSpec.Anonymous && fieldSpecType == webCtxType {
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
											val, ok = ctx.GetRequestValueForType(def.preferredName, fieldSpecType)
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
								outData := reflect.ValueOf(instPtr).MethodByName(methodName).Call([]reflect.Value{param1InstVal})
								generateOutputData(c, methodName, outData)
							},
						})
					} else {
						if len(reqMethod) == 0 {
							continue
						}
						fmt.Printf("[*%v]'s method %d: func %v(%s)\n", instPtrType.Elem().Name(), i, methodName, param1Typ.Name())
						wrappers = append(wrappers, HandlerWrapper{
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

type FieldDefOfHandlerBaseParam struct {
	index         int
	fieldSpec     reflect.StructField
	existsTag     bool           // exists http tag
	attrs         StructTagAttrs // come from http tag
	preferredName string
}

func (c *FieldDefOfHandlerBaseParam) preferredText(key string, allowedValOnly bool, allowedFieldName bool) string {
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

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
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var reqRegex, middleRegex, codeRegex *regexp.Regexp
var TypeOfWebError reflect.Type

func init() {
	reqRegex = regexp.MustCompile(`^(get|post|put|patch|delete|head|connect|options|trace)`)
	middleRegex = regexp.MustCompile(`^(handle)`)
	codeRegex = regexp.MustCompile(`^((\w*[\D+|^][2-5]\d{2})|default|([2-5]\d{2}))$`)
	TypeOfWebError = reflect.TypeOf(WebError{})
}

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
	Purpose         HandlerWrapperPurpose
	MethodType      reflect.Method
	BaseParamType   reflect.Type
	Method          string // lower case, ex: get, post, etc.
	Path            string // lower cast path, does it need to support case-sensitive?
	InFields        []BaseParamField
	MiddlewareNames []string
	OutFields       []BaseParamField
	CtxAttrs        StructTagAttrs // tag attr come from the base field in InFields
	Description     string         // description comes from the base field in InFields
}

func (s *HandlerSpec) UpperMethod() string {
	return strings.ToUpper(s.Method)
}

type HandlerWrapperPurpose int

const (
	HandlerForReq HandlerWrapperPurpose = iota // for http request
	HandlerForMid                              // for middleware
)

func (h HandlerWrapperPurpose) Regexp() *regexp.Regexp {
	switch h {
	case HandlerForReq:
		return reqRegex
	case HandlerForMid:
		return middleRegex
	default:
		return nil
	}
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

type BaseParamField struct {
	Index         int
	FieldSpec     reflect.StructField
	ExistsTag     bool           // exists http tag
	Attrs         StructTagAttrs // come from http tag
	PreferredName string
	Description   string
}

func (c *BaseParamField) preferredText(key string, allowedFirstValOnly bool, allowedFieldName bool) string {
	if c.ExistsTag {
		if name, ok := c.Attrs.PreferredName(key, allowedFirstValOnly); ok {
			return name
		}
	}
	if allowedFieldName {
		return c.FieldSpec.Name
	}
	return ""
}

func (c *BaseParamField) SupportedMediaTypesForRequest() []spec.MediaTypeSupport {
	var list []spec.MediaTypeSupport
	for _, v := range c.Attrs.AttrsWithValOnly() {
		if support, ok := spec.GetSupportedMediaType(v.Val); ok && support.Req {
			list = append(list, support)
		}
	}
	return list
}

func (c *BaseParamField) PreferredMediaTypeTitleForResponse() spec.MediaTypeTitle {
	for _, v := range c.Attrs.AttrsWithValOnly() {
		if support, ok := spec.GetSupportedMediaType(v.Val); ok && support.Resp {
			return support.Title
		}
	}
	return GetPreferredResponseFormat(c.FieldSpec.Type)
}

func GetPreferredResponseFormat(typ reflect.Type) spec.MediaTypeTitle {
	switch typ.Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float64, reflect.Float32,
		reflect.String:
		return spec.PlainText
	case reflect.Struct,
		reflect.Array, reflect.Slice:
		return spec.JsonObject
	case reflect.Interface:
		return spec.JsonObject
	case reflect.Pointer:
		return GetPreferredResponseFormat(typ.Elem())
	}
	return spec.PlainText
}

type WebError struct {
	error
	Message string `json:"message"`
	Code    string `json:"code"`
}

func ToWebError(err error, code string) WebError {
	return WebError{
		error:   err,
		Message: err.Error(),
		Code:    code,
	}
}

// GenerateHandlerWrappers generates handler for the instance
// TODO: consider to cache result for same instance.
func GenerateHandlerWrappers(instPtr any, purpose HandlerWrapperPurpose, refPtr dij.DependencyReferencePtr) []HandlerWrapper {
	wrappers := make([]HandlerWrapper, 0)
	instPtrType := reflect.TypeOf(instPtr)
	handleMethodRegex := purpose.Regexp()
	rtEnv := ((*refPtr)[RefKeyForWebConfig].(*WebConfig)).RtEnv
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

					if purpose == HandlerForReq {
						switch methodType.NumOut() {
						case 0: // ignore
						case 1:
							analyzeOutBaseParam(methodType.Out(0), purpose, &hdlSpec)
							// TODO: below handler should deal response
						default:
							log.Fatalf("Handler function(%s) can not return more than one value.(%d)", methodType.Name(), methodType.NumOut())
						}
					}

					if baseParamType != WebCtxType {
						// this part extends WebContext, so it also should process tag information
						fmt.Printf("[*%v]'s method %d: func %v(%s)\n", instPtrType.Elem().Name(), i, methodName, baseParamType)
						//fmt.Printf("\t%s\n", baseParamType.Name())
						analyzeInBaseParam(baseParamType, purpose, &hdlSpec)
						if envOnly, ok := hdlSpec.CtxAttrs.FirstAttrsWithKey("env"); ok {
							if !rtEnv.IsInOnlyEnv(envOnly.Val) {
								continue
							}
						}

						valid := (*refPtr)[RefKeyForWebValidator].(*validator.Validate)

						wrappers = append(wrappers, HandlerWrapper{
							hdlSpec,
							func(c *gin.Context) {
								baseParamInstPtrVal := reflect.New(baseParamType)
								baseParamInstVal := baseParamInstPtrVal.Elem()
								ctx := WebContext{c}
								for _, def := range hdlSpec.InFields {
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
								if err := valid.Struct(baseParamInstPtrVal.Interface()); err != nil {
									webErr := ToWebError(err, strconv.Itoa(http.StatusBadRequest))
									c.JSON(http.StatusBadRequest, webErr)
								} else {
									outData := reflect.ValueOf(instPtr).MethodByName(methodName).Call([]reflect.Value{baseParamInstVal})
									generateOutputData(c, methodName, outData, hdlSpec)
								}
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
								generateOutputData(c, methodName, outData, hdlSpec)
							},
						})
					}
				}
			}
		}
	}
	return wrappers
}

func analyzeInBaseParam(baseParamType reflect.Type, purpose HandlerWrapperPurpose, hdlSpec *HandlerSpec) {
	fieldsCnt := baseParamType.NumField()
	if fieldsCnt == 0 {
		return
	}
	baseKey := purpose.BaseKey()
	handleMethodRegex := purpose.Regexp()
	hdlSpec.InFields = make([]BaseParamField, 0, fieldsCnt)
	for f := 0; f < fieldsCnt; f++ {
		field := baseParamType.Field(f)
		tag, existsTag := field.Tag.Lookup(HttpTagName)
		diTag := ParseStructTag(tag)
		doc := field.Tag.Get(DescriptionTagName)
		def := BaseParamField{
			Index:       f,
			FieldSpec:   field,
			ExistsTag:   existsTag,
			Attrs:       diTag,
			Description: doc,
		}
		if field.Anonymous {
			if field.Type == WebCtxType {
				// extended/embedded struct, retrieve request name and method from http tag
				if existsTag {
					if path := def.preferredText(baseKey, true, false); len(path) > 0 {
						hdlSpec.Path = path
					}
					if attr, b := diTag.FirstAttrsWithKey("method"); b {
						if len(attr.Val) > 0 {
							hdlSpec.Method = strings.ToLower(string(handleMethodRegex.Find([]byte(attr.Val))))
						}
					}
					if attr, exists := diTag.FirstAttrsWithKey("middleware"); exists {
						middlewares := strings.Split(attr.Val, "&")
						hdlSpec.MiddlewareNames = append(hdlSpec.MiddlewareNames, middlewares...)
						//fmt.Printf("middlewares: %v, diTag: %v\n", middlewares, diTag)
					}
				}
				hdlSpec.Description = doc
				hdlSpec.CtxAttrs = def.Attrs
				//if doc != "" {
				//	log.Printf("I Got doc: %s\n", doc)
				//}
			} else {
				log.Fatal("only can embedded WebContext struct.")
			}
		} else {
			def.PreferredName = def.preferredText("name", true, true)
			//fmt.Printf("\t%d[%s][%s] %v\n", def.Index, def.PreferredName, def.FieldSpec.Name, def.FieldSpec.Type)
		}
		hdlSpec.InFields = append(hdlSpec.InFields, def)
	}
	return
}

func analyzeOutBaseParam(baseParamType reflect.Type, _ HandlerWrapperPurpose, hdlSpec *HandlerSpec) {
	if baseParamType.Kind() != reflect.Struct {
		log.Fatalf("only support to return a struct instead of '%v'\n", baseParamType)
	}
	fieldsCnt := baseParamType.NumField()
	if fieldsCnt == 0 {
		return
	}
	hdlSpec.OutFields = make([]BaseParamField, 0, fieldsCnt)
	for f := 0; f < fieldsCnt; f++ {
		field := baseParamType.Field(f)
		if field.Anonymous || !field.IsExported() {
			log.Fatalf("field(%s.%s) of returned struct should not be anonymous or un-exported", baseParamType.Name(), field.Name)
		}
		fieldType := field.Type
		switch fieldType.Kind() {
		case reflect.Array, reflect.Slice, reflect.Map:
		case reflect.Pointer:
			elemType := fieldType.Elem()
			switch spec.GetVariableKind(elemType) {
			case spec.VarKindBase:
				// ok
			case spec.VarKindObject:
				// ok
			default:
				log.Fatalf("unsupport response type: %s, try to use pointer of struct or base type", fieldType.Name())
			}
		case reflect.Interface:
			if !IsError(fieldType) {
				log.Fatalf("unsupport response type: %s, try to use pointer of struct or base type", fieldType.Name())
			}
		default:
			log.Fatalf("unsupport response type: %s, try to use pointer of struct or base type", fieldType.Name())
		}

		tag, existsTag := field.Tag.Lookup(HttpTagName)
		diTag := ParseStructTag(tag)
		doc := field.Tag.Get(DescriptionTagName)
		def := BaseParamField{
			Index:       f,
			FieldSpec:   field,
			ExistsTag:   existsTag,
			Attrs:       diTag,
			Description: doc,
		}

		code := strings.ToLower(def.preferredText("status", true, true))
		re := codeRegex
		if v := re.Find([]byte(strings.ToLower(code))); len(v) >= 3 {
			code = string(v)
			if code != "default" {
				code = code[len(code)-3:]
			}
		} else {
			code = ""
		}
		if code == "" || code == "default" {
			if IsError(fieldType) {
				code = "400"
			} else {
				code = getPreferredResponseCode(hdlSpec.Method)
			}
		}
		def.PreferredName = code
		hdlSpec.OutFields = append(hdlSpec.OutFields, def)
	}
}

func getPreferredResponseCode(method string) string {
	switch method {
	case "get", "head", "trace":
		return "200"
	case "post", "put":
		return "201"
	}
	return "200"
}

func generateOutputData(c *gin.Context, method string, output []reflect.Value, hdlSpec HandlerSpec) {
	if len(output) != 1 || len(hdlSpec.OutFields) == 0 {
		return
	}
	resultValue := output[0]

OutputData:
	for _, field := range hdlSpec.OutFields {
		fieldValue := resultValue.Field(field.Index)
		if fieldValue.IsNil() {
			continue
		}
		format := field.PreferredMediaTypeTitleForResponse()
		code, _ := strconv.Atoi(field.PreferredName)

		var v any
		if typ := field.FieldSpec.Type; typ.Kind() == reflect.Interface {
			// interface kind should only be error
			if IsError(typ) {
				v = ToWebError(fieldValue.Interface().(error), field.PreferredName)
			}
		} else {
			v = fieldValue.Elem().Interface()
		}

		// text format
		if text, ok := v.(string); ok {
			log.Printf("******* %s *********\n", text)
			c.Data(code, string(format), []byte(text))
			break OutputData
		}

		// byte array
		if binary, ok := v.([]byte); ok {
			c.Data(code, string(format), binary)
			break OutputData
		}

		switch format {
		case spec.UrlEncoded:
			// TODO: implement marshal struct
		case spec.PlainText:
			c.String(code, fmt.Sprint(v))
			break OutputData
		case spec.JsonObject:
			c.JSON(code, v)
			break OutputData
		case spec.XmlObject:
			c.XML(code, v)
			break OutputData
		case spec.HtmlPage:
			c.Data(code, string(format), []byte(fmt.Sprint(v)))
			break OutputData
		case spec.OctetStream, spec.PngImage, spec.JpegImage:
			// TODO: implement reader or something?
		}

		log.Printf("unsupported response format(%s)", format)
		break
	}
}

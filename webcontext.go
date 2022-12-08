// Copyright 2022 Yuchi Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dij_gin

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"reflect"
)

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

type InWay = string

const (
	InHeaderWay InWay = "header" // one kind of way for parameter
	InCookieWay InWay = "cookie" // one kind of way for parameter
	InQueryWay  InWay = "query"  // one kind of way for parameter
	InPathWay   InWay = "path"   // one kind of way for parameter
	InBodyWay   InWay = "body"   // one kind of way for request body
)

func IsCorrectInWay(way InWay) bool {
	switch way {
	case InHeaderWay, InQueryWay, InPathWay, InCookieWay, InBodyWay:
		return true
	}
	return false
}

type WebContextSpec interface {
	iAmAWebContext()

	GetRequestHeader(key string) string
}

type WebContext struct {
	*gin.Context
}

func (c *WebContext) iAmAWebContext() {

}

func (c *WebContext) GetRequestValue(key string, instPtr any) (exists bool) {
	var text string
	var err error
	if text, exists = c.GetQuery(key); !exists {
		if text, exists = c.GetPostForm(key); !exists {
			if text = c.Param(key); len(text) == 0 {
				if text = c.GetHeader(key); len(text) == 0 {
					if text, err = c.Cookie(key); err != nil {
						return false
					}
				}
			}
		}
	}
	if str, ok := instPtr.(*string); ok {
		*str = text
		return true
	} else {
		if err := json.Unmarshal([]byte(text), instPtr); err != nil {
			fmt.Printf("parse key:'%s' with value:'%s' incorrectly, %v\n", key, text, err)
		}
		return true
	}
}

func (c *WebContext) GetRequestValueForType(key string, typ reflect.Type, inWay InWay) (data any, exists bool) {
	var text string
	var err error
	switch inWay {
	case InHeaderWay:
		text = c.GetHeader(key)
		if exists = len(text) > 0; !exists {
			return nil, false
		}
	case InQueryWay:
		if text, exists = c.GetQuery(key); !exists {
			return nil, false
		}
	case InPathWay:
		text = c.Param(key)
		if exists = len(text) > 0; !exists {
			return nil, false
		}
	case InCookieWay:
		if text, err = c.Cookie(key); err != nil {
			return nil, false
		}
	case InBodyWay:
		if text, exists = c.GetPostForm(key); !exists {
			return nil, false
		}
	default:
		if len(inWay) > 0 {
			log.Fatalln("Not support data come from this way: " + inWay)
		}
		// guess
		if text, exists = c.GetQuery(key); !exists {
			if text, exists = c.GetPostForm(key); !exists {
				if text = c.Param(key); len(text) == 0 {
					if text = c.GetHeader(key); len(text) == 0 {
						if text, err = c.Cookie(key); err != nil {
							return nil, false
						}
					}
				}
			}
		}
	}

	switch typ.Kind() {
	case reflect.String:
		return text, true
	default:
		instPtrVal := reflect.New(typ)
		if err := json.Unmarshal([]byte(text), instPtrVal.Interface()); err != nil {
			fmt.Printf("parse key:'%s' with value:'%s' incorrectly, %v\n", key, text, err)
		}
		return instPtrVal.Elem().Interface(), true
	}
}

func (c *WebContext) GetRequestHeader(key string) string {
	return c.Request.Header.Get(key)
}

var WebCtxType reflect.Type

func init() {
	WebCtxType = reflect.TypeOf(WebContext{})
}

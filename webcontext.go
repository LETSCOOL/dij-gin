// Copyright 2022 Yuchi Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dij_gin

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
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

type WebContextSpec interface {
	iAmAWebContext()
}

type WebContext struct {
	*gin.Context
}

func (c *WebContext) iAmAWebContext() {

}

func (c *WebContext) GetRequestValue(key string, instPtr any) (exists bool) {
	var text string
	if text, exists = c.GetQuery(key); !exists {
		if text, exists = c.GetPostForm(key); !exists {
			if text = c.Param(key); len(text) == 0 {
				return false
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

func (c *WebContext) GetRequestValueForType(key string, typ reflect.Type) (data any, exists bool) {
	var text string
	if text, exists = c.GetQuery(key); !exists {
		if text, exists = c.GetPostForm(key); !exists {
			if text = c.Param(key); len(text) == 0 {
				return nil, false
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

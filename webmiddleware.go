package dij_gin

// Copyright 2022 Yuchi Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

import "reflect"

func IsTypeOfWebMiddleware(typ reflect.Type) bool {
	if typ.Kind() == reflect.Pointer {
		if elemTyp := typ.Elem(); elemTyp.Kind() == reflect.Struct {
			return IsTypeOfWebMiddleware(elemTyp)
		}
		return false
	}
	instPtrValue := reflect.New(typ)
	instIf := instPtrValue.Interface()
	_, ok := instIf.(WebMiddlewareSpec)
	return ok
}

type WebMiddlewareSpec interface {
	iAmAWebMiddleware()
}

type WebMiddleware struct {
}

func (m *WebMiddleware) iAmAWebMiddleware() {

}

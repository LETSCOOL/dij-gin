// Copyright 2022 Yuchi Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dij_gin

import "reflect"

func IsTypeOfWebServer(typ reflect.Type) bool {
	if typ.Kind() == reflect.Pointer {
		if elemTyp := typ.Elem(); elemTyp.Kind() == reflect.Struct {
			return IsTypeOfWebServer(elemTyp)
		}
		return false
	}
	instPtrValue := reflect.New(typ)
	instIf := instPtrValue.Interface()
	_, ok := instIf.(WebServerSpec)
	return ok
}

type WebServerSpec interface {
	WebControllerSpec
	iAmAWebServer()
}

type WebServer struct {
	WebController
}

func (w *WebServer) iAmAWebServer() {

}

// Copyright 2022 Yuchi Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dij_gin

import (
	"github.com/letscool/dij-gin/spec"
	"reflect"
)

func IsTypeOfWebController(typ reflect.Type) bool {
	if typ.Kind() == reflect.Pointer {
		if elemTyp := typ.Elem(); elemTyp.Kind() == reflect.Struct {
			return IsTypeOfWebController(elemTyp)
		}
		return false
	}
	instPtrValue := reflect.New(typ)
	instIf := instPtrValue.Interface()
	_, ok := instIf.(WebControllerSpec)
	return ok
}

type WebControllerSpec interface {
	iAmAWebController()

	// SetupRouter only be implemented for dynamic routing in runtime.
	SetupRouter(router WebRouter, others ...any)
}

type WebController struct {
	Spec *spec.WebSiteSpec `di:"webserver.spec.record"`
}

func (w *WebController) iAmAWebController() {

}

func (w *WebController) SetupRouter(_ WebRouter, _ ...any) {

}

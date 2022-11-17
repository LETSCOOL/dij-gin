package dij_gin

import (
	"github.com/gin-gonic/gin"
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

type WebRouter interface {
	gin.IRouter
}

type WebControllerSpec interface {
	iAmAWebController()

	// SetupRouter only be implemented for dynamic routing in runtime.
	SetupRouter(router WebRouter, others ...any)
}

type WebController struct {
}

func (w *WebController) iAmAWebController() {

}

func (w *WebController) SetupRouter(_ WebRouter, _ ...any) {

}

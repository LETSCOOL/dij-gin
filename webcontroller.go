package dij_gin

import "reflect"

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
	//SetupRouter(router gin.IRouter, others ...any)
}

type WebController struct {
}

func (w *WebController) iAmAWebController() {

}

package dij_gin

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

package spec

import "reflect"

type VariableKind int

const (
	VarKindUnsupported VariableKind = iota
	VarKindBase
	VarKindObject
	VarKindArray
)

func GetVariableKind(t reflect.Type) VariableKind {
	switch t.Kind() {
	case reflect.Bool:
		return VarKindBase
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return VarKindBase
	case reflect.Float64, reflect.Float32:
		return VarKindBase
	case reflect.String:
		return VarKindBase
	case reflect.Struct, reflect.Map: // TODO: can support map ?
		return VarKindObject
	case reflect.Array, reflect.Slice:
		return VarKindArray
	default:
		return VarKindUnsupported
	}
}

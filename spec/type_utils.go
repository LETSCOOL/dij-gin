// Copyright 2022 Yuchi Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package spec

import (
	. "github.com/letscool/lc-go/lg"
	"reflect"
)

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
	case reflect.Struct /*, reflect.Map*/ : // TODO: can support map ?
		if t == TypeOfTime {
			return VarKindBase
		}
		return VarKindObject
	case reflect.Array, reflect.Slice:
		return VarKindArray
	//case reflect.Pointer:
	//	return GetVariableKind(t.Elem())
	default:
		return VarKindUnsupported
	}
}

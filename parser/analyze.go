package parser

import (
	"fmt"
	"go/types"
)

func analyzeFile(f *file) (*Doc, error) {
	var enums []*Enum
	var curEnum *Enum

	for _, c := range f.Consts {
		if curEnum == nil || c.Type().String() != curEnum.Name {
			t, ok := c.Type().Underlying().(*types.Basic)
			if !ok {
				return nil, fmt.Errorf("gobuf: unsupported const basic type \"%s\"", c.Type())
			}

			kindOk, kind := kindOfType(t)
			if !kindOk && kind != "" {
				continue
			}
			if kind == "" {
				return nil, fmt.Errorf("gobuf: unsupported const type \"%s\"", c.Type())
			}

			curEnum = &Enum{Kind: kind, Name: c.Type().String()}
			enums = append(enums, curEnum)
		}

		curEnum.Values = append(curEnum.Values, &Value{
			Name:  c.Name(),
			Value: c.Val().ExactString(),
		})
	}

	var structs []*Struct

	for name, s := range f.Structs {
		st := &Struct{
			Name:   name,
			Fields: make([]*Field, s.NumFields()),
		}
		structs = append(structs, st)

		for i := 0; i < s.NumFields(); i++ {
			field := s.Field(i)

			st.Fields[i] = &Field{
				Name: field.Name(),
				Type: analyzeType(field.Type()),
			}

			if st.Fields[i].Type == nil {
				return nil, fmt.Errorf("gobuf: unsupported field type \"%s\"", field.Type().String())
			}
		}
	}

	return &Doc{f.Name, f.Package, enums, structs, f.OtherPkg}, nil
}

func analyzeType(t types.Type) *Type {
	array := analyzeArray(t)
	if array != nil {
		return array
	}

	mapType := analyzeMap(t)
	if mapType != nil {
		return mapType
	}

	pointer := analyzePointer(t)
	if pointer != nil {
		return pointer
	}

	return analyzeScalar(t)
}

func analyzeNamed(t types.Type) (string, types.Type) {
	if named, ok := t.(*types.Named); ok {
		return named.String(), named.Underlying()
	}
	return "", t
}

func analyzeArray(t2 types.Type) *Type {
	name, t := analyzeNamed(t2)
	_, isArray := t.(*types.Array)
	_, isSlice := t.(*types.Slice)

	if isArray || isSlice {
		type Array interface {
			Elem() types.Type
		}

		if array, ok := t.(Array); ok {
			var length int

			if array2, ok := t.(*types.Array); ok {
				length = int(array2.Len())
			}

			elem := analyzeType(array.Elem())
			if elem.Kind == UINT8 {
				return &Type{Kind: BYTES, Len: length, Name: name}
			}

			return &Type{Kind: ARRAY, Elem: elem, Len: length, Name: name}
		}
	}

	return nil
}

func analyzeMap(t types.Type) *Type {
	if mapType, ok := t.(*types.Map); ok {
		key := analyzeScalar(mapType.Key())
		elem := analyzeType(mapType.Elem())
		if key != nil && elem != nil {
			return &Type{Kind: MAP, Key: key, Elem: elem}
		}
	}
	return nil
}

func analyzePointer(t types.Type) *Type {
	if pointer, ok := t.(*types.Pointer); ok {
		if t := analyzeScalar(pointer.Elem()); t != nil {
			return &Type{Kind: POINTER, Elem: t}
		}

		if t := analyzeArray(pointer.Elem()); t != nil {
			return &Type{Kind: POINTER, Elem: t}
		}

	}
	return nil
}

func analyzeScalar(t types.Type) *Type {
	name, t2 := analyzeNamed(t)
	switch t2.(type) {
	case *types.Basic:
		basic := t2.(*types.Basic)
		_, kind := kindOfType(basic)
		if kind != "" {
			return &Type{Kind: kind, Name: name}
		}
	case *types.Struct:
		return &Type{Kind: STRUCT, Name: name}
	}

	return nil
}

func kindOfType(t *types.Basic) (bool, string) {
	switch t.Kind() {
	case types.Int:
		return true, INT
	case types.Uint:
		return true, UINT
	case types.Int8:
		return true, INT8
	case types.Uint8:
		return true, UINT8
	case types.Int16:
		return true, INT16
	case types.Uint16:
		return true, UINT16
	case types.Int32:
		return true, INT32
	case types.Uint32:
		return true, UINT32
	case types.Int64:
		return true, INT64
	case types.Uint64:
		return true, UINT64
	case types.Float32:
		return true, FLOAT32
	case types.Float64:
		return true, FLOAT64
	case types.String:
		return true, STRING
	case types.Bool:
		return true, BOOL
	case types.UntypedInt:
		return false, INT
	case types.UntypedString:
		return false, STRING
	}
	return false, ""
}

package cssource

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/ChaimHong/gobuf/parser"

	"strings"
)

func Gen(jsonData []byte) ([]byte, error) {
	var doc parser.Doc

	if err := json.Unmarshal(jsonData, &doc); err != nil {
		return nil, err
	}

	var o writer

	o.Writef("using System;")
	o.Writef("using System.Collections.Generic;")
	o.Writef("using System.IO;")

	for _, e := range doc.Enums {
		o.Writef("enum %s : %s{", e.Name, e.Kind)
		for _, v := range e.Values {
			o.Writef("%s = %s,", v.Name, v.Value)
		}
		o.Writef("}")
	}

	for _, s := range doc.Structs {
		o.Writef("class %s {", s.Name)

		for _, field := range s.Fields {
			if field.Type.Kind == parser.ARRAY {
				if field.Type.Len != 0 {
					o.Writef("public %s %s = new %s[%d];",
						typeName(field.Type), field.Name, typeName(field.Type.Elem), field.Type.Len)
				} else {
					o.Writef("public %s %s = new %s();",
						typeName(field.Type), field.Name, typeName(field.Type))
				}
			} else if field.Type.Kind == parser.MAP {
				o.Writef("public %s %s = new %s();",
					typeName(field.Type), field.Name, typeName(field.Type))
			} else if field.Type.Kind == parser.BYTES && field.Type.Len != 0 {
				o.Writef("public %s %s = new byte[%d];",
					typeName(field.Type), field.Name, field.Type.Len)
			} else {
				o.Writef("public %s %s;", typeName(field.Type), field.Name)
			}
		}

		o.Writef("public int Size() {")
		o.Writef("int size = 0;")
		for _, field := range s.Fields {
			genSizer(&o, "this."+field.Name, field.Type, 1)
		}
		o.Writef("return size;")
		o.Writef("}")

		o.Writef("public int Marshal(byte[] b, int n) {")
		for _, field := range s.Fields {
			genMarshaler(&o, "this."+field.Name, field.Type, 1)
		}
		o.Writef("return n;")
		o.Writef("}")

		o.Writef("public int Unmarshal(byte[] b, int n) {")
		for _, field := range s.Fields {
			genUnmarshaler(&o, "this."+field.Name, field.Type, 1)
		}
		o.Writef("return n;")
		o.Writef("}")

		o.Writef("public void UnmarshalReader(BinaryReader reader) {")
		for _, field := range s.Fields {
			genUnmarshaler(&o, "this."+field.Name, field.Type, 1)
		}
		o.Writef("return ;")
		o.Writef("}")

		o.Writef("}")
	}

	return o.Bytes(), nil

}

type writer struct {
	deepth int
	bytes.Buffer
}

func (w *writer) Writef(format string, args ...interface{}) {
	format = strings.TrimLeft(format, "\t ")

	if format[0] == '}' {
		w.deepth--
	}

	for i := 0; i < w.deepth; i++ {
		w.WriteByte('\t')
	}

	if format[len(format)-1] == '{' {
		w.deepth++
	}

	w.WriteString(fmt.Sprintf(format, args...))
	w.WriteByte('\n')
}

func isNullable(t *parser.Type) bool {
	return t.Kind == parser.POINTER && t.Elem.Kind != parser.STRUCT && t.Elem.Kind != parser.STRING
}

func typeName(t *parser.Type) string {
	if t.Name != "" {
		return t.Name
	}
	switch t.Kind {
	case parser.INT:
		return "long"
	case parser.UINT:
		return "ulong"
	case parser.INT8:
		return "sbyte"
	case parser.UINT8:
		return "byte"
	case parser.INT16:
		return "short"
	case parser.UINT16:
		return "ushort"
	case parser.INT32:
		return "int"
	case parser.UINT32:
		return "uint"
	case parser.INT64:
		return "long"
	case parser.UINT64:
		return "ulong"
	case parser.FLOAT32:
		return "float"
	case parser.FLOAT64:
		return "double"
	case parser.STRING:
		return "string"
	case parser.BYTES:
		return "byte[]"
	case parser.BOOL:
		return "bool"
	case parser.MAP:
		return fmt.Sprintf("Dictionary<%s, %s>", typeName(t.Key), typeName(t.Elem))
	case parser.POINTER:
		if t.Elem.Kind == parser.STRUCT {
			return typeName(t.Elem)
		}
		if t.Elem.Kind == parser.STRING {
			return "string"
		}
		return fmt.Sprintf("Nullable<%s>", typeName(t.Elem))
	case parser.ARRAY:
		if t.Len != 0 {
			return fmt.Sprintf("%s[]", typeName(t.Elem))
		}
		return fmt.Sprintf("List<%s>", typeName(t.Elem))
	default:
		panic("do not support this kind")
	}
	return ""
}

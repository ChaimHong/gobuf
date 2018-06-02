package cssource

import (
	"fmt"

	"github.com/ChaimHong/gobuf/parser"
)

func genUnmarshaler(o *writer, name string, t *parser.Type, n int) {
	if genArrayUnmarshaler(o, name, t, n) {
		return
	}
	if genMapUnmarshaler(o, name, t, n) {
		return
	}
	if genPointerUnmarshaler(o, name, t, n) {
		return
	}
	genScalarUnmarshaler(o, name, t)
}

func genArrayUnmarshaler(o *writer, name string, t *parser.Type, n int) bool {
	if t.Kind == parser.ARRAY {
		if t.Len == 0 {
			o.Writef("{")
			o.Writef("	%s = new %s((int)Gobuf.ReadUvarint(b, ref n));", name, typeName(t))
			o.Writef("	for (var i%d = 0; i%d < %s.Capacity; i%d ++) {", n, n, name, n)
			o.Writef("		%s v%d;", typeName(t.Elem), n)
			genUnmarshaler(o, fmt.Sprintf("v%d", n), t.Elem, n+1)
			o.Writef("		%s.Add(v%d);", name, n)
			o.Writef("	}")
			o.Writef("}")
		} else {
			o.Writef("	%s = new %s[%d];", name, typeName(t.Elem), t.Len)
			o.Writef("	for (var i%d = 0; i%d < %d; i%d ++) {", n, n, t.Len, n)
			genUnmarshaler(o, fmt.Sprintf("%s[i%d]", name, n), t.Elem, n+1)
			o.Writef("	}")
		}
		return true
	}
	return false
}

func genMapUnmarshaler(o *writer, name string, t *parser.Type, n int) bool {
	if t.Kind == parser.MAP {
		o.Writef("{")
		o.Writef("	var cap%d = (int)Gobuf.ReadUvarint(b, ref n);", n)
		o.Writef("	%s = new %s(cap%d);", name, typeName(t), n)
		o.Writef("	for (var i%d = 0; i%d < cap%d; i%d ++) {", n, n, n, n)
		o.Writef("		%s key%d;", typeName(t.Key), n)
		o.Writef("		%s val%d;", typeName(t.Elem), n)
		genScalarUnmarshaler(o, fmt.Sprintf("key%d", n), t.Key)
		genUnmarshaler(o, fmt.Sprintf("val%d", n), t.Elem, n+1)
		o.Writef("		%s.Add(key%d, val%d);", name, n, n)
		o.Writef("	}")
		o.Writef("}")
	}
	return false
}

func genPointerUnmarshaler(o *writer, name string, t *parser.Type, n int) bool {
	if t.Kind == parser.POINTER {
		o.Writef("if (b[n++] != 0) {")
		genScalarUnmarshaler(o, name, t.Elem)
		o.Writef("} else {")
		o.Writef("	%s = null;", name)
		o.Writef("}")
		return true
	}
	return false
}

func genScalarUnmarshaler(o *writer, name string, t *parser.Type) {
	switch t.Kind {
	case parser.BOOL:
		o.Writef("%s = b[n++] == 1;", name)
	case parser.INT8, parser.UINT8:
		o.Writef("%s = (%s)b[n++];", name, typeName(t))
	case parser.INT16, parser.UINT16:
		o.Writef("%s = (%s)Gobuf.ReadUint16(b, ref n);", name, typeName(t))
	case parser.INT32, parser.UINT32:
		o.Writef("%s = (%s)Gobuf.ReadUint32(b, ref n);", name, typeName(t))
	case parser.INT64, parser.UINT64:
		o.Writef("%s = (%s)Gobuf.ReadUint64(b, ref n);", name, typeName(t))
	case parser.INT:
		o.Writef("%s = Gobuf.ReadVarint(b, ref n);", name)
	case parser.UINT:
		o.Writef("%s = Gobuf.ReadUvarint(b, ref n);", name)
	case parser.FLOAT32:
		o.Writef("%s = Gobuf.ReadFloat32(b, ref n);", name)
	case parser.FLOAT64:
		o.Writef("%s = Gobuf.ReadFloat64(b, ref n);", name)
	case parser.BYTES:
<<<<<<< HEAD
		o.Writef("%s = Gobuf.ReadBytes(b, ref n);", name)
=======
		o.Writef("%s = System.Text.Encoding.UTF8.GetString ( Gobuf.ReadBytes(b, ref n) ) ;", name)
		// o.Writef("%s = Gobuf.ReadString(b, ref n);", name)
>>>>>>> fix
	case parser.STRING:
		o.Writef("%s = Gobuf.ReadString(b, ref n);", name)
	case parser.STRUCT:
		o.Writef("%s = new %s();", name, typeName(t))
		o.Writef("n = %s.Unmarshal(b, n);", name)
	}
}

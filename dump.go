package vars

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strconv"
)

func Dump(v ...interface{}) {
	for _, vv := range v {
		fmt.Println(dump(vv, 0))
	}
}

func Sdump(v ...interface{}) string {
	buf := bytes.NewBuffer(nil)
	for _, vv := range v {
		buf.WriteString(dump(vv, 0))
		buf.WriteString("\n")
	}
	return buf.String()
}

func dump(v interface{}, depth int) string {
	var output string
	switch vv := v.(type) {
	case int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64:
		output = dumpInt(vv, depth)
	case string:
		output = dumpString(vv, depth)
	case float32, float64:
		output = dumpFloat(vv, depth)
	case complex64, complex128:
		output = dumpComplex(vv, depth)
	default: // map struct slice channel
		switch reflect.TypeOf(v).Kind() {
		case reflect.Map:
			output = dumpMap(reflect.ValueOf(v), depth)
		case reflect.Struct:
			output = dumpStruct(reflect.ValueOf(v), depth)
		case reflect.Slice:
			output = dumpSlice(reflect.ValueOf(v), depth)
		case reflect.Chan:
			output = dumpChannel(v, depth)
		case reflect.Ptr:
			output = dumpPtr(reflect.ValueOf(v), depth)
		}
	}

	return output
}

func dumpInt(v interface{}, depth int) string {
	buf := bytes.NewBuffer(nil)
	ident(depth, buf)
	switch t := v.(type) {
	case int:
		buf.WriteString("(int) ")
		buf.WriteString(strconv.FormatInt(int64(t), 10))
	case uint:
		buf.WriteString("(uint) ")
		buf.WriteString(strconv.FormatUint(uint64(t), 10))
	case int8:
		buf.WriteString("(int8) ")
		buf.WriteString(strconv.FormatInt(int64(t), 10))
	case uint8:
		buf.WriteString("(uint8) ")
		buf.WriteString(strconv.FormatUint(uint64(t), 10))
	case int16:
		buf.WriteString("(int16) ")
		buf.WriteString(strconv.FormatInt(int64(t), 10))
	case uint16:
		buf.WriteString("(uint16) ")
		buf.WriteString(strconv.FormatUint(uint64(t), 10))
	case int32:
		buf.WriteString("(int32) ")
		buf.WriteString(strconv.FormatInt(int64(t), 10))
	case uint32:
		buf.WriteString("(uint32) ")
		buf.WriteString(strconv.FormatUint(uint64(t), 10))
	case int64:
		buf.WriteString("(int64) ")
		buf.WriteString(strconv.FormatInt(int64(t), 10))
	case uint64:
		buf.WriteString("(uint64) ")
		buf.WriteString(strconv.FormatUint(uint64(t), 10))
	default:
		buf.WriteString("unknown interger type")
	}

	return buf.String()
}

func dumpString(s string, depth int) string {
	buf := bytes.NewBuffer(nil)
	ident(depth, buf)
	buf.WriteString("(string:")
	buf.WriteString(strconv.FormatInt(int64(len(s)), 10))
	buf.WriteString(") ")
	buf.WriteString(`"`)
	buf.WriteString(s)
	buf.WriteString(`"`)
	return buf.String()
}

func dumpFloat(v interface{}, depth int) string {
	buf := bytes.NewBuffer(nil)
	ident(depth, buf)
	switch t := v.(type) {
	case float32:
		buf.WriteString("(float32) ")
		buf.WriteString(strconv.FormatFloat(float64(t), 'f', -1, 32))
	case float64:
		buf.WriteString("(float64) ")
		buf.WriteString(strconv.FormatFloat(float64(t), 'f', -1, 64))
	default:
		buf.WriteString("unkown float type")
	}
	return buf.String()
}

func dumpComplex(v interface{}, depth int) string {
	buf := bytes.NewBuffer(nil)
	ident(depth, buf)
	switch t := v.(type) {
	case complex64:
		buf.WriteString("(complex64) ")
		buf.WriteString("(")
		buf.WriteString(strconv.FormatFloat(float64(real(t)), 'f', -1, 32))
		i := imag(t)
		if i > 0 {
			buf.WriteString("+")
		}
		buf.WriteString(strconv.FormatFloat(float64(imag(t)), 'f', -1, 32))
		buf.WriteString("i")
		buf.WriteString(")")
	case complex128:
		buf.WriteString("(complex128) ")
		buf.WriteString("(")
		buf.WriteString(strconv.FormatFloat(float64(real(t)), 'f', -1, 64))
		i := imag(t)
		if i > 0 {
			buf.WriteString("+")
		}
		buf.WriteString(strconv.FormatFloat(float64(imag(t)), 'f', -1, 64))
		buf.WriteString("i")
		buf.WriteString(")")
	}
	return buf.String()
}

func dumpStruct(v reflect.Value, depth int) string {
	typ := v.Type()
	buf := bytes.NewBuffer(nil)
	ident(depth, buf)
	buf.WriteString("struct(")
	buf.WriteString(typ.Name())
	buf.WriteString(") {\n")
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).CanInterface() {
			ident(depth, buf)
			buf.WriteString("\t")
			buf.WriteString(fmt.Sprintf(`["%s"] =>`, typ.Field(i).Name))
			buf.WriteString("\n")
			buf.WriteString(dump(v.Field(i).Interface(), depth+1))
			buf.WriteString("\n")
		}
	}
	for i := 0; i < depth; i++ {
		buf.WriteString("\t")
	}
	buf.WriteString("}")
	return buf.String()
}

func dumpMap(v reflect.Value, depth int) string {
	typ := v.Type()
	buf := bytes.NewBuffer(nil)
	ident(depth, buf)
	buf.WriteString(fmt.Sprintf("map[%s]%s{ \n", typ.Key().String(), typ.Elem().String()))
	for _, key := range v.MapKeys() {
		ident(depth+1, buf)
		buf.WriteString(fmt.Sprintf("[%s] => \n", printMapKey(key)))
		buf.WriteString(dump(v.MapIndex(key).Interface(), depth+1))
		buf.WriteString("\n")
	}
	buf.WriteString("}")
	return buf.String()
}

func printMapKey(v reflect.Value) string {
	//typ := v.Type()
	//switch typ.Kind() {
	//
	//}
	// TODO: how to print strutc?
	return fmt.Sprintf("%v", v.Interface())
}

func dumpSlice(v reflect.Value, depth int) string {
	typ := v.Type()
	buf := bytes.NewBuffer(nil)
	ident(depth, buf)
	buf.WriteString("slice(")
	buf.WriteString(typ.Elem().String())
	buf.WriteString(": ")
	buf.WriteString(strconv.FormatInt(int64(v.Len()), 10))
	buf.WriteString(": ")
	buf.WriteString(strconv.FormatInt(int64(v.Cap()), 10))
	buf.WriteString(") {\n")
	for i := 0; i < v.Len(); i++ {
		ident(depth+1, buf)
		buf.WriteString(fmt.Sprintf("[%d] => \n", i))
		buf.WriteString(dump(v.Index(i).Interface(), depth+1))
		buf.WriteString("\n")
	}
	ident(depth, buf)
	buf.WriteString("}")
	return buf.String()
}

func dumpChannel(v interface{}, depth int) string {
	buf := bytes.NewBuffer(nil)
	typ := reflect.TypeOf(v)
	val := reflect.ValueOf(v)
	ident(depth, buf)
	buf.WriteString("(")
	buf.WriteString(typ.String())
	buf.WriteString(": ")
	buf.WriteString(strconv.FormatInt(int64(val.Len()), 10)) // len
	buf.WriteString(": ")
	buf.WriteString(strconv.FormatInt(int64(val.Cap()), 10)) //cap
	buf.WriteString(") ")
	buf.WriteString(fmt.Sprintf("%v", v))

	return buf.String()
}

func dumpPtr(v reflect.Value, depth int) string {
	typ := v.Type()
	if typ.Kind() != reflect.Ptr {
		panic("parameter must be reflect.Ptr")
	}

	return dump(v.Elem().Interface(), depth)
}

func ident(depth int, writer io.Writer) {
	for i := 0; i < depth; i++ {
		writer.Write([]byte("\t"))
	}
}

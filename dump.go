package dump

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"unsafe"
)

func Dump(v ...interface{}) {
	//_, file, line, _ := runtime.Caller(1)
	//fmt.Printf("%s:%d ", file, line)
	for _, vv := range v {
		Fdump(os.Stdout, vv)
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

func Fdump(w io.Writer, v ...interface{}) (int, error) {
	return w.Write([]byte(Sdump(v...)))
}

func dump(v interface{}, depth int) string {
	var output string
	switch vv := v.(type) {
	case bool:
		output = dumpBool(vv, depth)
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
			output = dumpStruct(reflect.ValueOf(v), depth, false)
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

func dumpBool(v bool, depth int) string {
	buf := bytes.NewBuffer(nil)
	indent(depth, buf)
	buf.WriteString("(bool) ")
	if v {
		buf.WriteString("true")
	} else {
		buf.WriteString("false")
	}
	return buf.String()
}

func dumpInt(v interface{}, depth int) string {
	buf := bytes.NewBuffer(nil)
	indent(depth, buf)
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
		panic("unknown interger type")
	}

	return buf.String()
}

func dumpString(s string, depth int) string {
	buf := bytes.NewBuffer(nil)
	indent(depth, buf)
	buf.WriteString("(string: ")
	buf.WriteString(strconv.FormatInt(int64(len(s)), 10))
	buf.WriteString(") ")
	buf.WriteString(`"`)
	buf.WriteString(s)
	buf.WriteString(`"`)
	return buf.String()
}

func dumpFloat(v interface{}, depth int) string {
	buf := bytes.NewBuffer(nil)
	indent(depth, buf)
	switch t := v.(type) {
	case float32:
		buf.WriteString("(float32) ")
		buf.WriteString(strconv.FormatFloat(float64(t), 'f', -1, 32))
	case float64:
		buf.WriteString("(float64) ")
		buf.WriteString(strconv.FormatFloat(float64(t), 'f', -1, 64))
	default:
		panic("unkown float type")
	}
	return buf.String()
}

func dumpComplex(v interface{}, depth int) string {
	buf := bytes.NewBuffer(nil)
	indent(depth, buf)
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

func dumpStruct(v reflect.Value, depth int, isPtr bool) string {
	typ := v.Type()
	buf := bytes.NewBuffer(nil)
	indent(depth, buf)
	buf.WriteString("struct(")
	if isPtr {
		buf.WriteString("*")
	}
	buf.WriteString(typ.Name())
	buf.WriteString(") {\n")
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		indent(depth, buf)
		buf.WriteString("\t")
		buf.WriteString(fmt.Sprintf(`[%s] =>`, sprintfStructField(typ.Field(i).Name)))
		buf.WriteString("\n")
		if v.Field(i).CanInterface() {
			buf.WriteString(dump(field.Interface(), depth+1))
		} else {
			buf.WriteString(dump(Interface(v, i), depth+1))
		}
		buf.WriteString("\n")
	}
	for i := 0; i < depth; i++ {
		buf.WriteString("\t")
	}
	buf.WriteString("}")
	return buf.String()
}

func sprintfStructField(field string) string {
	if t := field[0]; t >= 'a' && t <= 'z' {
		return fmt.Sprintf("%s:unexported", field)
	}
	return field
}

// Interface copy the unexported field in a struct, so the value could be `Interfaceable`
func Interface(rv reflect.Value, field int) interface{} {
	var val interface{}
	if rv.CanAddr() {
		rv = rv.Field(field)
		val = reflect.NewAt(rv.Type(), unsafe.Pointer(rv.UnsafeAddr())).Elem().Interface()
	} else {
		rv2 := reflect.New(rv.Type()).Elem()
		rv2.Set(rv)
		rf := rv2.Field(field)
		rf = reflect.NewAt(rf.Type(), unsafe.Pointer(rf.UnsafeAddr())).Elem()
		val = rf.Interface()
	}
	return val
}

func dumpMap(v reflect.Value, depth int) string {
	typ := v.Type()
	buf := bytes.NewBuffer(nil)
	indent(depth, buf)
	buf.WriteString(fmt.Sprintf("map[%s]%s{ \n", typ.Key().String(), typ.Elem().String()))
	for _, key := range v.MapKeys() {
		indent(depth+1, buf)
		buf.WriteString(fmt.Sprintf("[%s] => \n", printMapKey(key)))
		buf.WriteString(dump(v.MapIndex(key).Interface(), depth+1))
		buf.WriteString("\n")
	}
	indent(depth, buf)
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
	indent(depth, buf)
	buf.WriteString("slice(")
	buf.WriteString(typ.Elem().String())
	buf.WriteString(": ")
	buf.WriteString(strconv.FormatInt(int64(v.Len()), 10))
	buf.WriteString(": ")
	buf.WriteString(strconv.FormatInt(int64(v.Cap()), 10))
	buf.WriteString(") {\n")
	for i := 0; i < v.Len(); i++ {
		indent(depth+1, buf)
		buf.WriteString(fmt.Sprintf("[%d] => \n", i))
		buf.WriteString(dump(v.Index(i).Interface(), depth+1))
		buf.WriteString("\n")
	}
	indent(depth, buf)
	buf.WriteString("}")
	return buf.String()
}

func dumpChannel(v interface{}, depth int) string {
	buf := bytes.NewBuffer(nil)
	typ := reflect.TypeOf(v)
	val := reflect.ValueOf(v)
	indent(depth, buf)
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
	if v.Elem().Kind() == reflect.Struct {
		return dumpStruct(v.Elem(), depth, true)
	}
	return dump(v.Elem().Interface(), depth)
}

func indent(depth int, writer io.Writer) {
	for i := 0; i < depth; i++ {
		writer.Write([]byte("\t"))
	}
}

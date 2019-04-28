package vars

import (
	"fmt"
	"testing"
)

func TestDumpInt(t *testing.T) {
	var (
		a        = 1
		b int    = 2
		c uint   = 3
		d int8   = 4
		e uint8  = 5
		f int16  = 6
		g uint16 = 7
		h int32  = 8
		i uint32 = 9
		j int64  = 10
		k uint64 = 11
	)
	Dump(a, b, c, d, e, f, g, h, i, j, k)
}

func TestDumpString(t *testing.T) {
	var str string = "hello, 世界  \t\naaa"
	Dump(str)
}

func TestDumpFloat(t *testing.T) {
	var (
		a float64 = 1.244211
		b float64 = -12.31231
		c float32 = 1.21123
		d float32 = -1.12123121
		e         = -1231.23
		f         = 11.1423
	)
	Dump(a, b, c, d, e, f)
}

func TestDumpComplex(t *testing.T) {
	var a complex128 = 1i + 2
	var b complex64 = 12.1 - 1231.12i
	var c complex128 = -1231.12 + 123i

	Dump(a, b, c)
}

func TestDumpChannel(t *testing.T) {
	type T struct {
		F1 int
		F2 string
		F3 float64
	}

	a := make(chan int, 1)
	b := make(chan bool)
	c := make(chan struct{}, 101)
	d := make(chan T, 10)
	c <- struct{}{}
	c <- struct{}{}
	c <- struct{}{}

	Dump(a, b, c, d)
}

func TestDumpStruct(t *testing.T) {
	type Home struct {
		price   int
		Address string
	}
	type Person struct {
		Name string
		Age  int
		Home Home
	}

	p := Person{
		"aaa",
		10,
		Home{
			Address: "aaaaa",
		},
	}

	Dump(p)
}

func TestDumpPtr(t *testing.T) {
	type T struct {
		F1 string
		F2 int
	}

	Dump(&T{})
}

func TestDumpSlice(t *testing.T) {
	var arr []int = []int{1, 2, 3}
	var arr2 []interface{} = []interface{}{1, 2.1, 1 + 2i, "hello world", []interface{}{0.11, 2, 3}}
	Dump(arr, arr2)
}

func TestDumpMap(t *testing.T) {
	m := make(map[interface{}]interface{})
	m[100] = "helloworld"
	m["aaa"] = 1
	m[1.1] = 1 + 2i
	Dump(m)
}

func BenchmarkDump(b *testing.B) {
	var a string = "hello world"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Dump(a)
	}

	//100000	     15661 ns/op
}

func BenchmarkDump2(b *testing.B) {
	var a string = "hello world"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fmt.Println(a)
	}
	//  200000	      9949 ns/op
}

func TestSdump(t *testing.T) {
	t.Log(Sdump(1,"2",3.123))
}

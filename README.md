# dump

格式化打印go变量，类似于php的var_dump()函数。

#### number
```go
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
//  (int) 1
//  (int) 2
//  (uint) 3
//  (int8) 4
//  (uint8) 5
//  (int16) 6
//  (uint16) 7
//  (int32) 8
//  (uint32) 9
//  (int64) 10
//  (uint64) 11
```


#### string
```go
var str string = "hello, 世界  \t\naaa"
Dump(str)

//(string: 20) "hello, 世界  	
//  aaa"
```

#### complex
```go
var a complex128 = 1i + 2
var b complex64 = 12.1 - 1231.12i
var c complex128 = -1231.12 + 123i

Dump(a, b, c)
//(complex128) (2+1i)
//(complex64) (12.1-1231.12i)
//(complex128) (-1231.12+123i)
```


#### channel 
```go
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

//  (chan int: 0: 1) 0xc0000160e0
//  (chan bool: 0: 0) 0xc00003c600
//  (chan struct {}: 3: 101) 0xc00003c660
//  (chan vars.T: 0: 10) 0xc0000340c0
```

#### struct
```go
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

// struct(Person) {
//   	["Name"] =>
//   	(string: 3) "aaa"
//   	["Age"] =>
//   	(int) 10
//   	["Home"] =>
//   	struct(Home) {
//   		["Address"] =>
//   		(string: 5) "aaaaa"
//   	}
// }

```


#### slice
```go
var arr []int = []int{1, 2, 3}
var arr2 []interface{} = []interface{}{1, 2.1, 1 + 2i, "hello world", []interface{}{0.11, 2, 3}}
Dump(arr, arr2)

/**
slice(int: 3: 3) {
    [0] => 
    (int) 1
    [1] => 
    (int) 2
    [2] => 
    (int) 3
}
slice(interface {}: 5: 5) {
    [0] => 
    (int) 1
    [1] => 
    (float64) 2.1
    [2] => 
    (complex128) (1+2i)
    [3] => 
    (string:11) "hello world"
    [4] => 
    slice(interface {}: 3: 3) {
        [0] => 
        (float64) 0.11
        [1] => 
        (int) 2
        [2] => 
        (int) 3
    }
}
**/
```


#### map
```go
m := make(map[interface{}]interface{})
m[100] = "helloworld"
m["aaa"] = 1
m[1.1] = 1 + 2i
Dump(m)

/**
map[interface {}]interface {}{ 
	[100] => 
	(string: 10) "helloworld"
	[aaa] => 
	(int) 1
	[1.1] => 
	(complex128) (1+2i)
}
*/
```

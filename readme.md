
# Golang interface internal

## Golang interface assignment and assertion

Golang interface is implicit and does a lot of complicated stuff under hood. 

Normally, We use interface for assignment and assertion, e.g.

```
type myStruct struct{ name string }
type myInterface interface{ method() }

func (myStruct) method() {} // myStruct implement myInterface

func main() {

	var s myStruct = myStruct{name: "Stone"}

	// interface assingment
	var i myInterface = s // assignment of struct
	var any interface{} = i // assignment of interface

	// interface assertion
	v1 := i.(myStruct)  // assertion of struct
	v2 := any.(myInterface) // assertion of interface

	fmt.Printf("v1 = %v, v2 = %v", v1, v2)
}
``` 

Please reference [Go Data Structures: Interfaces](https://research.swtch.com/interfaces) first

## real code to show how it runs
```
package main

import "fmt"

type heator interface {
	heat()
}

type coolor interface {
	cool()
}

type someThing struct {
	name string
}

func (s someThing) heat() {
	fmt.Println("I am heating with name = " + s.name)
}

func (s someThing) cool() {
	fmt.Println("I am cooling with name = " + s.name)
}

func main() {
	some := someThing{name: "stone"}

	var h heator = some
	h.heat()

	var c coolor = some // NOTE: var c cooler = h, compile failed
	c.cool()

	var any interface{} = some // NOTE: var any interface{} = h, can be compiled and has the same result

	v0, ok0 := any.(someThing)
	if ok0 {
		fmt.Printf("any, empty interface, assert struct someThing, type = %T, val = %v\n", v0, v0)
	}

	v1, ok1 := any.(coolor)
	if ok1 {
		fmt.Printf("any, empty interface, assert interface coolor, type = %T, val = %v\n", v1, v1)
	}

	v2, ok2 := any.(heator)
	if ok2 {
		fmt.Printf("any, empty interface, assert interface heator, type = %T, val = %v\n", v2, v2)
	}

	// any = c, compile OK
	// h = c, compile fail
	v3, ok3 := h.(coolor)
	if ok3 {
		fmt.Printf("h, interface heator, assert interface coolor, type = %T, val = %v\n", v3, v3)
	}
}
```

## Run Result

go run main.go 

```
I am heating with name = stone
I am cooling with name = stone

any, empty interface, assert struct someThing, type = main.someThing, val = {stone}
any, empty interface, assert interface coolor, type = main.someThing, val = {stone}
any, empty interface, assert interface heator, type = main.someThing, val = {stone}
h, interface heator, assert interface coolor, type = main.someThing, val = {stone}
```

NOTE: github has more code than above, so it needs to run like this  

go run interface_internal.go pointer_interface.go 

## My Opinion & Guess

### Interface composed of two internal references 

Under hood, the interface is composed of two referecnes, each reference is one word size in memory. 

You can imagine the references similiar to Java reference or C++ smart pointer. 

### The first reference: itable + value type

#### dynamic itable

The first reference is the dynamic itable which is built dynamiclly, i.e. in runtime.

You can imagine itable similiar to Java Interface or C++ vtable, but Java and C++ implement it staticlly, i.e. in compile time.

Itable will be built first time with assignment, then be cached. 

The complexity is O(m+n), m == the number of concrete struct methods, n == the number of interface methods. 

#### value type

The first reference has one more field, the value type, which is paired with the second reference to describe the data.

After assignment, the value type in first reference is unchangeable until another assignment to the interface variable, just like itable.

### The second reference: the copy of concrete value

The second reference is the copy of the concrete value, which can not be an interface. [Interfaces do not hold interface values](https://blog.golang.org/laws-of-reflection)

I do not try the pointer to interface, but I think it is illegal as interface itself. Note: pointer to interface is rarely used.

Concrete value could be:
1. the copy of the memory with the type of struct OR 
2. the copy of the memory with the pointer, which points to an entity of type of struct 
3. can not be any interface (and pointer to an interface)

The following code which is part of the above can demonstrate
```
	v3, ok3 := h.(coolor)
	if ok3 {
		fmt.Printf("h, interface heator, assert interface coolor, type = %T, val = %v\n", v3, v3)
	}
```

If h does not hold the copy of real value with the type struct of someThing, 

it can not assert sucessuflly for the interface coolor.

### assginment to interface with concrete value or another interface

```
some := someThing{name: "Stone:}

var h1 heator
var h2 heator
var any interface{}

h1 = some
any = some

h2 = any  
```

h1 and h2 have the same concrete value, i.e. "Stone", but they are different copies

#### the concrete value can be pointer to struct or struct itself if the receiver is struct

check pointer_interface.go for demonstration

```
type itfer interface {Dummy()}

type foo struct{}

func (f foo) Dummy() {}

func main() {
	var f1 foo
	var f2 *foo = &foo{}

	var i1 itfer = f1
	fmt.Printf("i1 from f1, type = %T\n", i1)
	var i2 itfer = f2
	fmt.Printf("i2 from f2, type = %T\n", i2)
}
```

#### the concrete value must be pointer to struct if the receiver is a pointer

check pointer_interface.go for demonstration

```
type itfer interface {Dummy()}

type bar struct{}

func (b *bar) Dummy() {}

func main() {
	var b1 bar
	var b2 *bar = &b1

	// var i3 itfer = b1	// NOTE: compile fail
	var i4 itfer = b2
	fmt.Printf("i4 from b2, type = %T\n", i4)
}
```

### Type assertion with interface, 

e.g. v := i.(Type) when Type is interface

From the above, it is a matching game.

If Type is an interface, Golang tries to match the total methods in Type, which is interface, to the concrete value of i. 

If all matched, it returns the concrete value without panic, i.e. a new copy to v.

### interface assignment with interface

Although h, variable of heator interface, can asssert type of cooler interface,

You can not assign c to h when in compile time. It will fail.

But you can assgin c to empty interface, it passes the compilation.

Golang compilation just check validation in this senario by checking the latter interface include the method of the prev one.

### implict assignment

When it runs with function parameter or return result, the interface assignment is implicit.

e.g. 

```
type foo struct{ name string }

func help(i interface{}) string {
	return "def"
}

func toInterface(f foo) interface{} {
	return f  // an implicit assignment for return
}

func main() {
	var f foo = foo{"abc"}

	fmt.Println(help(f))  // an implimit assignment for the parameter in help()

	i := toInterface(f)
	if v, ok := i.(foo); ok {
		fmt.Println(v)
	}
}
```

### More trick when the concrete is a pointer

check pointer_interface.go

### other build-in type

other build-in types like stirng, slice, map, numerics, is same like struct.

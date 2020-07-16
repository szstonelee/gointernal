
# interface internal

## code
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
	some := someThing{"stone"}

	var h heator = some
	h.heat()

	var c coolor = some // NOTE: var c cooler = h, can not be compiled because interface can not be receiver
	c.cool()

	var any interface{} = h // NOTE: var any interface{} = some, can be compiled and has the same result

	v0, ok0 := any.(someThing)
	if ok0 {
		fmt.Printf("any has concrete someThing, type = %T, val = %v\n", v0, v0)
	}

	v1, ok1 := any.(coolor)
	if ok1 {
		fmt.Printf("any has interface coolor, type = %T, val = %v\n", v1, v1)
	}

	v2, ok2 := any.(heator)
	if ok2 {
		fmt.Printf("any has interface heator, type = %T, val = %v\n", v2, v2)
	}

	v3, ok3 := h.(coolor)
	if ok3 {
		fmt.Printf("heator has interface coolor, type = %T, val = %v\n", v3, v3)
	}
}
```

## Run Result

```
I am heating with name = stone
I am cooling with name = stone
any has concrete someThing, type = main.someThing, val = {stone}
any has interface coolor, type = main.someThing, val = {stone}
any has interface heator, type = main.someThing, val = {stone}
heator has interface coolor, type = main.someThing, val = {stone}
```

## My Opinion

### Interface composed of two references internal 

Under hood, each reference is one word size in memory. 

You can imagine reference similiar to Java reference or C++ pointer. 

The first reference is the dynamic itable which is built dynamiclly, i.e. in runtime.

You can imagine itable similiar to Java Interface or C++ vtable, but Java and C++ implement it staticlly in compile time.

Itable is built first time with usage, i.e. the assignment to the interface variable, then be cached. 

The complexity is O(m+n), m == the number of concrete struct methods, n == the number of interface methods. 

### The second reference is always the concrete value

The second reference is the copy of the concrete value, which is not an interface. 

It could be a type of struct or pointer which could point to a struct (but not an interface).

The part of code above can demonstrate
```
	v3, ok3 := h.(coolor)
	if ok3 {
		fmt.Printf("heator has interface coolor, type = %T, val = %v\n", v3, v3)
	}
```

### Type assertion, e.g. v = i.(Type) when Type is interface

From the above, it is a matching game.

If Type is an interface, go try to match the total methods in Type to the second reference of i. 

If it can match all the methonds in the concrete value in i, i.e the second reference, return the concrete value without panic.

### Special for pointer

check pointer_interface.go
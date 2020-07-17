
# Golang interface internal

## Golang interface assign and assertion

Golang interface is something which run not as simple as the code looks like, for example

```
// NOTE: the following code is pseduo
var s struct {} = {...}
var i1 interface one {...}
var i2 interface two {...}

// And empty interface is more complicated 
var any interface{}

// assingment, two modes
// mode1: struct to interface
i1 = s
// mode2: interface to interface
i2 = i1

// interface type asserttion, two modes
// mode1: interface assertion with struct (In future, there are points)
v1 := i1.(Type of struct)
// mode2: interface assertion with interface 
v2 := i1.(Type of interface)
``` 

Please reference [Go Data Structures: Interfaces](https://research.swtch.com/interfaces) first

## real code to expose the internal
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
	fmt.Printf("reveal interface assignment and assertion ......\n\n")

	some := someThing{name: "stone"}

	var h heator = some
	h.heat()

	var c coolor = some // NOTE: var c cooler = h, can not be compiled because interface can not be receiver
	c.cool()
	fmt.Println()

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

go run interface_internal.go (if you use the above code)

or 

go run interface_internal.go pointer_interface.go (if you use the code in GitHub which has one more source file)

```
I am heating with name = stone
I am cooling with name = stone
any has concrete someThing, type = main.someThing, val = {stone}
any has interface coolor, type = main.someThing, val = {stone}
any has interface heator, type = main.someThing, val = {stone}
heator has interface coolor, type = main.someThing, val = {stone}
```

## My Opinion & Guess

### Interface composed of two references internal 

Under hood, interface is composed of two referecnes, each reference is one word size in memory. 

You can imagine Golang interface internal reference similiar to Java reference or C++ smart pointer. 

The first reference is the dynamic itable which is built dynamiclly, i.e. in runtime.

You can imagine itable similiar to Java Interface or C++ vtable, but Java and C++ implement it staticlly, i.e. in compile time.

Itable will be built first time with usage, i.e. at the time of assignment to the interface variable, then be cached. 

The complexity is O(m+n), m == the number of concrete struct methods, n == the number of interface methods. 

### The second reference is always the concrete value

The second reference is the copy of the concrete value, which is not an interface. 

It could be a type of struct or pointer which points to a struct (but not an interface).

The part of code above can demonstrate
```
	v3, ok3 := h.(coolor)
	if ok3 {
		fmt.Printf("heator has interface coolor, type = %T, val = %v\n", v3, v3)
	}
```

### Type assertion, e.g. v := i.(Type) when Type is interface

From the above, it is a matching game.

If Type is an interface, Golang tries to match the total methods in Type to the second reference of i. 

If type assertion can match all the methonds in the concrete value in i, i.e the second reference, it returns the concrete value without panic.

### interface assign to interface

Although h, variable of heator interface, can asssert type of cooler interface,

You can not assign c to h. In compile time, it fails.

But you can assgin c to empty interface, it pass compile.

Golang just check validation of the interface type assign between each other 

by which one has more methods in interface definition, not the real concrete value within the interface variable.

### More trick when the concrete is a pointer

check pointer_interface.go
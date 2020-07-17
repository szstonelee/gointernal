
# Golang interface internal

## Golang interface assignment and assertion

Golang interface is implicit and does a lot of complicated stuff under hood. 

Normally, We use interface for assignment and assertion, e.g.

```
// NOTE: the following code is pseduo
var s struct {} = {...}
var i1 interface one {...}
var i2 interface two {...}

// And empty interface makes more complicated 
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
	fmt.Printf("reveal interface assignment and assertion ......\n\n")

	some := someThing{name: "stone"}

	var h heator = some
	h.heat()

	var c coolor = some // NOTE: var c cooler = h, compile failed
	c.cool()

	var any interface{} = some // NOTE: var any interface{} = h, can be compiled and has the same result

	v0, ok0 := any.(someThing)
	if ok0 {
		fmt.Printf("\nany, empty interface, assert struct someThing, type = %T, val = %v\n", v0, v0)
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

go run interface_internal.go pointer_interface.go (if you use the code in GitHub which has one more source file, pointer_interface.go)

```
reveal interface assignment and assertion ......

I am heating with name = stone
I am cooling with name = stone

any, empty interface, assert struct someThing, type = main.someThing, val = {stone}
any, empty interface, assert interface coolor, type = main.someThing, val = {stone}
any, empty interface, assert interface heator, type = main.someThing, val = {stone}
h, interface heator, assert interface coolor, type = main.someThing, val = {stone}
```

## My Opinion & Guess

### Interface composed of two internal references 

Under hood, the nterface is composed of two referecnes, each reference is one word size in memory. 

You can imagine the references similiar to Java reference or C++ smart pointer. 

### The first reference is itable (plus value type)

The first reference is the dynamic itable which is built dynamiclly, i.e. in runtime.

You can imagine itable similiar to Java Interface or C++ vtable, but Java and C++ implement it staticlly, i.e. in compile time.

Itable will be built first time with assignment, then be cached. 

The complexity is O(m+n), m == the number of concrete struct methods, n == the number of interface methods. 

NOTE 1: when it runs with function parameter or return result, the assignment may be implicit.

like 

```
type foo struct{ name string }

func help(i interface{}) string {
	return "def"
}

func toInterface(f foo) interface{} {
	return f
}

func main() {
	var f foo = foo{"abc"}

	fmt.Println(help(f))

	i := toInterface(f)
	if v, ok := i.(foo); ok {
		fmt.Println(v)
	}
}
```

NOTE 2:

The first reference has one more field, i.e. the value type, which is paired with the second reference to describe the data.

After assignment, the value type in first reference is static until another assignment to the interface variable. 

### The second reference is always the copy of concrete value

The second reference is the copy of the concrete value, which can not be an interface. [Interfaces do not hold interface values](https://blog.golang.org/laws-of-reflection)

I do not try the value of pointer to interface, but I think it is illegal with interface. And note: pointer to interface is rarely used.

Concrete value could be:
1. type of struct OR 
2. pointer to a struct 
3. but can not be interface (and pointer to an interface) (check pointer_interface.go)

The following code which is part of the above can demonstrate
```
	v3, ok3 := h.(coolor)
	if ok3 {
		fmt.Printf("h, interface heator, assert interface coolor, type = %T, val = %v\n", v3, v3)
	}
```

If h not save the real value, it can not assert sucessuflly for the interface coolor.

### assginment to interface with concrete value or another interface is same

```
some := someThing{name: "Stone:}

var h1 heator
var h2 heator
var any interface{}

h1 = some
any = some

h2 = any  // NOTE: h2 has the same concrete value of h1, which is "Stone"
```

h1 and h2 have the same concrete value, i.e. "Stone", but they are different copies

#### the concrete value can be pointer to struct or struct itself if the receiver is struct it self

check pointer_interface.go for demonstration

#### the concrete value must be pointer to struct if the receiver is a pointer

check pointer_interface.go for demonstration

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
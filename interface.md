
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

## Code to show how it runs
```
package main

import "fmt"

type heater interface {
	heat()
}

type cooler interface {
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

	var h heater = some
	h.heat()

	var c cooler = some // NOTE: var c cooler = h, compile failed
	c.cool()

	var any interface{} = some // NOTE: var any interface{} = h, can be compiled and has the same result

	v0, ok0 := any.(someThing)
	if ok0 {
		fmt.Printf("any, empty interface, assert struct someThing, type = %T, val = %v\n", v0, v0)
	}

	v1, ok1 := any.(cooler)
	if ok1 {
		fmt.Printf("any, empty interface, assert interface cooler, type = %T, val = %v\n", v1, v1)
	}

	v2, ok2 := any.(heater)
	if ok2 {
		fmt.Printf("any, empty interface, assert interface heater, type = %T, val = %v\n", v2, v2)
	}

	// any = c, compile OK
	// h = c, compile fail
	v3, ok3 := h.(cooler)
	if ok3 {
		fmt.Printf("h, interface heater, assert interface cooler, type = %T, val = %v\n", v3, v3)
	}
}
```

## Run Results

go run main.go 

```
I am heating with name = stone
I am cooling with name = stone

any, empty interface, assert struct someThing, type = main.someThing, val = {stone}
any, empty interface, assert interface cooler, type = main.someThing, val = {stone}
any, empty interface, assert interface heater, type = main.someThing, val = {stone}
h, interface heater, assert interface cooler, type = main.someThing, val = {stone}
```

NOTE: github has more code than above, so it needs to run like this  

go run interface_internal.go pointer_interface.go 

## My Opinion & Guess

### Interface composed of two internal references 

Under hood, the interface is composed of two referecnes, each reference is one word size in memory. 

You can imagine the references similar to Java reference or C++ smart pointer. 

### The first reference: itable + value type

#### dynamic itable

The first reference is the dynamic itable which is built dynamically, i.e. in runtime.

You can imagine itable similar to Java Interface or C++ vtable, but Java and C++ implement it statically, i.e. in compile time.

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
	v3, ok3 := h.(cooler)
	if ok3 {
		fmt.Printf("h, interface heater, assert interface cooler, type = %T, val = %v\n", v3, v3)
	}
```

If h does not hold the copy of real value with the type struct of someThing, 

it can not assert sucessuflly for the interface cooler.

### Assginment to interface with concrete value or another interface

```
some := someThing{name: "Stone:}

var h1 heater
var h2 heater
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

[The reason for the difference treatment for pointer and concrete is here](https://golang.org/doc/faq#Functions_methods)

### Type assertion with interface

e.g. v := i.(Type) when Type is interface type

From the above, it is a matching game.

If Type is an interface, Golang tries to match the total methods in Type, which is interface, to the concrete value of i. 

If all matched, it returns the concrete value without panic, i.e. a new copy to v.

A trick for anonymous interface

```
v, ok := x.(interface{ F() (int, error) })
if ok {
	v.F()	// no panic
}
```
Anonymous interface is better than named interface when named interface is from other package.

Because other package could changed the signature, we can use anonymous interface for the decouple.

Another sample is from [Effitive Go](https://golang.org/doc/effective_go.html#interface_conversions), code with some change shows like this
```
package main

// Stringer is an interface
type Stringer interface {
	String() string
}

type myStruct struct{}

func (myStruct) String() string {
	return "it is of myStruct type"
}

func f(value interface{}) string {
	switch str := value.(type) {
	case string:
		// a, ok := str.(Stringer)	// invalid because str is not an interface
		return str
		
	case Stringer:
		_, ok := str.(myStruct) // right now, str is an interface!!!! and myStruct must implement String(), otherwise, compile failed
		if ok {
			return "also myStruct " + str.String()
		} else {
			return str.String()
		}

	case int:
		return "abc"

	default:
		return "value is not string-like"
	}
}

func main() {

}
```

### Interface assignment with interface

Although h, variable of heater interface, can asssert type of cooler interface,

you can not assign c to h when in compile time. It will fail.

But you can assgin c to empty interface, it passes the compilation.

Golang compilation just check validation in this senario by checking the latter interface include the method of the prev one.

### Be careful of implict assignment

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

There are more info about implicit for interface, check [Golang nil](nil.md)

### More trick when the concrete is a pointer

check pointer_interface.go

### Other build-in type

other build-in types like stirng, slice, map, numerics, bool, byte are same as struct.

### Interface with nil

[check here](nil.md)

### interface comparison

```
package main

import "fmt"

type myStruct struct {
	id int
}

type myStruct2 struct {
	id int
}

type myInterface interface {
	method()
}

func (myStruct) method() {}

func main() {
	var a myStruct = myStruct{id: 1}
	var b myStruct = myStruct{id: 2}
	var c myStruct = myStruct{id: 1}
	var d myStruct2 = myStruct2{id: 1}

	if a == b {
		fmt.Println("a == b")
	} else {
		fmt.Println("a != c")
	}

	if a == c {
		fmt.Println("a == c")
	}

	// if a == d {	// compile failed

	var ia myInterface = a
	var ib myInterface = b
	var ic myInterface = c

	if ia == ib {
		fmt.Println("ia == ib")
	} else {
		fmt.Println("ia != ib")
	}

	if ia == ic {
		fmt.Println("ia == ic")
	}

	if ia == c {
		fmt.Println("ia == c")
	}

	// if i1 == d {	// compile failed
	fmt.Println("d = ", d)

	var any interface{} = myStruct{id: 1}
	if any == ia {
		fmt.Println("empty interface comparison with another interface")
	}
	if any == a {
		fmt.Println("empty interface comparison with another struct, case 1, equal")
	}
	if any == d {
		fmt.Println("empty interface comparison with another struct, case 2, equal")
	} else {
		fmt.Println("empty interface comparison with another struct, case 2, not equal")
	}
}
```

Run result 
```
a != c
a == c
ia != ib
ia == ic
ia == c
d =  {1}
empty interface comparison with another interface
empty interface comparison with another struct, case 1, equal
empty interface comparison with another struct, case 2, not equal
```

Interface comparison with another interface or concrete value, 

it will return true by two condition is matched

1. value type same
2. concrete value same (though different memory address)

But note, can not compare interface with func concrete 

the following code will panic at run time

```
	var i1 interface{} = func(int) {}
	var i2 interface{} = func(int) {}

	if i1 == i2 {
		fmt.Println("i1 == i2")
	}
```

# name for interface

## er way

As a convention, we suffix "er" to a phrase to declare it as interface. For example, Reader to read, Writer to write.

If the word which is appended by "er" is a verb, it is natural. But if it is a noun or adj, it is a little wierd.

That is why sometimes, we swap the noun word for the verb word in a phrase
```
type ByteReader interface {
    ReadByte() (c byte, err error)
}
```

It is simple, but some times it is not easy because some words are not easy to suffix "er". For example, 

```
Server or Info
```

or for Generate, which one is natural, Generater or Generator?

I choose Generater for consistency.

Sometimes, you can imagine a noun as a verb. String has a usage as verb, but in common sense, it is a noun. Stringer is a popular interface in Golang because we treat string as a verb, like Java's Object::toString(). If you treat noun as verb, it is easier to append "er".

## patch way

If appending-er looks not good, there are two patch ways for the naming of interface

1. Prefix I
```
Request -> IRequest
```
Requester is more consistent than IRequest. But someone feel Requester sounds like a person, not an interface.

I prefer Requester.

2. Suffix Interface
```
Request -> RequestInterface
Info -> InfoInterface
```

If 'er' suffix is not good, I will use 'Interface' suffix.

## no way

Sometimes, no er, no I, no Interface for interface, even in the standard library.

For examp, [net/Conn](https://golang.org/pkg/net/#Conn)

Otherwise, it would be Conner or Connectioner.

No way is not a good way. I treat it as non-usual-way. For the legacy code or common interface like net/Conn.

## fake way

Sometimes, er does not mean interface.

For example, in bufio package

```
type ReadWriter struct {
	*Reader  // *bufio.Reader
  *Writer  // *bufio.Writer	
}
```

It makes me confused. It is very rare to use pointer to interface. Actually you can treat pointer to interface as wrong code.

[But after the package](https://golang.org/pkg/bufio/#Reader), we see
```
type Reader struct {}
type Writer struct {}
```

The Reader and Writer in bufio package is struct, no interface. 

So the pointers in ReadWriter, which is a struct too in package bufio, make sense.

But I do not feel good. It is a violation of er convention. I do not like the fake way.

Sometimes, Golang prefers short to verbose to sacrifice the meaning of long phrase. I think Golang goes too far for the idea of short.

It would be better if using #or# to replace #er#.


# struct conversion

```
package main

import "fmt"

type myStruct1 struct {
	name string
}

type myStruct2 struct {
	name string
}

type myStruct3 struct {
	name string
	age  int
}

type myStruct4 struct {
	nameAlias string
}

func main() {
	a1 := myStruct1{name: "aaa"}
	a2 := myStruct2{name: "bbb"}

	a1 = myStruct1(a2)
	fmt.Println(a1) // will print bbb

	// a1 = myStruct1(myStruct3{name: "xxx"})	// illegal
	a3 := myStruct4{nameAlias: "ccc"}
	// a1 = myStruct1(a3)	// illegal
	fmt.Print(a3)
}

```

# function converstion

```
package main

import "fmt"

type myFunc func(int, string) int

func main() {
	var f1 myFunc

	f2 := func(int, string) int { return int(1) }
	f1 = f2
	fmt.Println(f1) // will print function address

	f3 := func(int32, string) int { return int(1) }
	// f1 = f3	// illegal, int32 != int
	// fmt.Println(f3)	// can not print anonymous func
	fmt.Println(f3(5, "abc")) // anonymous func must be used at least once
}
```

[More about func](func.md)

# interface conversion

interface convsion is a big topic, [check here](interface.md)

# explicit type different, but internal type same

```
package main

import "fmt"

func main() {
	type myString string

	a := "abc" // a is string type

	var b myString

	// b = a	// illegal
	b = myString(a)

	fmt.Println(b) // will print "abc"

	fmt.Printf("type of a = %T, type of b = %T\n", a, b) // type is different, a is "string" while b is "main.myString"
}
```
You now can know why it is different for the interface implementation.

Because interface is based on type, check [inteface internal](interface.md) for more info.

There is a trick. You can convert explicit type to other type which implement an interface to [invoke the method](https://golang.org/doc/effective_go.html#conversions).


# Golang pointer is like C++'s pointer

```
var p *int
var b int = 1
p = &b
*p = 2
fmt.Println(b)  // will print 2
```

# Puzzle: pointer with interface

```
package main

import "fmt"

type myStruct struct {
	id int
}

type myInterface interface {
	method()
}

func (myStruct) method() {}

func main() {
	var b myStruct = myStruct{id: 1}

	var p1 *myStruct = &a

	var i1 myInterface = b
	var i2 myInterface = p1

	v1 := i1.(myStruct)
	v1.id = 2
	fmt.Println(b.id)

	v2 := i2.(*myStruct)  // v2 := i2.(myStruct) will panic
	v2.id = 3
	fmt.Println(b.id)
}
```

The run result is

```
1
3
```

Interface's concrete value is a copy.

If the copy is the struct, after new assignment, it has nothing to do with the original one, i.e. b.

If the copy is the pointer, though interface assertion is another copy, but all copies are related to the address of b.

It is similar to the following c++'s code
```
int b = 1
int* p1 = &b  // p1 hold the address of b
int* p2 = p1  // p2 is another pointer, but it holds the address of b
int* p3 = p2
*p3 = <new val>
b will be the <new val>
```
# pointer with nil

[Check here](nil.md)
# Golang pointer is like C++'s pointer

```
	var p *int
	var a int = 1
	p = &a
	*p = 2
	fmt.Println(a)  // will print 2
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
	var a myStruct = myStruct{id: 1}

	var p1 *myStruct = &a

	var i1 myInterface = a
	var i2 myInterface = p1

	v1 := i1.(myStruct)
	v1.id = 2
	fmt.Println(a.id)

	v2 := i2.(*myStruct)  // v2 := i2.(myStruct) will panic
	v2.id = 3
	fmt.Println(a.id)
}
```

The run result is

```
1
3
```

Interface's concrete value is a copy.

If the copy is the struct, after new assignment, it has nothing to do with the original one, i.e. a.

If the copy is the pointer, though interface assertion is another copy, but all copies are reltated to the address of a.

It is similar to the following c++'s code
```
int a = 1
int* p1 = &a  // p1 hold the address of a
int* p2 = p1  // p2 is another pointer, but it holds the address of a
int* p3 = p2
*p3 = <new val>
a will be the <new val>
```

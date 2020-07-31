# Golang construction and initializaion of object

In other language, like C++, you need the `new` operator to construct object in heap (if not optimized by compiler), 

or no new to construct the object in stack. 

```
// C++ code
class MyClass {
public:
  MyClass(int v) : val_(v) {}
private:
  int val_;
}

void foo() {
  MyClass* p = new MyClass(100); // in heap

  MyClass b(200);  // in stack, after exit foo(), b will disappear

  delete p; // if no delete, memory will leak
}
```

In Java, every object (no primitives like int, bool, long) is allocated in heap by new (if not optimized by compiler for escape analyzation)

Otherwise, a reference in Java is null which means there is no object associated with the reference.

NOTE: You can treat Java Clone(), ClassLoader, Reflection, Auto-box, assingnment of String literal as a special construction way.

In Golang, it is a little different.

We can construct the object by new or var or literal

## new

```
struct myStruct struct {
  name string,
  value int
}

var a *myStruct = new(myStruct)
```

## var

```
var b myStruct
```

## literal assignment

```
c := myStruct{} // equivalent to var c myStruct
```
NOTE: var c []int is not equivalent to c := []int{}, see the following section of make().

The difference is that *new* returns pointer while *var* and *literal* return no pointer. You can do

```
func init() *myStruct {
  p := new(myStruct)
  p.name = "Stone"
  p.value = 100
  return p
}
```

the same as 

```
func init() *myStruct {
  return &myStruct{name: "Stone", value: 100}
}
```

## new with pointer

For pointer, it is trivial and subtle.

```
var p1 *myStruct
fmt.Println(*p1)   // will panic, the object of myStruct does not exist, but the pointer p1 exists

var p2 *myStruct = new(myStruct)
fmt.Println(*p2)   // will not panic, the object of myStruct is constructed by new, and p2 is the address to the object

var p3 *int
fmt.Println(*p3)  // will panic

var p4 *int = new(int)
fmt.Prointln(*p4) // will not panic

var p5 *string
fmt.Println(*p5)  // will panic

var p6 *string = new(string)
fmt.Println(*p6)  // will not panic
```

The pointer itself in Golang is an object. But its purpose is pointed to another object. If it is not pointed to another object, the pointer is nil. When we reference the pointer of nil, it will panic.

So if we use new() for a pointer, actually, there are two objects constructed. One is the new object, the other is the pointer.

## heap or stack

When object constructed in Golang, where it lives, stack or heap?

It depends. For Golang, it will try its best to construct the objects in stack, because stack costs much less than heap. Even with new(), the object may exist in stack. That is one reason for why Golang is quicker and has less overhead of GC than Java.

NOTE: Java tries to do the same efficient way by escape analyzation like Golang, but it is harder for Java to construct the object in stack because Golang is simpler than Java.

## make() is for the underlying initialization for map, slice and channel

Because map, slice and channel has two layers, after the construction, there may be a following initialization for these kinds of object.

Check 
1. [slice internal](slice.md)
2. [chanel internal](channel.md)

The top layer is a data structure for abstraction or logic description. 

The top layer always exists.

The underlying layer is the real data for the type. For slice, it is an array. For map, it is a hash map.

When map, slice, channel are declared at the beginning, they only has the top layer, but no underlying layer.

If only the top layer exist, it is nil.

After initialization, like make() does or by literal assignment, the underlying layer is constructed.

```
var a []myStruct = make(myStruct, 5, 10)      // length = 5, capcaity = 10
// will print false
fmt.Println(a == nil) 

var b []int
// will print true, the top layer exist, but the underlying layer, i.e. the int array, does not exist
fmt.Println(b == nil)  

var c []int = []int{} // it equals to c := []int{} which is idiomatic
// will print false, which equals to c := make([[]int, 0, 0]) which is idiomatic
fmt.Println(c == nil) 
```

check [nil](nil.md) for more details

# Golang init order

Each variable in a source file which is outside of function is like a global instance in C in one module.

In C/C++, there are no guarantee for the order of the glabal instance construction. So C/C++ suggests you init and use them from class methods as static variables. [one exmaple for C++](https://stackoverflow.com/questions/3746238/c-global-initialization-order-ignores-dependencies/3746249#3746249), 

[For example for Golang](https://stackoverflow.com/questions/24790175/when-is-the-init-function-run) 
```
package main

import "fmt"

var WhatIsThe = AnswerToLife()

func AnswerToLife() int {
	return 2
}

func init() {
	WhatIsThe = 1
}

func main() {
	if WhatIsThe == 1 {
		fmt.Println("WhatIsThe value is from init()")
	} else if WhatIsThe == 2 {
		fmt.Println("WhatIsThe value is is from var, which call AnswerToLife()")
	} else {
		fmt.Printf("WhatIsThe is anything else, %v\n", WhatIsThe)
	}
}
```
It will print WhatIsThe value is from init().

The order of var in global and init() for each source file with multi import is [the same post in StackOverflow](https://stackoverflow.com/questions/24790175/when-is-the-init-function-run)







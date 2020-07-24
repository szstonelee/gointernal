
# Slice internal

[Check the post first: Go Slices: usage and internals](https://blog.golang.org/slices-intro#:~:text=Slice%20internals&text=It%20consists%20of%20a%20pointer,referred%20to%20by%20the%20slice.)

Slice has two layers.

The first layer is always there, if imagined in C++, it looks like
```
// C++ code to imagine the slice data structure in Golang
struct slice {
  void* _ptr;     // which points to an backed array
  int _capacity;  // the real size of the array
  int _start;     // the start index in the array, which is zero index of the slice
  int _len;       // the size for the slice, NOTE: not the size of array
};
```

The second layer is the backed array. It is pointed to by _ptr.

It may not exist. If array does not exist, i.e. _ptr == nullptr, the slice is nil.
```
// Golang code
var a []int   // a is nil
```

a is nil which looks like in C++ code
```
// C++ code
a = struct slice {
  _ptr = nullptr;   // so in Golang, a is nil
  _capacity = 0;
  _start = 0;
  _len = 0;
};
```

After make(), we create the backed array
```
// Golang code
a := make([]int, 5, 10)
```

If in C++, it looks like
```
// C++ code
a = struct slice {
  _ptr = new int[10];   // the allocated memory is the second layer, which is pointed to by _ptr
  _capacity = 10;
  _start = 0;
  _len = 5;
};
```

After a initiazlized by make(), if we code a[2:9] in Golang, it looks like
```
// C++ code
a[2:9] = struct slice {
  _ptr = a->_ptr; // the backed array does not change
  _capacity = a->_capacity;   // capacity does not change
  _start = 2; // start index is from 2
  _len = 7;   // 9-2=7, the slice's size is 7, only have room for 7 elements
};
```

Then 
```
// Golang code
a = a[2:9]  // means assignment(copy) the a[2:9] struct to a
```

So a now looks like 
```
// in C++ code
a = struct {
  _ptr = the original backed array with size of 10
  _capacity = the original of 10
  _start = new position of 2 in the array as 0 index of the slice
  _len = new len of 7 for the slice
};
```

so 
```
// in Golang
a[1] = 99
```

equals to 
```
// C++ code
a->ptr_[2+1] = 99;
```

# example 1
## code
```
package main

import "fmt"

func f(s []int) {
	s = s[:6]
	for i := 3; i < 6; i++ {
		s[i] = i + 1  // write 4, 5, 6
	}
}

func printSlice(s []int) {
	fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
}

func main() {
	a := make([]int, 6, 6)

	for i := 0; i < 3; i++ {
		a[i] = i + 1
	}

	printSlice(a)
	a = a[:3]
	printSlice(a)

	f(a)

	printSlice(a)
	a = a[:6]
	printSlice(a)
}
```

## result
```
len=6 cap=6 [1 2 3 0 0 0]
len=3 cap=6 [1 2 3]
len=3 cap=6 [1 2 3]
len=6 cap=6 [1 2 3 4 5 6]
```

In f(), s is diffenent from a in terms of memory address. In Golang, ervery parameter is passed by copy.

But s has the same internal value of a, i.e. s->_ptr == a->_ptr.

When we save 4, 5, 6 to s, it saves the value to the backed array, which is the same array of a.

So the last line output is [1 2 3 4 5 6].

# example 2

## code 
```
package main

import "fmt"

func f(s []int) {
	s = make([]int, 6, 6)   // NOTE: the only changed code if compared with example 1
	for i := 3; i < 6; i++ {
		s[i] = i + 1  // write 4, 5, 6
	}
}

func printSlice(s []int) {
	fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
}

func main() {
	a := make([]int, 6, 6)

	for i := 0; i < 3; i++ {
		a[i] = i + 1
	}

	printSlice(a)
	a = a[:3]
	printSlice(a)

	f(a)

	printSlice(a)
	a = a[:6]
	printSlice(a)
}
```

## result
```
len=6 cap=6 [1 2 3 0 0 0]
len=3 cap=6 [1 2 3]
len=3 cap=6 [1 2 3]
len=6 cap=6 [1 2 3 0 0 0]
```

The result is differnet from example 1. The last line output is [1 2 3 0 0 0] if compared to example 1's [1 2 3 4 5 6].

Why?

Because in f(), the _ptr in s is changed to totally new array which is created by make().

s->_ptr in f() is differnt from a->_ptr in main(), so the backed arrays are different.

In example 1, the _ptr in s and a is same, i.e. the backed array does not change.


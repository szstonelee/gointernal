
# Slice internal

[Check the post first: Go Slices: usage and internals](https://blog.golang.org/slices-intro#:~:text=Slice%20internals&text=It%20consists%20of%20a%20pointer,referred%20to%20by%20the%20slice.)

Slice has two layers.

The first layer is always there, if imagined in C++, it looks like
```
// C++ code to imagine the slice data structure in Golang
struct slice {
  void* _ptr;     // which points to a backed array
  int _capacity;  // the real size of the array
  int _start;     // the start index in the array, which is zero index of the slice
  int _len;       // the size for the slice, NOTE: not the size of array
};
```

The second layer is the backed array. It is pointed to by _ptr.

It may not exist. If array does not exist, i.e. _ptr == nullptr, the slice is nil.
```
var b []int   // b is nil
```

b is nil which looks like in C++ code
```
// C++ code
b = struct slice {
  _ptr = nullptr;   // so in Golang, b is nil
  _capacity = 0;
  _start = 0;
  _len = 0;
};
```

After make(), we create the backed array
```
// Golang code
b := make([]int, 5, 10)
```

If in C++, it looks like
```
// C++ code
b = struct slice {
  _ptr = new int[10];   // the allocated 10-integer memory is the second layer, which is pointed to by _ptr
  _capacity = 10;
  _start = 0;
  _len = 5;
};
```

After b is initiazlized by make(), if we code b[2:9] in Golang, it looks like
```
// C++ code
b[2:9] = struct slice {
  _ptr = b->_ptr; // the backed array does not change
  _capacity = b->_capacity;   // capacity does not change
  _start = 2; // start index is from 2
  _len = 7;   // 9-2=7, the slice's size is 7, i.e. room for 7 elements
};
```

Then 
```
// Golang code
b = b[2:9]  // means assig (copy) the b[2:9] struct to b
```

So b now looks like after b = b[2:9]
```
// in C++ code
b = struct {
  _ptr = the original backed array with size of 10
  _capacity = the original of 10
  _start = new position of 2 in the array as 0 index of the slice
  _len = new len of 7 for the slice
};
```

so after b = b[2:9], if we do
```
// in Golang
b[1] = 99
```
it equals to 
```
// C++ code
b->ptr_[2+1] = 99;
```

# Example one
## Code
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
	b := make([]int, 6, 6)

	for i := 0; i < 3; i++ {
		b[i] = i + 1
	}

	printSlice(b)
	b = b[:3]
	printSlice(b)

	f(b)

	printSlice(b)
	b = b[:6]
	printSlice(b)
}
```

## Result and explanation
```
len=6 cap=6 [1 2 3 0 0 0]
len=3 cap=6 [1 2 3]
len=3 cap=6 [1 2 3]
len=6 cap=6 [1 2 3 4 5 6]
```

In f(), s is diffenent from b in terms of memory address. 

In Golang, every parameter is passed by value, i.e. a copy.

But s has the same internal value of b, i.e. s->_ptr == b->_ptr.

When we save 4, 5, 6 to s, it saves to the same internal backed array of b.

So the last line output is [1 2 3 4 5 6].

# Example two

## Code 
```
package main

import "fmt"

func f(s []int) {
	s = make([]int, 6, 6)   // NOTE: the only changed code if compared with example one
	for i := 3; i < 6; i++ {
		s[i] = i + 1  // write 4, 5, 6
	}
}

func printSlice(s []int) {
	fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
}

func main() {
	b := make([]int, 6, 6)

	for i := 0; i < 3; i++ {
		b[i] = i + 1
	}

	printSlice(b)
	b = b[:3]
	printSlice(b)

	f(b)

	printSlice(b)
	b = b[:6]
	printSlice(b)
}
```

## Result and explanation
```
len=6 cap=6 [1 2 3 0 0 0]
len=3 cap=6 [1 2 3]
len=3 cap=6 [1 2 3]
len=6 cap=6 [1 2 3 0 0 0]
```

The result is differnet from example one. The last line output is [1 2 3 0 0 0] if compared to example one's [1 2 3 4 5 6].

Why?

Because in f(), the _ptr in s is changed to a totally new array which is created by make().

After make(), s->_ptr in f() is different from b->_ptr, so the backed arrays are different.

In example one, the _ptr in s and b is same, i.e. the backed array does not change.

# Slice nil

Note:

var b int[] is differnt from b := []int{}

[Check nil for more details](nil.md)
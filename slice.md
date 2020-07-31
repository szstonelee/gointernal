
# Slice internal

[Check the post first: Go Slices: usage and internals](https://blog.golang.org/slices-intro#:~:text=Slice%20internals&text=It%20consists%20of%20a%20pointer,referred%20to%20by%20the%20slice.)

Slice has two layers.

The first layer is always there, if imagined in C++, it looks like
```
// C++ code for imagining the data structure
struct slice {
	void* _ptr;	// which points to a backed array
	int _start;	// the start index in the array which means zero index of the slice
	int _len;	// the length of the slice, how many elements can be saved in the slice 
	int _capacity;  // the capacity for the slice, it equals capacity of array minus _start
};
```
NOTE: _start can not be decremented in future, and always _len <= _capacity.

The second layer is the backed array. It is pointed to by _ptr.

The backed array may not exist. If it does not exist, i.e. _ptr == nullptr, the slice is nil.
```
var b []int   // b is nil
```

b is nil which looks like in C++ code
```
// C++ code
b = struct slice {
	_ptr = nullptr; 
	_start = 0;
	_len = 0;
	_capacity = 0;
};
```

After make(), we create the backed array.
```
// Golang code
b := make([]int, 5, 10)
```

If in C++, it looks like
```
// C++ code
b = struct slice {
	_ptr = new int[10];   // the allocated 10-integer memory is the second layer, which is pointed to by _ptr
	_start = 0;
	_len = 5;
	_capacity = 10;
};
```

After b is initialized by make(), b[2:9] in Golang looks like
```
// Pseudo C++ code for imagination of b[2:9]
b[2:9] = struct slice {
	_ptr = b->_ptr; // the backed array does not change
	_start = 2; // start index is from 2
	_len = 7;	// 9-2, which means you can save 7 elements in the slice
	_capacity = 8;   // array size - _start, i.e. 10 - 2 = 8
};
```

Then 
```
// Golang code
b = b[2:9]  // means assign (copy) the b[2:9] struct to b
```

So b now looks like after b = b[2:9]
```
// Pseudo C++ code
b = struct {
	_ptr = the original backed array with size of 10
	_start = new position of 2 in the array as 0 index of the slice
	_len = ths size of slice, which is 7
	_capacity = the capacity of slice, which is 8
};
```

So after b = b[2:9], if we do
```
// in Golang
b[1] = 99
```
it equals to 
```
// C++ code
b->_ptr[2+1] = 99;
```

You can check the trick in the following Golang code.
```
b := make([]int, 5, 10)
b = b[2:9]
b = b[4:]
fmt.Println(len(b), cap(b))
```

The output is: len(b) = 3, cap(b) = 4.

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
		b[i] = i + 1	// write 1, 2, 3
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
		b[i] = i + 1	// write 1, 2, 3
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

Because in f(), the _ptr in s is pointed to another array, a totally new array, which is created by make().

After make(), s->_ptr in f() is different from b->_ptr, so the backed arrays are different.

In example one, the _ptr in s and b is same, i.e. the backed array does not change.

# Slice nil

Note:

var b int[] is totally different from b := []int{}

[Check nil for more details](nil.md)
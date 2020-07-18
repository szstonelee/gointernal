
# Prerequisite

[nil video From Google](https://www.youtube.com/watch?v=ynoY2xz-F8s)

[Interface Internal](interface.md)

# Which types can use nil

Only pointer, map, slice, function, channel, interface can use nil

```
var a int
frm.Println(a == nil) // illegal

type myStruct struct {}
var b myStruct
frm.Println(b == nil)  // illegal
```

# slice with nil

## slice internal

A slice variable has three items (fields), 

1. _ptr, internal pointer, point to the backed array 
2. _len, internal length, int size
3. _cap, internal capacity, int size

NOTE: _ptr is not a Golang pointer type which will be described below.

## When slice is nil

When a slice is constructed (but no assignment occur or not assigned with nil), 

every item is zero, which means _ptr == nullptr from C++'s view.

If the internal _ptr is nullptr, the slice is nil, which means the backed array does not exist. 

```
var a []int
fmt.Println(a == nil) // will print true
```

## After assignment of {}, it is different
```
a := []int{}
fmt.Println(a == nil) // will print false
```

Why? Becuase []int is different from []int{}

For []int{}, the backed array is allocated, though the size of the array is zero!

At this time, the internal _ptr is not nullptr, it has the memory address to the zero-sized array. 

## Think it as C++ or Java code

You can treat the above like the C++ code
```
_ptr = nullptr;   // when constructed

// when assigned with {}
int* buf = new int[0];
_ptr = buf;
assert(_ptr != nullptr);
```

Or From Java's view, it looks like
```
_ptr = null;  // when constructed

// when assigned with {}
_ptr = new Int[0];
assert _ptr != null;
```

## slice index out of bound

e.g. 1
```
var a []int
fmt.Println(a[0]) // will panic with index out of range
```
e.g. 2
```
a := []int{}
fmt.Println(a[0]) // will panic with index out of range
```

Example 1 and 2 trigger the same panic, but it is a little different.

The first one is like: because _ptr == 0, so panic

The second one is like: because_ptr->size() == 0, so panic

# map with nil

Map is similiar to slice. Check the following code

```
var a map[int]string
fmt.Println(a == nil) // will print true

a := map[int]string{}
fmt.Println(a == nil) // will print false
```

You can imagine there is _ptr in map. 

_ptr is like C++'s pointer, similar to slice.

The _ptr points to an allocated memory which is the real (or backed) hash map data structure.

When constructed but not assigned any value, _prt == nullptr, i.e. zero

When assigned an empty hash map of {}, the _ptr is not zero. It is the memory address of the backed empty hash map.

# pointer with nil

We can treat pointer in Golang similar to slice and map.

It means there is an internal _ptr in pointer.

e.g.
```
var p *[]int
fmt.Println(p == nil) // will print true
```
When construted, _ptr is nullptr. So if _ptr == nullptr, the pointer is nil.

```
var a []int // a is nil
p = &a
fmt.Println(p == nil) // will print false
```
because a has been constructed, the & opertator will assign the address of a to _ptr.

It means: _ptr != nullptr

```
fmt.Println(*p == nil) // will print true
```
because  *p == nil equals to a == nil

# interface with nil

## interface internal items

When it comes to interface, it is tricky.

You can imagine there are two interal items in interface

1. _ptr_to_type
2. _concrete_val

NOTE: _concrete_val can not be an interface.

When _ptr_to_type == nullptr && _concrete_val == nil, the interface is nil.

Otherwise, it is not nil, even the _concrete_val may be nil.

For assignment of interface, there are three modes:

## Assignment mode 1: i = ob

where ob is not interface, it means
```
_ptr_to_type = type of ob // remember Golang is strictly typed, so type exists
_concrete_val = ob  // So when ob is nil, it is OK
```

## Assignment mode 2: i = iOther

where iOther is another interfae, it means
```
i._ptr_to_type = iother._ptr_to_type
i._concrete_val = iother._concrete_val
```
## Assignment mode 3: i = nil
```
_ptr_to_type = nullptr
_concrete_val = nil
```

## Sample code

```
	var any interface{}
	fmt.Println(any == nil) // will print true

	var a []int = []int{1, 2, 3}
	fmt.Println(a == nil) // will print false

	a = nil
	fmt.Println(a == nil) // will print true

	any = a
	fmt.Println(any == nil) // will print false

	any = nil
	fmt.Println(any == nil) // will print true

	var i interface{ method() }
	fmt.Println(i == nil) // will print true

	any = i
	fmt.Println(any == nil) // will print true
```



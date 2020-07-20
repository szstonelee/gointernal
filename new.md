# Golang construction of object

In other language, like C++, you need the `new` operator to construct object in heap (if not optimized by compilor), 

or no new to construct the object in stack. 

```
// C++ code
class MyClass {
public:
  MyClass(int v) : val_(v) {}

private:
  int val_;
}

void func() {
  MyClass* p = new MyClass(100); // in heap

  MyClass a(200);  // in stack, after exit func(), a will be deleted

  delete p;
}
```

In Java, every object (not primitives) is allocated in heap by new (Except the compile optimization for escape)

Otherwise, the reference to the object in Java is null.

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

## literal

```
c := myStruct{} // equivalent to var c myStruct
```

The difference is that new returns pointer while var and literal return no pointer. You can do

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

For pointer, it is trivial and subtle

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
fmt.Println(*p6)  // will panic
```

## make() is for the underlying initialization for map, slice and channel

Because map, slice and channel has two layers. 

The top layer is a data structure for abstraction or logic description.

The underlying layer is the real data structure for the type. For slice, it is an array. For map, it is a hash map.

When map, slice, channel is constructed, it only has the top layer, but no underlying layer.

After initialization, like make() does or with literal assinment, the underlying layer is constructed.

```
var a []myStruct = make(myStruct, 5, 10)      // length = 5, capcaity = 10

var b []int
// will print true, the top layer exist, but the underlying layer, i.e. the int array, does not exist
fmt.Println(b == nil)  

var c []int = []int{} // it equals to c := []int{} which is idiomatic
// will print false, which equals to var c []int = make([int, 0, 0])
fmt.Println(c == nil) 
```

check [nil](nil.md) for more details


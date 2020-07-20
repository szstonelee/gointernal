# Function internal

# code
```
package main

import "fmt"

type myFunc func(int, []int) (string, error)

type myFunc2 func(int, []int) error

func main() {
	var a int = 1
	var b string = "I am a string!!!"

	f1 := func(ar1 int, ar2 []int) (string, error) {
		fmt.Println("closure sample", a, b) // closure
		fmt.Println("argument print", ar1, ar2)
		a = 9
		return "abc", nil
	}

	f1(2, []int{4, 5, 6})
	fmt.Printf("f1 type = %T, %T\n", f1, myFunc(f1))

	var any interface{} = f1
	if _, ok := any.(myFunc); ok {
		fmt.Println("f1 is myFunc")
	} else {
		fmt.Println("f1 is not myFunc")
	}

	var f2 myFunc = func(ar1 int, ar2 []int) (string, error) {
		fmt.Println("closure sample", a, b) // closure
		fmt.Println("argument print", ar1, ar2)
		fmt.Println("f2 print one more line....")
		return "abc", nil
	}
	fmt.Printf("f2 type = %T, %T\n", f2, myFunc(f2))

	any = f2
	if _, ok := any.(myFunc); ok {
		fmt.Println("f2 is myFunc")
	}

	// f3 := myFunc2(f2) // illegal
	f1 = f2
	any = f1
	if _, ok := any.(myFunc); ok {
		fmt.Println("f1 assign from f2, f1 is myFunc")
	} else {
		fmt.Println("f1 assign from f2, f1 is not myFunc though f2 is myFunc")
		fmt.Println("call f1() which is assigned from f2")
		f1(111, []int{333, 444, 555, 666})
	}

	// if f1 == f2 {	// illegal, type not match though signature is same
	// if any == f1 // illegal, interface can not be compared with func
}
```

# Run result

```
closure sample 1 I am a string!!!
argument print 2 [4 5 6]
f1 type = func(int, []int) (string, error), main.myFunc
f1 is not myFunc
f2 type = main.myFunc, main.myFunc
f2 is myFunc
f1 assign from f2, f1 is not myFunc though f2 is myFunc
call f1() which is assigned from f2
closure sample 9 I am a string!!!
argument print 111 [333 444 555 666]
f2 print one more line....
```

# Conclusion

## Support closure

```
f1 := func(ar1 int, ar2 []int) (string, error) {
	fmt.Println("closure sample", a, b) // closure, can use varirable out of func scope
```

Golang func has the attribute of closure.

It is like C++ & Java Lambda Function. 

But it is more like Python clousure, no need to declaration. Easy to use.

## Closure capture: passed by reference.

```
f1() {
	a = 9
}

closure sample 9 I am a string!!!
```

## Same signature, but different type

f1 and f2 have the same signature, but they are different types.

f1 type: func(int, []int) (string, error)

f2 type: main.myFunc

It is important for interface implementation.

## Type converstion 

### legal when signatures are same

```
fmt.Printf("f1 type = %T, %T\n", f1, myFunc(f1))	// will print main.myFunc
```

### illegal when signatures are not same, signature includes return

```
type myFunc func(int, []int) (string, error)

type myFunc2 func(int, []int) error

// f3 := myFunc2(f2)    // illegal
```

## Assignment 

Assginment is legal for same signature, but the type does not change.

```
f1 = f2
if _, ok := any.(myFunc); ok {
	...
} else {
	fmt.Println("f1 assign from f2, f1 is not myFunc though f2 is myFunc")
}
```

# Usage 1: func for interface in http

[Error handling and Go](https://blog.golang.org/error-handling-and-go)

Check: http.Handle()

I have different idea about the http.HandleFunc in the 'Error Handling and Go'.

Code could look like this

```
func init() {
	http.HandleFunc("/view", viewRecordForAll)
}

func viewRecordForAll(w http.ResponseWriter, r *http.Request) {
	if err := viewRecordWithoutError(w, r); err != nil {
		// deal with error in one place
		http.Error(w, err.Error(), 500)
	}
}

func viewRecordWithoutError(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)
	key := datastore.NewKey(c, "Record", r.FormValue("id"), 0, nil)
	record := new(Record)

	if err := datastore.Get(c, key, record); err != nil {
		return err
	}

	if err := viewTemplate.Execute(w, record); err != nil {
		return err
	}

	DoOtherThingsWhenErrorJustReturn()

	return nil
}
```

# Usage 2: func for dealing with error

[From the blog](https://blog.golang.org/errors-are-values)

## Problem: do not like so many if err != ni {}

If you do not like the error return, like the following code

```
_, err = fd.Write(p0[a:b])
if err != nil {
    return err
}
_, err = fd.Write(p1[c:d])
if err != nil {
    return err
}
_, err = fd.Write(p2[e:f])
if err != nil {
    return err
}
// and so on
```

## Solution 1: wrap with a func and a variable err

You can wrap it with a func
```
var err error
write := func(buf []byte) {
    if err != nil {
        return
    }
    _, err = w.Write(buf)
}
```

Then, the code will look like this
```
write(p0[a:b])
write(p1[c:d])
write(p2[e:f])
// and so on
if err != nil {
    return err
}
```

## Solution 2, wrap more, the err varble to struct and use interface

Even more, you can wrap the err variable to a struct, then make the struct implement an interface method

like this
```
type errWriter struct {
    w   io.Writer
    err error
}

func (ew *errWriter) write(buf []byte) {
    if ew.err != nil {
        return
    }
    _, ew.err = ew.w.Write(buf)
}
```

Then the code looks like this
```
ew := &errWriter{w: fd}
ew.write(p0[a:b])
ew.write(p1[c:d])
ew.write(p2[e:f])
// and so on
if ew.err != nil {
    return ew.err
}
```

## Comparison

```
	// if f1 == f2 {	// illegal, type not match though signature is same
	// if any == f1 // illegal, interface can not be compared with func
```

Note: function object can only compare to nil
```
	var f1 func(int) = func(int) {}
	var f2 func(int) = func(int) {}

	// fmt.Println(f1 == f2)	// illegal
```

## Function as receiver

Usually we use struct or pointer (which usually points to a struct) as a receiver.

We can use function type as a receiver too.

```
package main

import "fmt"

type myFunc func(string) string

type myInterface interface{ method() }

func (f myFunc) method() {
	fmt.Println("inteface mehtod() be called!")
}

func main() {
	var f myFunc = func(a string) string { return a + ":suffix" }
	a := f("abc")
	fmt.Println(a)
	var b myInterface = f
	b.method()
}
```

NOTE: interface can not be a receiver.
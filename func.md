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

f1 and f2 have the same signature, but they are different type.

f1 type: func(int, []int) (string, error)

f2 type: main.myFunc

It is important for Interface.

## Type converstion 

### legal when signatures are same

```
fmt.Printf("f1 type = %T, %T\n", f1, myFunc(f1))	// which print main.myFunc
```

### illegal when signatures are not same, signature includes return

```
type myFunc func(int, []int) (string, error)

type myFunc2 func(int, []int) error

// f3 := myFunc2(f2)    // illegal
```

## Assignment 

Assginment is legal for same signature, but the types do not change

```
f1 = f2
if _, ok := any.(myFunc); ok {
	...
} else {
	fmt.Println("f1 assign from f2, f1 is not myFunc though f2 is myFunc")
}
```

# One Usage

[Error handling and Go](https://blog.golang.org/error-handling-and-go)

http.Handle()

I have different idea about the http.HandleFunc in the 'Error Handling and Go', code could look like this

```
func init() {
	http.HandleFunc("/view", viewRecordForError)
}

func viewRecordForError(w http.ResponseWriter, r *http.Request) {
	if err := viewRecordForGood(w, r); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func viewRecordForGood(w http.ResponseWriter, r *http.Request) error {
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


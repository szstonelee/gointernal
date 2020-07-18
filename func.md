# Function internal

# code
```
package main

import "fmt"

type myFunc func(int, []int) (string, error)

type myFunc2 func() (string, error)

func main() {
	var a int = 111
	var b string = "I am a string!!!"
	f1 := func(ar1 int, ar2 []int) (string, error) {
		fmt.Println("closure sample", a, b) // closure
		fmt.Println("argument print", ar1, ar2)
		return "abc", nil
	}
	_, err := f1(11, []int{22, 33, 44})
	if err != nil {
		return
	}

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

	// f3 := myFunc2(f2)    // illegal
	f1 = f2
	any = f1
	if _, ok := any.(myFunc); ok {
		fmt.Println("f1 assign from f2, f1 is myFunc")
	} else {
		fmt.Println("f1 assign from f2, f1 is not myFunc though f2 is myFunc")
		fmt.Println("call f1() which is assigned from f2")
		f1(111, []int{222, 333, 444, 555})
	}
}
```

# Run result

```
closure sample 111 I am a string!!!
argument print 11 [22 33 44]
f1 type = func(int, []int) (string, error), main.myFunc
f1 is not myFunc
f2 type = main.myFunc, main.myFunc
f2 is myFunc
f1 assign from f2, f1 is not myFunc though f2 is myFunc
call f1() which is assigned from f2
closure sample 111 I am a string!!!
argument print 111 [222 333 444 555]
f2 print one more line....
```
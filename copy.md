Golang interanl use copy everywhere.

# Code one

```
package main

import "fmt"

type myStruct struct {
	age int
}

type changer interface {
	Change()
}

func (s myStruct) Change() {
	s.age = 99
}

func main() {
	b := myStruct{age: 20}
	fmt.Println(b)

	copy := b

	b.Change()
	fmt.Println(b)

	p := &copy
	p.Change()
	fmt.Println(p)
}

```

# Result one

```
{20}
{20}
&{20}
```

p.Change() actually is (\*p).Change(), then a copy of (\*p) to method Change() as argument for paramenter s.

So after p.Change(), actually no modification occur to p.

# Code two

```
package main

import "fmt"

type myStruct struct {
	age int
}

type changer interface {
	Change()
}

func (s *myStruct) Change() {
	s.age = 99
}

func main() {
	b := myStruct{age: 20}
	fmt.Println(b)

	c := b

	b.Change()
	fmt.Println(b)

	fmt.Println(c)
	c.Change()
	fmt.Println(c)
}
```

# Result two

```
{20}
{99}
{20}
{99}
```

c.Change() translates as (&c).Change(), then s in Change() is the copy value of &c. i.e. s holds the address of c.

So after Change(), c's age has been changed.

# Code three
```
package main

type myStruct struct {
	age int
}

func (s myStruct) Change() {
	s.age = 99
}

func returnMyStruct() myStruct {
	return myStruct{age: 30}
}

func main() {
	returnMyStruct().Change()
}
```

The above code is legal, though it outputs nothing.

But if we changed the receiver from s myStruct to s *myStruct, the above code is not legal.

Because the return from returnMySttruct() can not be addressed. Imagine the return of returnMySttruct() maybe exisit in the register of CPU. It is like the C++ rvalue.

For the legal version, the return of returnMySttruct() may be in the CPU register and can not be addressed, but for the receiver is not a pointer, a copy for the return CPU register value is created in memory (or created in another CPU register) and it is the s in Change() which can be used.
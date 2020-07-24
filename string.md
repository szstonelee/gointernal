
# string internal

Golang string is like Java String

## pass reference value to function

```
package main

import "fmt"

func f(s string) {
	fmt.Println(len(s))
}

func main() {
	a := "abc"
	f(a)
}
```

In the above example, if a is assigned with a big string, e.g. one million length string, when we go into fmt.Println(len(s)), there is no one more million char be allocated for s. a and s share the same underlying array of char (Acutually, in Golang, not char, but rune).

So it is different from C++
```
// C++ code
void f(std::string s) {
  int l = static_cast<int>(s.size());
}
```
From the above C++ code, unlike Golang, in f(), there are a new memory allocated for s. If a is string of one million char, when in f(), one million more memory is allocated for s. It is a big burden for C++ if pass the std::string for the parameter this way. 

So usually, like Java, we pass the string directly as the argument in Golang.

## string is immutable like Java
```
package main

func f(s string) {
	s += " :tail"
}

func main() {
	a := "abc"
	f(a)
	println(a) // will print "abc", not "abc :tail"
}
```

Similar to Java, string in Golang is immutable. For s += " :tail" in f(), it creates a new string, which is "abc :tail" for s. But because s is something like the reference of Java, it can not change the variable a. So in main(), print(a) will output "abc", not "abc :tail".

It is like the Java way
```
// Java code
void f(String s) {
  s += " :tail";  // s reference will point to a new string "abc :tail" which is created by an internal StringBuilder
}

void main() {
  String a = "abc";
  f(a);
  System.out.Println(a);    // will print "abc"
}
```

# We can use stirng the C++ way

But differnt from Java, in Golang, we can use pointer. This way, we can make a and s share the same underlying array of char.

```
package main

func f(s *string) {
	*s += " :tail"
}

func main() {
	a := "abc"
	f(&a)
	println(a) // will print "abc :tail"
}
```

It is similar to the C++ code
```
// C++ code
void f(std::string *s) {  // or using reference like void f(std::string& s) { s += " :tail";}
  *s += " :tail";
}

void main() {
  std::string a = "abc";
  f(&a)
  std::cout << a << std::endl;  // will print "abc :tail"
}
```

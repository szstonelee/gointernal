
# string internal

Golang string is almost like Java String.

## Pass reference value to function

```
func f(s string) {
	fmt.Println(len(s))
}

func main() {
	b := "aaa"	// suppose we assign a huge string of one million length like "aaa.............."
	f(b)
}
```

If b is assigned with a huge string of one million length in main(), when it comes to f(), there is no more million char allocated for s. 

b and s share the same underlying array of char.

NOTE: We use char. Acutually, in Golang, the correct name for char is rune. 

If you code in Java, the memory allocation is the same.

But it is different for C++.
```
// C++ code
void f(std::string s) {
  size_t l = s.size();
}
```
In the C++ code above, unlike Golang, there is a new more memory allocated for s. If b is a string of one million char, when in f(), one million more memory is allocated and copied from b for s. It is a big burden for C++ if passing the huge std::string directly as value. In the following section, we will see how C++ deals with the issue.

So usually, like Java, we pass string directly as argument in Golang.

## string is immutable like Java
```
func f(s string) {
	s += " :tail"
}

func main() {
	b := "abc"
	f(b)
	fmt.Println(b) // will print "abc", not "abc :tail"
}
```

Similar to Java, string in Golang is immutable. For s += " :tail" in f(), Golang creates a new string for s, which is "abc :tail" . It can not change the variable b. So in main(), Println(b) will output "abc", not "abc :tail".

It is similar to what the Java does
```
// Java code
void f(String s) {
  s += " :tail";  // s will reference to a new string "abc :tail" which is created by an internal StringBuilder
}

void main() {
  String b = "abc";
  f(b);	// b does not change
  System.out.println(b);    // will print "abc"
}
```

# We can use stirng the C++ way of pointer

But differnt from Java, in Golang, we can use pointer for string. 

By pointer, we can make b and s share the same underlying array of char.
```
func f(s *string) {
	*s += " :tail"
}

func main() {
	b := "abc"
	f(&b)
	println(b) 	// will print "abc :tail"
}
```

It is similar to C++, which can use pointer to solve the burden issue for huge string in the example above.
```
// C++ code
void f(std::string *s) {  // or using reference like void f(std::string& s) { s += " :tail";}
  *s += " :tail";	// no new string memory allocated for s
}

void main() {
  std::string b = "abc";
  f(&b);
  std::cout << b << std::endl;  // will print "abc :tail"
}
```

For string in Python, you can imagine Python works like Java.
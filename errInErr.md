# Sample

[Check the source](https://golang.org/pkg/fmt/#Fprintf), please clieck the example in the page

The code is copied as follow

```
package main

import (
	"fmt"
	"os"
)

func main() {
	const name, age = "Kim", 22
	n, err := fmt.Fprintf(os.Stdout, "%s is %d years old.\n", name, age)

	// The n and err return values from Fprintf are
	// those returned by the underlying io.Writer.
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fprintf: %v\n", err)
	}
	fmt.Printf("%d bytes written.\n", n)

}
```

# error dealing

In Golang, we need deal with error because it may occurs, i.e. err != nil.

In the above example, after we first call Fprintf(), the following code checks err. If err is not nil, we output the info to os.Stderr.

But the interesting things is, if Fprintf() in the error handling block returns error, what can we do?

If we go on to code for that, it is like the for-ever recursive game.

Someone would say, 

1. The first Fprintf() is for os.Stdout, but the next one is for os.Stderr, so we do not need check error anymore.

2. The second Fptintf() use %v, which can guarantee no fault. Sometimes, if no more info, %v will print the memory address.

But if something happen like:

For 1, we want debug the %flag for Stderr

For 2, what if the error is not related to %v and is related to the memory which is allocated by Stderr.

# My View

For some errors, we do not need to deal with them.

It is not only for Golang, it is for everything.

For example, in C, we malloc(), if it failed, sometimes we can tolerate it, but in most cases, we can not go on.

It is the same thing for *new* in C++.

In Java, we use exceptions. Java uses checked and unchecked exception as the solution. It is a good idea. For example, if OOM, which is an unchecked exception, what can we do? The best way is to let the application crash and show the OOM message, which is the work of Java runtime. So we can stop our application, try: 1. optimize our code, 2. tune JVM, 3. add more memory to machine.

But Java has some little issues with the checked exception. For the IO, Java treats them as checked exception. But IO has two kinds. One is for network, the other is for file. If it is for file, mostly, we can not go on. Even for the network, there are two sub situations. If it is LAN, the error is not tolerant in some situations. For example, if you commit a transaction for LAN database, but no response return because of a network error, can you retry? No!

In Linux, if you request more memory, but the kernel can not do it, Linux just kills some other applications to make room for your memory request, or Linux just kills you to stop the request. 

The solution is too simple but reasonable. Sometimes it is the best approach.

In C++ we think construtor can throw exception. but destructor should not throw exception, because construtor apply for new resourcre and destructor works for closing resource which should always succeed.

When you code local procedure call (LPC), you do not code like this
```
void caller() {
  try {
    int res = lpc_sum(1, 2);
  } catch (CPUException e) {
    print("CPU is too hot or we bought a fake CPU!!!")
  }
}

int lpc_sum(int a, int b) {
  return a+b;
}
```

But if for RPC, we need code like the above to deal with other exceptions, because we assume:
1. there would be network issues
2. remote server maybe crash

That is why RPC are much more complicated than LPC. 

But if we assume
1. LAN is robust like local CPU
2. We treat remote server crash like the one of local machine 

If the assumptions are true, or if we code in limited error, it will make the RPC coding easier. It is the same way Linux does with OOM.

# More interesting 

frm.Println() actually call Fprintf() agian, [check here](https://golang.org/src/fmt/print.go?s=5840:5903#L202)
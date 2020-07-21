# Sample

[Check the source](https://golang.org/pkg/fmt/#Fprintf), please clieck the example

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

In Golang, we need deal with error when it may occur.

In the above example, when we first call Fprintf(), it returns err

Then the code check err, and if it not nil, output the info to os.Stderr.

But the interesting things is, if Fprintf() in the error handling block retturns error, what can we do?

It is like the for-ever recursive game.

Someone would say, 

1. the first Fprintf() is for os.Stdout, but the next one is for os.Stderr, so we do not need check error anymore.

2. the second Fptintf() using %v, which could guarantee no fault. (Sometimes, it prints out the address)

Yes. They are explanations. 

But if something happen like:

for 1, we want debug the %flag for Stderr
for 2, if the error is not related to %v, it is related to the mememry which is used by Stderr.

# My View

For some error, we do not need to deal with the error.

It is not only for Golang, it is for everything.

For example, in C, we malloc() something, if it failed, sometimes we can tolerate it, but in most cases, we can not go on.

It is the samething for new in C++.

In Java, we use exception. Java uses checked and unchecked exception fot the solution. It is a good idea. For example, if OOM, what can we do? The best way is to let the application crash and show the OOM message. So we can 1. optimize our code, 2. tuning JVM, 3. add more memory.

But Java has something wrong with the checked exception. For the IO, Java treat them like checked exception. But IO has two kinds. One is for network, the other is for file. If it is for file, mostly, we can not go on. Even for the network, there are two sub situations. If it is LAN, the error is not tolerant. For example, if you commit a transaction for database, no response return with network error, can you retry? No!

In Linux, if the process apply for allocation of new memory, but OS can not do it, it just kill some application to make room for the request. Or Linux just kill the applying process. 

This is the simple but the right way.

C++ think construtor can throw exception. but destructor should not throw exception, because construtor usually apply for new resourcre, destructor works for return resource. Return should always success.

When you code local procedure call, you do not code like this
```
void caller() {
  try {
    int res = lpc_sum(1, 2);
  } catch(CPUException e) {
    ...
  }
}

int lpc_sum(int a, int b) {
  return a+b;
}
```

But if cod RPC, we need deal with the exception, because we assume:
1. there are network issue
2. remote server can be wrong

That is why RPC much more complicated than LPC. 

But if we assume
1. LAN is robust like local CPU
2. remote server is same strong as the local machine

It will make the RPC coding easier. Or we just treat it like the Linux way for OOM.

# More interesting 

frm.Println() actually call Fprintf() agian, [check here](https://golang.org/src/fmt/print.go?s=5840:5903#L202)
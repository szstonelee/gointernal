# Golang schedule internal

## an interesting article

I happened to read an interesting article about Golang schedule internal, which is titled as [A pitfall of golang scheduler](http://www.sarathlakshman.com/2016/06/15/pitfall-of-golang-scheduler).

```
package main
import "fmt"
import "time"
import "runtime"

func main() {
    var x int
    threads := runtime.GOMAXPROCS(0)
    for i := 0; i < threads; i++ {
        go func() {
            for { x++ }
        }()
    }
    time.Sleep(time.Second)
    fmt.Println("x =", x)
}
```

I decided to try the code.

I supposed the code in the article is run in Linux. I tested it in my virtual Linux (by [Multipass](https://github.com/canonical/multipass)) which is Ubuntu 20.04.1 LTS in my Mac host. The go runtime version is go1.14.5 linux/amd64.

But the result is different. In my Mac, which is 4 cpu core, I tried to run like this
```
GOMAXPROCS=3 go run x.go
```
and with other GOMAXPROCS numbers like 0, 1, 2, 3, 4, 5, 6, 7, 8, 9 ...

The results of all test case are the same: 

**All can return with x printed as 0**

I believe what the article describe is what the author did. So I am curious.

Then I tried the code in Mac with macOS Catalina (10.15.6). It is the same.

Then I tried another go version 1.12.9, the result in Linux is different. (I do not test in MacOS) 

For 1.12.9, the result is: 

**No return, No termination**

There are three results for three versions of Golang runtime.

So I decided to dive deeper.

## Environment

Linux = Ubuntu 20.04.1 LTS

Golang Version = 1.12.9 (we will come to 1.14.5 later)

[Go 1.12.9 install manual in Linux](https://www.linode.com/docs/development/go/install-go-on-ubuntu/)

## new code fo dive deeper
```
package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	tNum := runtime.GOMAXPROCS(0)
	fmt.Println("max user space thread for go runtime = ", tNum)

	gNum := 4
	fmt.Println("go routine number = ", gNum)

	var x int
	for i := 0; i < gNum; i++ {
		go func() {
			for {
				x++
			}
		}()
	}

	fmt.Println("I am here after generating all goroutines")

	time.Sleep(time.Second)

	fmt.Println("x =", x)
}
```



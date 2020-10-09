# Golang schedule

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

But the result is different. In my Mac, which is 4 cpu core, it never terminates. I tried to run like this
```
GOMAXPROCS=3 go run x.go
```
and with other GOMAXPROCS numbers like 0, 1, 2, 3, 4, 5, 6, 7, 8, 9 ...

The result is the same, all return with x printed as 0.

I believe what the article describe is what the author did. So I am curious.

Then I tried the code in Mac with macOS Catalina (10.15.6). It is the same.

But if you use another go version like 1.12.9, the result in Linux is different. 

**Golang evolves. The runtime implementations are different.**

So I decided to dive deeper.

## new code for new test

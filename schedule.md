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

I tried the attractive code.

I supposed the code in the article is run in Linux. I tested it in my virtual Linux (by [Multipass](https://github.com/canonical/multipass)) which is Ubuntu 20.04.1 LTS in my Mac host. The go runtime version is go1.14.5 linux/amd64.

But the result is different. In my Mac, which is 4 cpu core, I tried to run like this
```
GOMAXPROCS=3 go run x.go
```
and with other GOMAXPROCS numbers like 0, 1, 2, 3, 4, 5, 6, 7, 8, 9 ...

The results of all test cases for Go 1.14.5 are the same: 

**All can return with x printed as 0**

I believe what the article describes is what the author did. So I am curious.

Then I tried the code in Mac with macOS Catalina (10.15.6). It is the same, i.e. all return with x = 0.

Then I tried another go version 1.12.9, the result in Linux is different. (I do not test in MacOS) 

For 1.12.9, all results of test case are the same:

**No return, No termination, Infinite loop**

There are three kinds of result based on what Go runtime version you have.

So I decided to dive deeper.

## new code fo dive deeper

### code
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

	fmt.Println("I am here after launching all goroutines")

	time.Sleep(time.Second)

	fmt.Println("x =", x)
}
```

### comments for the modified test code

[tNum is the max number of user space thread which are running simultaneously.](https://stackoverflow.com/questions/39245660/number-of-threads-used-by-go-runtime) 

gNum is the number of go routines to launch at the same time. Each of go routines run an infinite loop.

main() will sleep for one second after the launch, and if it can return to OS, it will print the value of x.

And from [stackoverflow](https://stackoverflow.com/questions/53388154/is-the-main-function-runs-as-a-goroutine), main() can be treated as an additional Go routine for the above code.

### NOTE about tNum and what I imagine of go routines

From the above link, tNum is the max number of the concurrent running threads. 

It hints that maybe there are more threads than tNum if the additional threads are not running, e.g. sleeping.

One Go routine is not corresponded to one thread, i.e. go routine != thread.

A Go routine is like a task job. When we create a Go routine, it just add a new task entry to a total task queue. Then Go runtime dynamiclly determines how many threads needs to finish these tasks. Go runtime has a bounded limit which is the GOMAXPROCS. But GOMAXPROCS only limits the running threads. 

After the creation of the working threads, go runtime will create a working queue for each working thread. Then runtime will allocate some tasks to each thread working queue from the total task queue because there is a depth limit for each thread working queue.

Each length of working thread queue is dynamic because some tasks finish quickly, or some tasks will block. So a working thread can take a task from the total task queue or steal a job from other working thread queue.

If the task is blocked for I/O of disk, runtime will create a new thread for it. So the thread is not runnable and not be counted for GOMAXPROCS.

If the task is blocked for I/O of network, because Linux kernel can use epoll for netwroking, runtime does not create a new thread for the blocked task. But the task is not running, it is blocked, waiting for an event, just like a task calling sleep().

If the runtime has chance, it will check the epoll or timer event for the blocked task. It is like Mutex in Linux. In Linux some threads are waiting for the synchronized primitive of Mutex, and can be scheduled on it. Blocked threads are just a queue for mutex. The difference is that mutex checking and thread schedule are handled by kernel and in kernel space, but the timer/epoll events are checked by Go runtime and the schedule of go routine is in user space and handled by Go runtime. 

Wnen will the runtime have chance to check? Any Go system call like runtime.Gosched() will do that. 

But if your go routine run an infinite loop, there is no chance to call Go runtime system call. So the timer event or epoll will be ignored.

e.g. 1

We create 3 Go routines, and Go runtime may create 3 threads for the tasks. After each Go routine finish in each thread. The threads will be reclaimed or return to a thread pool to rest.

e.g. 2

We create 10 go routines, and GOMAXPROCS = 5. 

Go runtime may create five threads, each thread may have a task queue with 2 Go rouitnes in each queue.

e.g. 3

In the above example, one Go routine call read() from disk. Go runtime will create a new thread for the read() Go routine. The number of threads is 6, but it meets the requirement of GOMAXPROCS, which is 5.

e.g. 4

From e.g. 2, a Go routine call sleep(). No more thread will be created. The sleep() routine will be marked or may be moved to another queue for the same thread. If the thread has chance to call Go runtime sys call, the runtime will check the timer event. If the event is coming, the runtime will resume the sleeping task.

But the Golang timer is not the same as the kernel timer. For kernel timer, because it is based on hardware and Linux is preemptive, each thread will have a chance to run or be checked. For Go timer, if the thread can not call into the runtime, which is invoked by any Go system call, the timer event will be ignored. So Go routines are co-operative, similiar to Python aysc framework.

## Test Environment for Go 1.12.9

### How setup

Linux = Ubuntu 20.04.1 LTS

Golang Version = 1.12.9 (we will come to 1.14.5 later)

[Go 1.12.9 install manual in Linux](https://www.linode.com/docs/development/go/install-go-on-ubuntu/)

### test cases

| tNum | gNum | Result |
| :--: | :--: | :-- |
| 2 | 1 | x = 0 |
| 2 | 2 | No Return |
| 2 | 4 | No Return |
| 4 | 2 | x = 0 |
| 4 | 3 | x = 0 |
| 4 | 4 | No Return |
| 4 | 5 | No Return |
| 5 | 3 | x = 0 |
| 5 | 4 | x = 0 |
| 5 | 5 | No Return |
| 5 | 6 | No Return |

### analyze for 1.12.9

If gNum is less than tNum, no infinite-loop Go routine task will be distributed to the queue of the main thread. It guarantees that the thread running main(), which we can call main thread, only has one Go routine task which is main(). After the task, main(), sleeps for one second, the main thread has chance to call into runtime, and the runtime found the timer event is coming, resume main() task, and then exectute the implicit exit(), which leads the whole process to terminate. 

When return, the main() will print x, but x is cached in cpu core, e.g. L1 cache, so it is zero. It means other threads which call the infinite-loop to increase x for ever, does call memory barrier instruction to sync x with the main thread.

If gNum is not less than tNum, an infinite-loop Go routine task will be assigned to queue of the main() thread. After sleep() be called in main(), the main task is makred as sleeping and/or be moved to another queue in the main thread, then the main thread will run the assigned infinite-loop task. Because it is infinite loop, no chance to call into runtime, it will not terminate because no chance to call implicit exit() in main(). The result is **No Return**.

If we modify the code like this
```
                        for {
                                x+
                                runtime.Gosched()
                        }
```
You will find the the result for **No Return** will return. And x will be printed for some value.

## Test Environment for Go 1.14.5

### How setup

Linux = Ubuntu 20.04.1 LTS

Golang Version = 1.14.5 

```
curl -O 
```
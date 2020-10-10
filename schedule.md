# Golang schedule internal

## an interesting article

I happened to read an interesting article about Golang schedule internal, which is titled as [A pitfall of golang scheduler](http://www.sarathlakshman.com/2016/06/15/pitfall-of-golang-scheduler).

The following code is copied from the article.

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

I supposed the code in the article is run in Linux. I tested the above code in my virtual Linux (by [Multipass](https://github.com/canonical/multipass)) which is Ubuntu 20.04.1 LTS in my Mac host. The Go runtime version is 1.14.5 linux/amd64.

But the result is different. In my Mac, which is 4 cpu core, I tried to run 
```
GOMAXPROCS=3 go run x.go
```
and with other GOMAXPROCS numbers like 0, 1, 2, 4, 5, 6, 7, 8, 9 ...

The outputs of all test cases for Go 1.14.5 are the same: 

**All can return with x printed as 0**

I believe what the article describes is what the author did. So I am curious.

Then I tried the code with macOS Catalina (10.15.6). It is the same, i.e. all return with x = 0.

Then I tried another Go version 1.12.9 in Linux, it is different. (I do not test 1.12.9 in MacOS) 

For 1.12.9, all outputs of test case are the same:

**No return, No termination, Infinite loop**

There are three kinds of output based on what Go runtime version you have.

1. In the article, some return, some not return.
2. In 1.12.9, no return.
3. In 1.14.5, all return.

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

### Comments and note for the modified test code

[tNum is the max number of user space thread which are running simultaneously.](https://stackoverflow.com/questions/39245660/number-of-threads-used-by-go-runtime) 

gNum is the number of Goroutines to launch at the same time. Each of Goroutines runs an infinite loop.

main() will sleep for one second after the launch, and if it can go on, it will print the value of x and call exit() implicitly to return to OS.

In Go, [exit() will terminate the whole process](https://stackoverflow.com/questions/25518531/ok-to-exit-program-with-active-goroutine), so all threads in the process will stop. It is different from Java. In Java, the main() thread exit, but if the Java app is not a daemon type, JVM runtime will wait for all other threads to stop.

And from [stackoverflow](https://stackoverflow.com/questions/53388154/is-the-main-function-runs-as-a-goroutine), main() can be treated as an special and additional Goroutine for the above code.

So basiclly, we have gNum+1 Goroutines. And we want to know the result for the Goroutine of main, i.e. whether it returns to OS.

From the above link, tNum is the max number of the concurrent running threads. 

It hints that maybe there are more threads than tNum if the additional threads are not running, e.g. sleeping or blocked.

### What I imagine of Goroutine

[A good article about the Go scheduler is here](https://www.ardanlabs.com/blog/2018/08/scheduling-in-go-part2.html)

One Goroutine is not corresponded to one thread, i.e. Goroutine != OS thread.

In some articles, the Goroutine is named as Go thread or green thread, and OS thread is named as machine thread or real thread. Here in this article, for simiplicity, thread is always OS thread. We try to view Goroutine from the view of OS thread.

A Goroutine is like a task job. When we create a Goroutine, we just add a new task entry to a global task queue(GRQ). Then Go runtime dynamiclly determines how many threads are needed to finish these tasks. Go runtime has a bounded limit for the number of threads which is the GOMAXPROCS. But GOMAXPROCS only limits the running threads. 

So here Goroutine == task.

Threads in Go will be scheduled to each cpu core. For each core, there is a task queue(LRQ). Tasks in GRQ will be distributed to LRQ and each task from LPQ will be executed in a thread. 

Because 
1. A thread can execute a lot of tasks in one LRQ, so overhead of thread context switch is low.  
2. One task does not need to be finished first for next task to be scheduled, so the tasks in one LRQ are exececuted in one thread concurrently.
3. Task schedule overhead is 1 tenth of thread overhead. Overhead for thread switch is a couple of us. Overhead for task switch is hundreds of ns.
4. If a task will be blocked, Go will deal with the blocked task specially to make the cpu core available for a running thread. i.e. The thread for the core never block. 

These are the keys of Go Scheduler.

Each length of LRQ is dynamic because some tasks finish quickly, or some one will block. So when a thread is scheduled to an empty LRQ, the Go runtime can steal some tasks from other LRQ so the lengths of LRQ are balanced.

If the task is blocked for I/O of networking, Go runtime applies epoll/IO Completion for the task. So the network-blocked Goroutine will be taken care by the Net Poller. i.e. The network-blocked task will no be in the LRQ at the point of time. And the thread will go on with other tasks in the LRQ.

If the task is blocked for I/O of disk in Linux, the thread will be blocked. But Go runtime knows that, so a new thread will take over the current LPQ. The blocked thread with the blocked task will wait to finish, i.e. when to be unblocked. After the blocked disk call finishs, the Goroutine will return to the LPQ, and the unblocked thread can return to the thread pool. [There is a proposal for an improvement for this strategy if you want to dive deeper](http://pages.cs.wisc.edu/~riccardo/assets/diskio.pdf).

NOTE: For Windows, because disk I/O can be treated as network I/O by IO Completion, the schedule in Windows is the same for disk and network I/O.

For Timer, [referenced from the article - Illustrated Tales of Go Runtime Scheduler](https://medium.com/@ankur_anand/illustrated-tales-of-go-runtime-scheduler-74809ef6d19b), we can treat timer something similiar to network I/O. [You can dive deeper from the implementation of timer](https://blog.gopheracademy.com/advent-2016/go-timers/). 

If the runtime has chance, it will check the epoll or timer event for the blocked task. 

It is like Mutex in Linux. In Linux some threads are waiting for the synchronized primitive of Mutex, and can be scheduled on it. Blocked threads are just a queue for mutex. The difference is that mutex checking and thread scheduling are handled by kernel and in kernel space, but the timer/epoll events are checked by Go runtime and the scheduling of Go routine is in user space and handled by Go Scheduler. 

Wnen will the runtime have chance to check? Any Go system call like runtime.Gosched() will do that. 

But if your Go routine run an infinite loop, there is no chance. So the event of timer or epoll will be ignored.

e.g. 1

We create 3 Go routines, and Go runtime may create 3 threads for the tasks. After each Go routine finish in each thread. The threads will be reclaimed or return to a thread pool to rest.

e.g. 2

We create 8 Go routines, and GOMAXPROCS = 4. 

Go runtime may create four threads, each thread run each core with a task queue of length 2.

e.g. 3

In the above example, one Go routine call read() from disk. Go runtime will move out the thread with the Go routine which calls read(). 

There could be two strategies for scheduling.

First:

A new thread could be created to replace the move-out thread for the left Go routines in the LRQ. This time, the thread number is 5, but it meets GOMAXPROCS = 4 because the move-out thread is blocked, not running.

Second:

The other goroutine in the same LRQ can be moved to another LRQ and be executed by another thread for the LRQ. The number of thread is 4 in this case.

e.g. 4

From e.g. 2, a Go routine call sleep(). No more thread will be created. The sleep() Go routine will be moved to a special queue which is for the timer event and be taken care by the Go runtime. If there is a chance to call into Go, the runtime will check the timer event. If the event is coming, the runtime will resume the sleeping task.

But the Golang timer is not the same as the kernel timer. For kernel timer, because it is based on hardware and Linux is preemptive, each thread will have a chance to run or be checked. For Go timer, if no chance to call into the runtime, which is invoked by any Go system call, the timer event will be ignored. So Go routines are co-operative, similiar to Python aysc framework.

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

If gNum is less than tNum, no infinite-loop Go routine task will be distributed to the queue of the main thread. It guarantees that the thread running main(), which we can call it as main thread, only has one Go routine task which is main(). After the task of main(), which sleeps for one second, the main thread has chance to call into runtime, and the runtime found the timer event is coming. The runtime will resume main() task, and the main thread will exectute the implicit exit() in main(), which leads the whole process to terminate. 

When return, the main() will print x, but x is cached in cpu core, e.g. L1 cache, so it is zero. It means other threads which call the infinite-loop to increase x for ever, do not call memory barrier instruction to sync x with the main thread.

If gNum is not less than tNum, an infinite-loop Go routine task will be assigned to queue of the main() thread. After sleep() be called in main(), the main task is makred as sleeping and/or be moved to another queue in the main thread, then the main thread will run the assigned infinite-loop task. Because it is infinite loop, no chance to call into runtime for the main thread, it will not terminate because no chance to call implicit exit() in main(). The result is **No Return**.

If we modify the code like this
```
                        for {
                                x++
                                runtime.Gosched()
                        }
```
You will find the the result for **No Return** will be changed to **return**. And x will be printed for some value because the main thread calls x++ before Gosched(). 

## Test Environment for Go 1.14.5

### How setup

Linux = Ubuntu 20.04.1 LTS

Golang Version = 1.14.5 

```
curl -O 
```
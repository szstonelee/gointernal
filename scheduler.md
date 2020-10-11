# Golang Scheduler internal

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

## new code to dive deeper

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

### Comments and notes for the modified test code

[tNum is the max number of user space thread which are running simultaneously.](https://stackoverflow.com/questions/39245660/number-of-threads-used-by-go-runtime) 

It hints that maybe there are more threads than tNum if the additional threads are not running, e.g. sleeping or blocked.

gNum is the number of Goroutines to launch at the same time. Each of Goroutines runs an infinite loop.

main() will sleep for one second after the launch, and if it can go on, it will print the value of x and call exit() implicitly to return to OS.

In Go, [exit() will terminate the whole process](https://stackoverflow.com/questions/25518531/ok-to-exit-program-with-active-goroutine), so all threads in the process will stop. It is different from Java. In Java, the main() thread exit, but if the Java app is not a daemon type, JVM runtime will wait for all other threads to stop.

And from [stackoverflow](https://stackoverflow.com/questions/53388154/is-the-main-function-runs-as-a-goroutine), main() can be treated as an special and additional Goroutine for the above code.

So basiclly, we have gNum+1 Goroutines. And we want to know the result for the Goroutine of main, i.e. whether it returns to OS.

### What I imagine of Goroutine

[A good article, Scheduling In Go, Part II - Go Scheduler, is here.](https://www.ardanlabs.com/blog/2018/08/scheduling-in-go-part2.html)

One Goroutine is not corresponded to one thread, i.e. Goroutine != OS thread.

In some articles, the Goroutine is named as Go thread or green thread or lightweight thread, and OS thread is named as machine thread or real thread. Here in this article, for simiplicity, thread is always OS thread. We try to check Goroutine from the view of OS thread.

Thread is the basic unit for OS to run code.

A Goroutine is a task job. When we create a Goroutine (go func()...), we just add a new task entry to a global task queue(GRQ). Then Go runtime dynamiclly determines how many threads are needed to finish these tasks. Go runtime has a bounded limit for the number of threads which is the GOMAXPROCS. But GOMAXPROCS only limits the running threads. The threads which controlled by Go runtime and do system work, are not counted for GOMAXPROCS, e.g. GC threads. 

So here Goroutine == task.

Threads in Go will be scheduled to each cpu core to run. For each core, there is a task queue(LRQ). Tasks in GRQ will be distributed to LRQ. And each task from LPQ will be consumed in a thread which goes to a core. 

The following 4 points are the keys of Go Scheduler.

1. A thread can execute a lot of tasks in one LRQ, so overhead of thread context switch is low. We do not need to context switch thread for each task. 

2. One task does not need to be finished first for next task to be scheduled, so the tasks in one LRQ are exececuted concurrently for the same thread. This is different from the [Java ExectutorService](http://tutorials.jenkov.com/java-util-concurrent/executorservice.html). In Java, each task must be finished then consumer threads pick next one from the synchronized queue. How Go does that? [Because in each function call, there is a chance for Go runtime to switch Goroutines](https://golang.org/doc/go1.2#preemption). 

3. Task schedule overhead is 1 tenth of thread switch overhead. Avarage overhead for thread switch is a couple of microseconds. Avrage overhead for task switch is hundreds of nanoseconds.

4. If a task will be blocked, Go will deal with the blocked task specially to make the cpu core available for a running thread to run next task. i.e. A thread running in the core never block. Please read the following details.

Each length of LRQ is dynamic because some tasks finish quickly, or some one will block. When a thread is scheduled to an empty LRQ, the Go runtime can steal some tasks from other LRQ so the lengths of LRQ are balanced.

If the task is blocked for I/O of networking, Go runtime applies epoll/IOCP(IO Completion Port) for the task. So the network-blocked Goroutine will be taken care by the [Net Poller](https://morsmachine.dk/netpoller). i.e. The network-blocked task will be moved out of the LRQ at the point of time. And the current thread will go on with other tasks in the LRQ. The specific thread of Net Poller will deal with the blocking task.

If the task is blocked for I/O of disk in Linux, the thread will be blocked. But Go runtime knows that, so a new thread will take over the current LPQ. The blocked thread with the blocked task will wait to finish, i.e. until unblocked. After the blocked disk call finishs, the Goroutine will return to the LPQ, and the unused thread can return to the thread pool or be destroyed. [There is a proposal for an improvement for this strategy if you want to dive deeper](http://pages.cs.wisc.edu/~riccardo/assets/diskio.pdf).

Again, a thread running in the core never block in Go.

The I/O for disk is so special in Linux because regular file descriptor are always blocked device.

1. [NON_BLOCKED has no effect for regular file descriptor](https://www.remlab.net/op/nonblock.shtml).
2. When poll for ready status of read and write of a regular file descriptor, it always return True, i.e. always ready.

The internal reason is related to the page cache. Even a page has been read from a disk file and located in the memory of page cache, the following read() could be blocked because the page could be evicted from cache just then. For write, even there are enough memory for it, the page in memory could need to write-back to disk. Page cache is unpredicatable because it is a cache. That is why Linux AIO needs DIRECT_IO, no buffered-IO.

NOTE: For Windows, because disk I/O can be treated as network I/O by IO Completion Port, the schedule in Windows is the same for disk and network I/O.

For Timer, [referenced from the article - Illustrated Tales of Go Runtime Scheduler](https://medium.com/@ankur_anand/illustrated-tales-of-go-runtime-scheduler-74809ef6d19b), we can treat timer something similiar to network I/O. [You can dive deeper from the implementation of Go timer](https://blog.gopheracademy.com/advent-2016/go-timers/). 

If the runtime has chance, it will re-schedule thre returned un-blocked task which are tiggered by timer, network I/O, regular file I/O. 

When will the runtime have chance to check? Any Go system call like runtime.Gosched() will do that. 

But if your Go routine run an infinite loop, there is no chance for Go Scheduler to take effect. So the returned tasks will be ignored.

e.g. 1

We create 3 Go routines, and Go runtime may create 3 threads for the tasks. After each Go routine finish in each thread. The threads will be reclaimed or return to a thread pool to rest.

e.g. 2

We create 8 Go routines, and GOMAXPROCS = 4. 

Go runtime may create four threads, each thread run each core with a task queue of length 2.

e.g. 3

In the above example, one Go routine call read() from disk. Go runtime will move out the thread with the Goroutine which calls read(). 

There could be two strategies for scheduling.

First:

A new thread could be created to replace the blocked thread for the other Goroutines in the LRQ. This time, the thread number is 5, but it meets GOMAXPROCS = 4 because the replaced thread is blocked, not running.

Second:

The other Goroutine in the same LRQ can be moved to another LRQ and be taken care by another thread. The number of thread is 4 in this case.

e.g. 4

From e.g. 2, a Go routine call sleep(). No more thread will be created. The sleep() Goroutine will be moved to a special queue which is for the timer event and be taken care by the Go runtime. 

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

### Analyze for 1.12.9

### For the above table

From the above table, we can conclude that

1. if tNum > gNum, return with x printed as 0
2. otherwise, no return

My Mac have four cores, how many LRQ we have? It equals to min(tNUm, 4). E.g. if tNum == 2, it is 2. if tNum == 5, it is 4.

There are gNum+1 Goroutines, the plus one Goroutine is main() itself.

Then gNum+1 Goroutines will be distributed to all LRQs.

The main Goroutine will run first. Because main call sleep for one second, the main Goroutine will be moved out of the LRQ. One second later, main Goroutine will come back to one LRQ. And the main Goroutine state is runnable at this point of time.

In the duration of one second, each thread of gNum threads will have a chance to run in one core because Linux is preemptive. Each thread needs to pick a not-running Goroutine from one LRQ and then run an infinite loop. So each Goroutine except the main Goroutine will change state from runnable to running. 

If tNum > gNum, there is at least one thread which state is not running, and when the non-running thread come to run, it needs pick a not-running Goroutine which is the come-back main Goroutine. At this point of time, it will call exit() implicitly and return to OS. 

x is printed as 0, because x is cached for each thread. For example, x is located in the register of a cpu core, and when a thread context switch, the register values are saved in the stack of the thread.

I do not think x is in L1 cache. Because there are thread switch, if x is in L1 cache, it will be flushed to main memory with some value. Then when the thread come to run main, it will get the updated value from main memory. It could not be zero in this situation.

If tNum <= gNum, no thread is available for the main Goroutine, because each thread is busy running an infinite loop. There is no chance to call into Go runtime which can schedule for the runnable main Goroutine.

### For adding runtime.Gosched()

If we modify the code like this
```
                        for {
                                x++
                                runtime.Gosched()
                        }
```

You will find the the result for **No Return** will be changed to **return**. And x will be printed for some value other than zero for tNum > gNum and tNum <= gNum.

For any condition, tNum > gNum or tNum <= gNum, because scheduler has chance to take effect by calling in untime.Gosched(), the main Goroutine has chance to be run. Note, even every thread is running an infinite loop, if Go Scheduler comes to play, every Goroutines has chance to run in the same thread, i.e. Goroutines are concurrent like threads in preemptive Linux OS.

That is why x is not zero. Because main Goroutine run after a Goroutine which run an infinite loop to increment x in the same thread, main Goroutine will see the updated x. 

## Test Environment for Go 1.14.5

### How setup

Linux = Ubuntu 20.04.1 LTS

Golang Version = 1.14.5 

```
curl -O https://dl.google.com/go/go1.14.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.14.5.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
go version
```

### test cases

| tNum | gNum | Result |
| :--: | :--: | :-- |
| 2 | 1 | x = 0 |
| 2 | 2 | x = 0 |
| 2 | 4 | x = 0 |
| 4 | 2 | x = 0 |
| 4 | 3 | x = 0 |
| 4 | 4 | x = 0 |
| 4 | 5 | x = 0 |
| 5 | 3 | x = 0 |
| 5 | 4 | x = 0 |
| 5 | 5 | x = 0 |
| 5 | 6 | x = 0 |
| 2 | 7999 | x = 0, sometimes 1-2 seconds duration, sometimes more than 70 seconds |
| 2 | 8000 | x = 0, sometimes 1-2 seconds duration, sometimes more than 80 seconds |
| 2 | 8001 | x = 0, sometimes 1-2 seconds duration, sometimes more than 70 seconds |

### Analyze for 1.14.5

From the above table, we can guess Go 1.14.5 adds new feature of preemption, e.g. adding a new system thread which can monitor all states of Goroutines in all working thread.

The following articles demonstrate this.

[For Go 1.13, add preemption](https://medium.com/a-journey-with-go/go-goroutine-and-preemption-d6bc2aa2f4b7).
[For Go 1.14, preemption improved as async preemption](https://medium.com/a-journey-with-go/go-asynchronous-preemption-b5194227371c).

The interesting stuff is:

1. How the schedule algorithm does, which leads some times, the main Goroutine return to OS very early, some times, the main Goroutine get chance to run very late.
2. x are always to be printed as 0. It seems when preemption takes effects, the context of main Goroutine is isolated. The main Goroutine must run after another Goroutine which increment x for ever, but the x for the two Goroutines are different objects, even they are in the same thread.

The above is two new problems to solve in future.
# Golang Scheduler internal

## An interesting article - A pitfall of golang scheduler

I happened to read an interesting article about Golang schedule internal, which is [A pitfall of golang scheduler](http://www.sarathlakshman.com/2016/06/15/pitfall-of-golang-scheduler).

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

From the author's pratice, if the number of Goroutines is one less than the number of cpu core of the machine by changing the code to 
```threads := runtime.GOMAXPROCS(0)-1```, x will be printed as 0. If the number of Gorouines is equal to the numbrer of cpu core, the program does not terminate, i.e. no return with x not be printed.

NOTE: for MAC OS, the number of cpu core reported by Go is virtual which is the double of real cpu core number, e.g., In Mac OS, if your machine has cpu core of 4, runtime.GOMAXPROCS(0) returns 8. But if you have virtual Linux in MAC, runtime.GOMAXPROCS(0) returns 4.

I supposed the code in the article is run in Linux. I tested the above code in my virtual Linux (by [Multipass](https://github.com/canonical/multipass)) which is Ubuntu 20.04.1 LTS in my Mac host. The Go runtime version is 1.14.5 linux/amd64.

In my Mac, which is 4 cpu core, I tried to run 

```
GOMAXPROCS=3 go run x.go
```
and with other GOMAXPROCS numbers like 0, 1, 2, 4, 5, 6, 7, 8, 9 ...

The outputs of all test cases for Go 1.14.5 are the same: 

**All can return with x printed as 0**

I believe what the article describes is what the author did. So I am curious.

Then I tried the code with macOS Catalina (10.15.6). It is the same, i.e. all return with x = 0.

Then I tried another Go version 1.12.9 in Linux, it is different. It includes the same result as the author said.

(I do not test 1.12.9 in MacOS) 

There are two kinds of output based on what Go runtime version you have.

1. In 1.12.9, some return, some no return.
2. In 1.14.5, all return.

That means Go runtime changed algorithm of Go Scheduler from 1.12 to 1.14.

So I decided to dive deeper.

## New code to dive deeper

### Code
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

### Comments and notes for the above modified code

[tNum is the max number of user space thread which are running simultaneously.](https://stackoverflow.com/questions/39245660/number-of-threads-used-by-go-runtime) 

It hints that maybe there are more threads than tNum if the additional threads are not running, e.g. sleeping or blocked.

tNum can easily be changed when run application in CLI. For example, set GOMAXPROCS to 9, type
```
GOMAXPROCS=9 go run yourCode.go
```

gNum is the number of Goroutines to launch at the same time. Each of Goroutines runs an infinite loop.

In the above modified code, I set gNum = 4. You can change it to any number you like for all test cases.

main() will sleep for one second after the launch of gNum Goroutines. If it can go on, it will print the value of x and call exit() implicitly to return to OS.

In Go, [exit() will terminate the whole process](https://stackoverflow.com/questions/25518531/ok-to-exit-program-with-active-goroutine). So all threads in the process will stop. It is different from Java. In Java, the main() thread exit, but if the Java app is not a daemon type, JVM runtime will wait for all other threads to stop.

And from [stackoverflow](https://stackoverflow.com/questions/53388154/is-the-main-function-runs-as-a-goroutine), main() can be treated as a special and additional Goroutine.

So basiclly, we have gNum+1 Goroutines. And we want to know the result for the Goroutine of main, i.e. whether it returns to OS.

### What I imagine of Goroutine

[A good article, Scheduling In Go, Part II - Go Scheduler, is here.](https://www.ardanlabs.com/blog/2018/08/scheduling-in-go-part2.html)

One Goroutine is not corresponded to one thread, i.e. Goroutine != OS thread.

In some articles, the Goroutine is named as Go thread or green thread or lightweight thread, and OS thread is named as machine thread or kernel thread or real thread. Here in this article, for simplicity, thread is always OS thread. We try to check Goroutine from the view of OS thread.

Thread is the basic unit for OS to run code.

A Goroutine is a task job. 

When we create a Goroutine like ```go func()...```, we just add a new task entry to a global task queue(**GRQ**). Then Go runtime dynamiclly determines how many threads are needed to finish these tasks. Go runtime has a bounded limit for the number of threads which is the GOMAXPROCS. But GOMAXPROCS only limits the running threads. The threads which are controlled by Go runtime, run Go system code and do system work, are not counted for GOMAXPROCS, e.g. GC threads. 

So here Goroutine == task.

Threads in Go will be scheduled to each cpu core to run. For each core (assuming all core are used), there is a task queue(**LRQ**). Tasks in GRQ will be distributed to LRQ. And each task from LRQ will be run in a thread which goes to a core. 

The following 4 points are the keys of Go Scheduler.

1. A thread can execute a lot of tasks in one LRQ without a thread context switch. So overhead of thread context switch is low. We do not need to context switch thread for each task. 

2. One task does not need to be finished first for next task to be scheduled. So the tasks in one LRQ are run concurrently (NOTE: not paralelly) for the same thread. This is different from the [Java ExectutorService](http://tutorials.jenkov.com/java-util-concurrent/executorservice.html). In Java, one task must be finished, then a consumer thread can pick next one from the synchronized queue. In Java, task execution are one by one for one thread. But in Go, task execution are concurrent for one thread. That is why Goroutine in Go is like thread in OS and sometimes is called as green thread. How does Go achieve that? [Because in each function call, there is a chance for Go runtime to switch Goroutines](https://golang.org/doc/go1.2#preemption). 

3. Task schedule overhead is 1 tenth of thread switch overhead. Average overhead for thread switch is a couple of microseconds. Average overhead for task switch is hundreds of nanoseconds.

4. If a task will be blocked, Go will deal with the blocked task specially to make the cpu core available for a running thread to run next task. i.e. **A thread running in the core (for current LRQ) never block**. Please read the following details.

Each length of LRQ is dynamic because some tasks finish quickly, some one will block. When a thread is scheduled to an empty LRQ, the Go runtime can steal some tasks from other LRQ. So the lengths of LRQ are balanced.

If the task is blocked for I/O of networking, Go runtime applies epoll/IOCP(IO Completion Port) for the task. So the network-blocked Goroutine will be taken care by the [Net Poller](https://morsmachine.dk/netpoller). The network-blocked task will be moved out of the LRQ at the point of time. And the current thread will go on with other tasks in the LRQ. The specific thread of Net Poller will deal with the blocking task.

If the task is blocked for I/O of disk in Linux, the thread will be blocked. But Go runtime knows that. So a new thread (probabally an idle thread) will take over the current LRQ. The blocked thread with the blocked task will wait to finish, i.e. until unblocked. After the blocked disk call finishs, the Goroutine will return to the LRQ, and the unused thread can return to the idle thread pool or be destroyed (usually no destroy). [There is a proposal for an improvement for this strategy if you want to dive deeper](http://pages.cs.wisc.edu/~riccardo/assets/diskio.pdf).

Again, a thread running in a cpu core, for the current LRQ, never block in Go.

The I/O for disk is so special in Linux because regular file descriptor are always blocked device.

1. [NON_BLOCKED has no effect for regular file descriptor](https://www.remlab.net/op/nonblock.shtml).
2. When poll for ready status of read and write of a regular file descriptor, it always return True, i.e. always ready.

The internal reason is related to the page cache. Even a page has been read from a disk file and located in the memory of page cache, the following read() could be blocked because the page could be evicted from cache just then. For write, even there are enough memory for it, the page in memory could need to write-back to disk. Page cache is unpredicatable because it is a cache. That is why Linux AIO needs DIRECT_IO, no buffered-IO.

NOTE: For Windows, because disk I/O can be treated as network I/O by IO Completion Port, the schedule in Windows is the same for disk and network I/O.

For timer, [referenced from the article - Illustrated Tales of Go Runtime Scheduler](https://medium.com/@ankur_anand/illustrated-tales-of-go-runtime-scheduler-74809ef6d19b), we can treat timer something similiar to network I/O. [You can dive deeper from the implementation of Go timer](https://blog.gopheracademy.com/advent-2016/go-timers/). 

If the runtime has chance, it will re-schedule the returned un-blocked task which has been triggered by timer, network I/O, regular file I/O. 

When will the runtime have chance to check? Any Go system call like runtime.Gosched() will do that. 

But if your Goroutine run an infinite loop, there is no chance for Go Scheduler to take effect. So the returned tasks will be ignored.

e.g. 1

Assuming a machine has 4 cpu cores and we set GOMAXPROCS = 3.

We create 3 Goroutines, and Go runtime may create 3 user-code threads for the tasks. After each Goroutine finish in each thread. The threads will be reclaimed or return to a thread pool to rest.

e.g. 2

Assuming a machine have 4 cpu cores, and we create 8 Goroutines, and GOMAXPROCS = 4. 

Go runtime may create four user-code threads, each thread run each core with a task queue of length 2.

e.g. 3

Assuming a machine have 4 cpu cores and we set GOMAXPROCS = 5.

There are four LRQs for the cores.

If Go runtime create 5 user-code threads, each thread has a chance to run in one/another core for the corresponded LRQ.

e.g. 4

Supposed from example 2, one Goroutine call read() from disk. Go runtime will move out the thread with the Goroutine which calls read(). 

There could be two strategies for scheduling.

First:

A new user-code thread could be created to replace the blocked thread for the other Goroutines in the LRQ. This time, the user-code thread number is 5, but it meets GOMAXPROCS = 4 because the replaced thread is blocked, i.e. not running.

Second:

The other Goroutine in the same LRQ can be moved to another LRQ and be taken care by another thread. The number of user-code thread is 4 in this case, i.e., one is blocked, the other 3 are running.

e.g. 5

Supposed from e.g. 2, a Goroutine call sleep(). No more user-code thread will be created. The sleep() Goroutine will be moved out to a special queue (accurately, a min-heap data structure for timer events) which is for the timer event and be taken care by the Go runtime. 

## Test Environment for Go 1.12.9

### How setup

Linux = Ubuntu 20.04.1 LTS

Golang Version = 1.12.9 (we will come to 1.14.5 later)

[Go 1.12.9 install manual in Linux](https://www.linode.com/docs/development/go/install-go-on-ubuntu/)

### Test cases

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

### Analysis for 1.12.9

#### For the above table

From the above table, we can conclude that

1. if tNum > gNum, return with x printed as 0
2. otherwise, no return

My Mac have four cores, how many LRQ do we have? It equals to min(tNUm, 4). For example, if tNum == 2, it is 2; if tNum == 5, it is 4.

There are gNum+1 Goroutines, the plus one Goroutine is main() itself.

Then gNum+1 Goroutines will be distributed to all LRQs.

The main Goroutine will run first. Because main call sleep() for one second, the main Goroutine will be moved out of the LRQ. One second later, main Goroutine will come back to one LRQ. And the main Goroutine state is runnable at this point of time.

In the duration of the one second, each thread of tNum threads will have a chance to run in one cpu core because Linux is preemptive. Each thread needs to pick a runnable Goroutine from one LRQ. All but main Goroutines run an infinite loop. So each Goroutine except the main Goroutine will change state from runnable to running if there are enough threads. And the states of corresponded thread are also running. 

If tNum > gNum, there is at least one thread which state is not running, and when the runnable thread come to run, it needs to pick a runnable Goroutine which is the come-back main Goroutine. At this point of time, it will call exit() implicitly and return to OS. 

x is printed as 0, because x is cached for each thread. For example, x is located in the register of a cpu core, and when a thread context switch, the register values are saved in the stack of the thread.

I do not think x is in L1 cache. Because there are thread switch, if x is in L1 cache, it will be flushed to main memory with some value. Then when the thread come to run main Goroutine, it will get the updated value from main memory. It could not be zero in this situation.

If tNum <= gNum, no runnable thread is available. Every thread is busy running. Although a runnable main Goroutine is in one LRQ, because each thread is busy running an infinite loop, there is no chance to call into Go runtime which can schedule for the runnable main Goroutine in the same thread. In other words, cuncurrency for Goroutine is disabled.

#### For adding runtime.Gosched()

If we modify the code like this
```
                        for {
                                x++
                                runtime.Gosched()
                        }
```

You will find the the result for **No Return** will be changed to **return**. And x will be printed for some value other than zero for tNum > gNum and tNum <= gNum.

For any condition, tNum > gNum or tNum <= gNum, because scheduler has chance to take effect by calling in runtime.Gosched(), the main Goroutine has chance to be run. Note, even every thread is running an infinite loop, if Go Scheduler comes to play, every Goroutines has chance to run in the same thread, i.e. Goroutines are concurrent for one thread like threads in preemptive Linux OS.

That is why x is not zero. Because main Goroutine run concurrently with a Goroutine which increments x in the same thread, main Goroutine will see the updated x. 

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

### Test cases

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

### Analysis for 1.14.5

From the above table, we can guess Go 1.14.5 adds a new feature like preemption. For example, Go runtime can add a new system thread which can monitor all states of Goroutines and all states of user-code threads, and do preemtion.

The following articles demonstrate this guess.

[For Go 1.13, add preemption](https://medium.com/a-journey-with-go/go-goroutine-and-preemption-d6bc2aa2f4b7).

[For Go 1.14, preemption improved as async preemption](https://medium.com/a-journey-with-go/go-asynchronous-preemption-b5194227371c).

The interesting stuff is:

1. How the schedule algorithm does, which leads some times, the main Goroutine return to OS very early, some times, the main Goroutine get chance to run very late.
2. x are always to be printed as 0. It seems when preemption takes effects, the context of main Goroutine is isolated. The main Goroutine must run after another Goroutine which increment x for ever, but the x for the two Goroutines are different objects, even they are in the same thread.

The two myths need to be answered in future.
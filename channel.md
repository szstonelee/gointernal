# channel

Golang channel is like slice and map. There are two layers data structure for them. 

[Check slice and map for two layers](new.md) and [nil help](nil.md) 

But channel has more states and more member fields. So channel is more complicated.

## Analog to Java

Golang channel is something analogous to Java Executor Framework. 

In the Java concurrency framework, there are three components:

1. the blocking queue, which has the task elements and has mutex to secure the concurrency
2. the producer, which enque one task into the blocking queue, usually one producer run in one thread
3. the consumer, which deque one task from the blocking queue, usually one consumer run in one thread

If the blocking queue is full, the producer, which is run in one thread, will blcok, i.e. be forced to sleep (or wait) by kernel. 

If nothing in the blocking queue, the consumer, will block.

When some producers are blocked, they are queued by kernel as a kernel queue of waiting threads for one mutex.

The same idea applies to blocked consuemers.

After the bloking queue changes from being full to having room, the framework wakes up one producer. Internally, it is an thread waken up by the kernel. So the producer can enque() its task to the blocking queue.

Java Exectutor Framework uses the kernel primitives like mutex to work the way described above. It has the overhead for the context switch from kernel space to user space. You can think one switch is around micro second.

Goang channel basically work in the user space, so it has much less cost. You can thinke if working always in the user space, the cost is about nano seconds (a couple, tens, hundreds? I do not know, but it is cheaper than micro second). 

It works like what the kernel does, so channel has the similar data structure as Java Execturor Framework.

## similar blocking queue

The Golang channel has an internal array like the blocking queue. Actually the interanl array is working as ring queue. If senders send something, the task stuff will be queued in the internal array.

Note, the internal array size can be zero. It means when a sender invokes enque(), it will block until a receiver invokes deque(). Or it will succeed at once if there has been a receiver called deque() before the enque(). 

## similar kernel thread queue

Like producer and consumer thread queues in kernel, Golang channel has two queues, one for the sender, the other for the receiver. If sender/receiver is blocked, it will be queued in the channel sender/reeiver queue.

## chanel state

### nil

When channel first constructed, it is nil. 

```
var ch chan int // ch is a channel with int, and ch is nil
```

After make() initialized, it has internal array which could be an zero-sized array and will be non-nil.

```
c1 := make(chan int, 10)  // c1 is channel with int, with room of 10 in the internal array
c2 := make(chan int)  // c2 is channel with int, the internal array is zero-sized, c2 != nil
var c3 chan int   // c3 is nil
c3 = make(chan int) // c3 now is not nil
```

If channel is nil, enque() and deque() will block

```
package main

import "fmt"

func main() {
	var ch chan int // ch is nil

	ch <- 5

	fmt.Println("can not print anything because ch <- 5 blocks.")	// will deadlock
}
```

### close

Golang channel has a special state, close. 

Note, close is effective when channel != nil.

At first, the channel's close state is false. After the sender call close(), it will be true.

When the channel state close == true, sender can not enque(). Otherwise, it will panic.

Whenever close is true or false, receiver can always deque().

if close == true and no element in the interan array, it will return an empty value to the receiver as quick as possible. 

NOTE: if the receiver is in range loop, it will ternimate the loop and will not get into the loop scope.

If you want to know the return value of a receiver is from the close state or just an empty value sent by a sender, you can test it like
```
// if v is int, we can not tell whether it is from the closed ch or just a zero int from the sender
// but we can use ok to tell it
v, ok := <-ch  
```

# channel with goroutine

We can use channel without goroutine, like 
```
package main

import "fmt"

func main() {
	ch := make(chan int, 1) // if ch := make(chan int), it will block and eventually deadlock

	ch <- 5

	a := <-ch

	fmt.Println(a) // will print 5
}
```

But if you want go run like thread switch, you need goroutine

In the above code, we can not make the ch as zero-size, it will deadlock,

but if we use goroutine, it will succeed
```
package main

import "fmt"

func gofunc(ch chan int) {
	ch <- 5
}

func main() {
	ch := make(chan int)

	go gofunc(ch)

	a := <-ch

	fmt.Println(a) // will print 5
}
```

Why? Because when execute go func(), it is like the main thread in main() create a new thread to do gofunc(). After that, the main thread will syncronize in the channel ch. The working thread, which execute the gofunc(), use the channel as an communication tool to sync the result. The result is returned from the channel, which is 5, to the variable a in main thread. 

Acutally, there is no working thread created, ervery thing is executed in one thread. It is the CSP feature of Golang runtime. [You can check the internal implementation](https://morsmachine.dk/go-scheduler). But you can imagine of working threads which maybe help you understand the new feature of Golang.

The is the philosophy of Golang:
```
Do not communicate by sharing memory; instead, share memory by communicating.
```

[Python asyncio framework has the similar idea.](https://realpython.com/async-io-python/)

# The number of goroutine is not the number of go exec

Do not forget, the main() is also a routine

```
package main

import "fmt"

func main() {
	ch := make(chan int, 2)
	ch <- 1
	ch <- 2
	ch <- 3
	fmt.Println(<-ch)
	fmt.Println(<-ch)
}
```

It will output: all goroutines are asleep - deadlock!

When the ch <-3 execute, the main() is a goroutine which is in the channel ch recevier queue.

# example for select
[From the Go tour about select](https://tour.golang.org/concurrency/5)
```
package main

import (
	"fmt"
	"time"
)

func fibonacci(c, quit chan int) {
	x, y := 0, 1
	for {
		select {
		case c <- x:
			x, y = y, x+y
		case <-quit:
			fmt.Println("quit")
			return
		}
	}
}

func main() {
	c := make(chan int)
	quit := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(<-c)
			time.Sleep(time.Second * 1)
		}
		fmt.Println("exit from here")
	}()
	go func() {
		time.Sleep(time.Second * 7)   // if more than 9 seconds, will print "exit from here", otherwise no error, try it!!!
		quit <- 1
	}()
	fibonacci(c, quit)
	time.Sleep(time.Second * 3)
}
```

You can use the above info to analyze the code, it is fun.

# use goroutine for parallelism

Goroutine with channel can be used for the parallel execution of multi cpu core. The following is an example.

Check sameBST.go, which is origiated from [A tour of go: Excercise: Equivalent Binary Tree](https://tour.golang.org/concurrency/8)

```
go run sameBST.go
```

It uses two channel and two go routine for parallelism.

Note: we use buffered channel (1K) for speed. Otherwise, the most portion of time is taken by the main memory seeking. You can try no buffered channel.

Check sameBST.cc, which is similar to sameBST.go but can only use one core.
```
g++ -std=c++17 sameBST.cc
./a.out
```

From the test results, you can see:

1. The time taken for building tree, which internally is related to memory allocation, are similar for Go and c++. Go is a little faster for its pre-allocated memory strategy.

2. The time taken for computing same(), which internally is a traverse of whole tree, is different. Go is around half of C++. That is the effect of parallelism where Go's goroutine with channecl can achieve. 

# the goroutine time order

Like thread, if you do not use sync way, the order of goroutine can not be guaranteed.

[Sample is from Go webiste](https://golang.org/ref/mem)
```
var a string

func hello() {
	go func() { a = "hello" }()
	print(a)
}
```

No guarantee when a will be "hello". Even it could never happen for the opitimized compilation. (E.g. If compiler found the goroutine has never been used, it can delete go goroutine() )

## 操作系统的信号量

信号量的概念是计算机科学家 **Dijkstra** （Dijkstra算法的发明者）提出来的，广泛应用在不同的操作系统中。系统中，会给每一个进程一个信号量，代表每个进程当前的状态，未得到控制权的进程，会在特定的地方被迫停下来，等待可以继续进行的信号到来。
如果信号量是一个任意的整数，通常被称为计数信号量（Counting semaphore），或一般信号量（general semaphore）；如果信号量只有二进制的0或1，称为二进制信号量（binary semaphore）。在linux系统中，二进制信号量（binary semaphore）又称[互斥锁](https://link.zhihu.com/?target=https%3A//zh.m.wikipedia.org/wiki/%25E4%25BA%2592%25E6%2596%25A5%25E9%2594%2581)（Mutex）
计数信号量具备两种操作动作，称为V（ `signal()` ）与P（ `wait()` ）（即部分参考书常称的“PV操作”）。V操作会增加信号量S的数值，P操作会减少它。
运行方式：

1. 初始化信号量，给与它一个非负数的整数值。
2. 运行P（ `wait()` ），信号量S的值将被减少。企图进入[临界区](https://link.zhihu.com/?target=https%3A//zh.m.wikipedia.org/wiki/%25E8%2587%25A8%25E7%2595%258C%25E5%258D%2580%25E6%25AE%25B5)的进程，需要先运行P（ `wait()` ）。当信号量S减为负值时，进程会被阻塞住，不能继续；当信号量S不为负值时，进程可以获准进入临界区。
3. 运行V（ `signal()` ），信号量S的值会被增加。结束离开[临界区段](https://link.zhihu.com/?target=https%3A//zh.m.wikipedia.org/wiki/%25E8%2587%25A8%25E7%2595%258C%25E5%258D%2580%25E6%25AE%25B5)的进程，将会运行V（ `signal()` ）。当信号量S不为负值时，先前被阻塞住的其他进程，将可获准进入[临界区](https://link.zhihu.com/?target=https%3A//zh.m.wikipedia.org/wiki/%25E8%2587%25A8%25E7%2595%258C%25E5%258D%2580%25E6%25AE%25B5)。




## Golang中的信号量

### golang 底层的信号量 runtime
该信号量实现无法在开发中使用。
Mutex中的sema就是该信号量。可以进行三个操作

``` 
func runtime_Semacquire(s *uint32) // 获取计数信号量

func runtime_SemacquireMutex(s *uint32, lifo bool, skipframes int) //获取互斥信号量

func runtime_Semrelease(s *uint32, handoff bool, skipframes int) // 释放信号量
```
实现：
https://www.cnblogs.com/ricklz/p/14610213.html
与chan的设计类似。当信号量不满足时，构建使用sudog结构，并放到信号量的队列，阻塞当前协程。在满足条件后，继续执行。

### 封装的信号量
在"golang.org/x/sync/semaphore" 包中，实现了一个高度封装的信号量。可以在开发中使用。
实现：https://zhuanlan.zhihu.com/p/337290029
该信号量设计思想都类似，只是这里的阻塞基于select。

在 `Go` 语言中信号量有时候也会被 `Channel` 类型所取代，因为一个 buffered chan 也可以代表 n 个资源。但是无法从chan中一次性获取多个资源，只能通过循环和计数来获取多个资源。即无法做到资源一次性分配，可能导致死锁。


不过既然 `Go` 语言通过 `golang.orgx/sync` 扩展库对外提供了 `semaphore.Weight` 这一种信号量实现，遇到使用信号量的场景时还是尽量使用官方提供的实现。在使用的过程中我们需要注意以下的几个问题：
-  `Acquire` 和  `TryAcquire` 方法都可以用于获取资源，前者会阻塞地获取信号量。后者会非阻塞地获取信号量，如果获取不到就返回 `false` 。
-  `Release` 归还信号量后，会以先进先出的顺序唤醒等待队列中的调用者。如果现有资源不够处于等待队列前面的调用者请求的资源数，所有等待者会继续等待。
- 如果一个 `goroutine` 申请较多的资源，由于上面说的归还后唤醒等待者的策略，它可能会等待比较长的时间。


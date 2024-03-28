## 应用
Go 语言中最常见的、也是经常被人提及的设计模式就是：**不要通过共享内存的方式进行通信，而是应该通过通信的方式共享内存**。**在很多主流的编程语言中，多个线程传递数据的方式一般都是共享内存，为了解决线程竞争，我们需要限制同一时间能够读写这些变量的线程数量**，然而这与 Go 语言鼓励的设计并不相同。
![alt text](image-19.png)


虽然我们在 Go 语言中也能使用共享内存加互斥锁进行通信，但是 Go 语言提供了一种不同的并发模型，即通信顺序进程（Communicating sequential processes，CSP）^[1](https://draveness.me/golang/docs/part3-runtime/ch06-concurrency/golang-channel/#fn:1)^。Goroutine 和 Channel 分别对应 CSP 中的实体和传递信息的媒介，Goroutine 之间会通过 Channel 传递数据。
![alt text](image-20.png)


## FIFO
1. Chan的数据流入和流出严格遵守FIFO，这是因为Chan的底层是RingBuffer。

2. Chan的发送和接受也严格遵守FIFO，这是因为Chan的发送和读取是队列。

- 先从 Channel 读取数据的 Goroutine 会先接收到数据；	
- 先向 Channel 发送数据的 Goroutine 会得到先发送数据的权利；

与之FIFO相对的模式是并发竞争模式：
- 发送方会向缓冲区中写入数据，然后唤醒接收方，多个接收方会尝试从缓冲区中读取数据，如果没有读取到会重新陷入休眠；
- 接收方会从缓冲区中读取数据，然后唤醒发送方，发送方会尝试向缓冲区写入数据，如果缓冲区已满会重新陷入休眠；
竞争模式容易导致惊群效应。（NGINX中惊群效用可以通过应用层锁accept_mutex或者操作系统层锁reuseport），并且无法保证先执行发送的G先发送数据到Chan，先执行读取的G先从Chan读取数据。

## 访问控制
多个G对Chan中数据并发访问（同时读或同时写）的控制策略：
1. **有锁（悲观锁）（官方）**：mutex。休眠和唤醒会带来额外的上下文切换，可能带来性能瓶颈。
2. 无锁（乐观锁）（社区）：CAS机制。不会休眠，性能较好。但是，目前社区通过 CAS 实现^[11](https://draveness.me/golang/docs/part3-runtime/ch06-concurrency/golang-channel/#fn:11)^的无锁 Channel 没有提供先进先出的特性。因为CAS不能保证，先得到数据的G能先获得向Chan发送数据的权利。

## 数据结构 （基于锁的并发访问实现）
```
type hchan struct {
	qcount   uint // total data in queue
	dataqsiz uint // size of the circular queue
	buf      unsafe.Pointer // circular queue buf
	elemsize uint16 
	closed   uint32 // flag of close
	elemtype *_type // element type
	sendx    uint  // buf send index
	recvx    uint  // buf recv index
	recvq    waitq // list of recv waiters 双向队列
	sendq    waitq // list of send waiters 双向队列

	lock mutex //有锁访问控制
}


type waitq struct {
	first *sudog //表示一个在等待列表中的 Goroutine和前后的指针
	last  *sudog
}

```
## makechan
make函数根据参数类型触发runtime.chan.go makechan函数：

- 如果当前 Channel 中不存在缓冲区，那么就只会为 [`runtime.hchan`](https://draveness.me/golang/tree/runtime.hchan) 分配一段内存空间；
- 在默认情况下会单独为 [`runtime.hchan`](https://draveness.me/golang/tree/runtime.hchan) 和缓冲区分配内存；

## chan <- i， G发送数据到chan
chan<-i被编译成chansend函数，发送数据的G调用chansend函数。函数流程如下（顺序执行，只有当上层逻辑不满足时，才继续执行下层）
### 1.访问控制和检查
如果channel为nil，那么调用gopark挂起该协程，并且没有入口恢复执行：`写nil chan会导致阻塞`。
如果 Channel已经关闭，那么会`panic` “send on closed channel”  中止程序。`写close chan会导致panic`。
在发送数据的逻辑执行之前会先为当前 Channel`加锁`，防止多个线程并发修改数据。
### 2.判断是否可以直接发送
Chan上已经有处于读等待的 Goroutine（说明chan中没有数据），那么 [`runtime.chansend`](https://draveness.me/golang/tree/runtime.chansend) 会从接收队列  `recvq`  中取出最先陷入等待的 Goroutine的sudog。 根据sudog直接向它发送数据。
- 调用 [`runtime.sendDirect`](https://draveness.me/golang/tree/runtime.sendDirect) 将发送的数据直接拷贝到  `x = <-c`  表达式中变量  `x`  所在的内存地址上；
- 调用 [`runtime.goready`](https://draveness.me/golang/tree/runtime.goready) 将等待接收数据的 Goroutine 标记成可运行状态  `Grunnable`  并把该 Goroutine 放到发送方所在的处理器的  `runnext`  上等待执行，该处理器在下一次调度时会立刻唤醒数据的接收方。发送数据的过程只是将接收方的 Goroutine 放到了处理器的  `runnext`  中，程序没有立刻执行该 Goroutine。
### 3.判断是否可以放入缓冲区
如果 Channel 存在缓冲区并且其中还有空闲的容量，会直接将数据存储到缓冲区  `sendx`  所在的位置上。
主要涉及sendx和qcount的计算，sendx需要进行取模等操作。
### 4.发送阻塞
```
func chansend(c *hchan, ep unsafe.Pointer, block bool, callerpc uintptr) bool {
...
if !block {
unlock(&amp;c.lock)
return false
}

gp := getg() // 获取发送数据的 Goroutine；
mysg := acquireSudog() // 获取sudog结构
mysg.elem = ep //设置sudog结构的相关信息
mysg.g = gp
mysg.c = c
gp.waiting = mysg // 将sudog设置到当前 Goroutine 的`waiting`上，表示 Goroutine 正在等待该  `sudog` 准备就绪；

c.sendq.enqueue(mysg) // 将sudog加入到channel的发送队列末尾
goparkunlock(&c.lock, waitReasonChanSend, traceEvGoBlockSend, 3) //阻塞当前Groutine，等待唤醒

// 被唤醒后执行
gp.waiting = nil
gp.param = nil
mysg.c = nil
releaseSudog(mysg) //释放sudog

return true 
}
```

1. 调用 [`runtime.getg`](https://draveness.me/golang/tree/runtime.getg) 获取发送数据使用的 Goroutine；
2. 执行 [`runtime.acquireSudog`](https://draveness.me/golang/tree/runtime.acquireSudog) 获取 [`runtime.sudog`](https://draveness.me/golang/tree/runtime.sudog) 结构并设置这一次阻塞发送的相关信息，例如发送的 Channel、是否在 select 中和待发送数据的内存地址等；
3. 将刚刚创建并初始化的 [`runtime.sudog`](https://draveness.me/golang/tree/runtime.sudog) 加入发送等待队列，并设置到当前 Goroutine 的  `waiting`  上，表示 Goroutine 正在等待该  `sudog`  准备就绪；
4. 调用 [`runtime.goparkunlock`](https://draveness.me/golang/tree/runtime.goparkunlock) 将当前的 Goroutine 陷入沉睡等待唤醒；
5. 被调度器唤醒后会执行一些收尾工作，将一些属性置零并且释放 [`runtime.sudog`](https://draveness.me/golang/tree/runtime.sudog) 结构体；

## i <- chan， G从chan接受数据
最终调用[`runtime.chanrecv`](https://draveness.me/golang/tree/runtime.chanrecv)函数。函数流程如下（顺序执行，只有当上层逻辑不满足时，才继续执行下层）
### 1.访问控制和检查
如果 Channel 为nil，那么会直接调用 [`runtime.gopark`](https://draveness.me/golang/tree/runtime.gopark) 挂起当前 Goroutine。并且没有入口恢复执行：`读nil chan会导致阻塞`。
如果 Channel 已经关闭并且缓冲区没有任何数据，[`runtime.chanrecv`](https://draveness.me/golang/tree/runtime.chanrecv) 会直接返回。`读close chan且缓冲区无数据时，返回零值和false`
### 2.判断是否可以直接接受
如果 Channel 的  `sendq`  队列中存在挂起的 Goroutine，会将  `recvx`  索引所在的数据拷贝到接收变量所在的内存空间上并将  `sendq`  队列中 Goroutine 的数据拷贝到缓冲区，释放一个阻塞的发送方。
运行时调用 [`runtime.goready`](https://draveness.me/golang/tree/runtime.goready) 将当前处理器的  `runnext`  设置成发送数据的 Goroutine，在调度器下一次调度时将阻塞的发送方唤醒。
### 3.判断是否可以放入缓冲区
如果 Channel 的缓冲区中包含数据，那么直接读取  `recvx`  索引对应的数据；
### 4.发送阻塞
在默认情况下会挂起当前的 Goroutine，将 [`runtime.sudog`](https://draveness.me/golang/tree/runtime.sudog) 结构加入  `recvq`  队列并陷入休眠等待调度器的唤醒。即：
使用 [`runtime.sudog`](https://draveness.me/golang/tree/runtime.sudog) 将当前 Goroutine 包装成一个处于等待状态的 Goroutine 并将其加入到接收队列中。
完成入队之后，上述代码还会调用 [`runtime.goparkunlock`](https://draveness.me/golang/tree/runtime.goparkunlock) 立刻触发 Goroutine 的调度，让出处理器的使用权并等待调度器的调度。


## close(chan)
### 1.检查
当 Channel 是一个空指针或者已经被关闭时，都会直接panic
### 2.逻辑

1. 将  `recvq`  和  `sendq`  两个队列中的数据加入到 Goroutine 列表  `gList`  中，与此同时该函数会清除所有 [`runtime.sudog`](https://draveness.me/golang/tree/runtime.sudog) 上未被处理的元素
2. 遍历gList，为所有被阻塞的 Goroutine 调用 [`runtime.goready`](https://draveness.me/golang/tree/runtime.goready) 触发调度。


## 细节
select 中对chan的操作，可以实现不阻塞Groutine

## 参考
https://draveness.me/golang/docs/part3-runtime/ch06-concurrency/golang-channel/#64-channel
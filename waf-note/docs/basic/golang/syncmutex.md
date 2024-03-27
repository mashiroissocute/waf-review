## 资料
https://draveness.me/golang/docs/part3-runtime/ch06-concurrency/golang-sync-primitives/

## Sync.Mutex
https://juejin.cn/post/7086756462059323429
### 注意点
Sync.Mutex是值类型
参数传递，调用非引用结构体方法的时候，都会导致复制。锁不能复制，因为lock和unlock可能操作不同对象
https://www.xiayinchang.top/post/6b348626.html
### 数据结构
``` 
type Mutex struct {
    state int32
    sema  uint32
}
const (
   mutexLocked = 1 << iota // mutex is locked  1<<0 = 1
   mutexWoken    							// 1<<1 = 2
   mutexStarving 							// 1<<2 = 4
   mutexWaiterShift = iota				    // 3
)
```
mutex 的 state 有 32 位，它的低 3 位分别表示 3 种状态：**唤醒状态**、**上锁状态**、**饥饿状态**，剩下的位数则表示当前阻塞等待的 goroutine 数量。
![enter image description here](/tencent/api/attachments/s3/url?attachmentid=4873038)

mutex 会根据当前的 state 状态来进入**正常模式**、**饥饿模式**或者是**自旋**。

sema时runtime的信号量：https://iwiki.woa.com/pages/viewpage.action?pageId=4007535298

### mutex 正常模式
当 mutex 调用 Unlock() 方法释放锁资源时，如果发现有等待唤起的 Goroutine 队列时，则会将队头的 Goroutine 唤起。
队头的 goroutine 被唤起后，会调用 CAS 方法去尝试性的修改 state 状态，如果修改成功，则表示占有锁资源成功。
(注：CAS 在 Go 里用 `atomic.CompareAndSwapInt32(addr *int32, old, new int32) `方法实现，CAS 类似于乐观锁作用，修改前会先判断地址值是否还是 old 值，只有还是 old 值，才会继续修改成 new 值，否则会返回 false 表示修改失败。)

### mutex 饥饿模式

由于上面的 Goroutine 唤起后并不是直接的占用资源，还需要调用 CAS 方法去**尝试性**占有锁资源。如果此时有新来的 Goroutine，那么它也会调用 CAS 方法去尝试性的占有资源。
但对于 Go 的调度机制来讲，会比较偏向于 CPU 占有时间较短的 Goroutine 先运行，而这将造成一定的几率让新来的 Goroutine 一直获取到锁资源，此时队头的 Goroutine 将一直占用不到，导致**饿死**。
针对这种情况，Go 采用了饥饿模式。即通过判断队头 Goroutine 在超过一定时间后还是得不到资源时，会在 Unlock 释放锁资源时，直接将锁资源交给队头 Goroutine，并且将当前状态改为**饥饿模式**。
后面如果有新来的 Goroutine 发现是饥饿模式时， 则会直接添加到等待队列的队尾。
如果一个 Goroutine 获得了互斥锁并且它在队列的末尾或者它等待的时间少于 1ms，那么当前的互斥锁就会切换回正常模式。

### mutex woken的作用

当一个goroutine在lock自旋的过程中，成功获取了锁，会将state的woken设置为1。

当一个goroutine在unlock时，发现state的woken是1，知道已经有一个goroutine通过自旋获取到了锁，就不在通过runtime_Semrelease去释放一个信号量，并唤醒信号量上阻塞的goroutine。

这个逻辑，仅当mutex正常模式才执行。 当mutex处于饥饿模式，会直接调用runtime_Semrelease，并唤醒阻塞在信号量上的goroutine。



### mutex 自旋
如果 Goroutine 占用锁资源的时间比较短，那么每次都调用信号量来阻塞唤起 goroutine，将会很**浪费**资源。
因此在符合一定条件后，mutex 会让当前的 Goroutine 去**空转** CPU，在空转完后再次调用 CAS 方法去尝试性的占有锁资源，直到不满足自旋条件，则最终会加入到等待队列里。
自旋的条件如下：

- 还没自旋超过 4 次
- 多核处理器
- GOMAXPROCS > 1
- p 上本地 Goroutine 队列为空
可以看出，自旋条件还是比较严格的，毕竟这会消耗 CPU 的运算能力。


### Lock() 过程
首先，如果 mutex 的 state = 0，即没有谁在占有资源，也没有阻塞等待唤起的 goroutine。则会调用 CAS 方法去尝试性占有锁，不做其他动作。
如果不符合 m.state = 0，则进一步判断是否需要自旋。
当不需要自旋又或者自旋后还是得不到资源时，此时会调用 runtime_SemacquireMutex 信号量函数，将当前的 goroutine 阻塞并加入等待唤起队列里。
如果有新来的 Goroutine 发现是饥饿模式时， 则会直接添加到等待队列的队尾。

### Unlock 过程
mutex 的 Unlock() 
如果当前是正常模式，如果有woken标志，简单的返回。没有woken标识，**唤起**队头 Goroutine，让其和其他正在加锁或者自旋的goroutine竞争。（woken标志，表示已经有goroutine自旋成功了）
如果是饥饿模式，则会**直接**将锁交给队头 Goroutine，然后唤起队头 Goroutine，让它继续运行。



## Golang死锁问题
在并发情况下，如果所有协程都因为等待资源而被阻塞，则会陷入死锁的状态。
容易导致golang死锁场景：

- 使用非缓冲的chan
- 使用mutex rwmutex, 协程1: mutex1.lock .. mutex2.lock 协程2:mutex2.lock .. mutex1.lock

死锁在编译的时候是没有办法发现的，并且也不好在代码中进行定位。
所以，遇到死锁的时候，通常会借助第三方的死锁检测工具来检查代码中存在的死锁。

常见的golang死锁检测工具： https://github.com/sasha-s/go-deadlock
例子： https://blog.csdn.net/DisMisPres/article/details/123402901



## Sync.RWLock
https://segmentfault.com/a/1190000039712353
https://juejin.cn/post/7002212463600992292
``` 
type RWMutex struct {
    w           Mutex  // held if there are pending writers
    writerSem   uint32 // semaphore for writers to wait for completing readers
    readerSem   uint32 // semaphore for readers to wait for completing writers
    readerCount int32  // number of pending readers
    readerWait  int32  // number of departing readers 防止写饿死
}
```

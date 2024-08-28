## 同步与互斥
通过制约方式，来实现进程、线程、协程间协作（同步）和对互斥资源访问（互斥）。

- **间接相互制约（互斥）**：因为进程在并发执行的时候共享临界资源而形成的相互制约的关系，需要对临界资源互斥地访问；
- **直接制约关系（同步）**：多个进程之间为完成同一任务而相互合作而形成的制约关系。
不管是同步还是互斥，制约方式都是一致的，只是使用方式不同。例如，使用信号量实现
互斥： （进程1 ：P 执行 V     进程2 ： P 执行 V）保证同一时刻进程2在进程1只有一个执行
同步： （进程1:   执行 V         进程2:    P 执行） 保证进程2在进程1后执行
https://cloud.tencent.com/developer/article/1803377

## 原子操作
原子操作是指，该操作没有中间状态。任何其他core都感知不到本次对变量的操作。
例如，在写int64变量的时候，32bit机器会执行两条指令分别写前后的32bit。该操作并非原子性，其他的core可能看到该变量只修改了前32bit的中间状态，而产生难以排查的错误。
https://blog.51cto.com/u_15127545/3436521

### 操作系统的原子操作

一条CPU指令是原子操作。

多条指令组合执行并且还能保证原子操作，称为原语。例如CompareAndSet（CAS）将多个指令合并成一个指令，在原语的执行过程中也是不响应中断的，使之成为原子操作。这个期间，等于是屏蔽中断。因此原语操作的指令长序应该是短小精悍的，这样才能保证系统的效率。

### 语言中的原子操作

如上例子，当多个Goroutine同时执行的时候，如果想要对一个变量修改。那么通过原子操作是可以避免中间状态错误的。想要实现原子操作，有多种方式。有互斥锁，有硬件层面的CAS原语等。Golang在sync.atomic包中实现了对int64等变量的原子操作，底层使用的是硬件层面的CAS原语。
https://gfw.go101.org/article/concurrent-atomic-operation.html

python没有提供原子操作。需要通过加锁实现保护。

## 操作系统通信机制（进程）


### 共享内存
将不同进程的一片虚拟内存地址映射到同一片物理内存地址。


### socket
基于本地或者远程socket网络通信

## 操作系统同步机制（进程和线程）

### 互斥锁

### 读写锁

### 自旋锁

当一个线程尝试获取已经被另一个线程持有的自旋锁时，该线程会持续检查锁的状态（"自旋"），而不是进入睡眠状态。减少上下文切换次数。
自旋锁适用于锁被持有的时间非常短的情况，因为线程在等待锁释放时会持续占用CPU资源。
自旋锁可能会导致"饥饿"现象，如果锁被频繁地获取和释放，等待的线程可能无法及时获得锁。

### 条件变量
实现wait和notifyall的功能。需要和互斥锁结合使用

互斥锁与条件变量的区别：
互斥锁主要用于保护共享资源，防止并发访问导致的问题。
条件变量用于协调线程或协程之间的执行，允许它们基于某些条件进行同步。例如构造生产者消费者模式。

### 信号量
信号量除了初始化外，仅能被通过两个标准的原子操作 `wait(S)` 和 `signal(S)` 来访问，这两个操作也被称为P、V操作，都是原语操作。
``` 
//记录型信号量
typedef struct{
int value;
struct process_control_block * list;
}semaphore;
// P操作
wait(semaphore *S)
{
S->value --;
if(S->Value < 0){
block( S->list );
}
]
signal(semaphore *S){
S->value++;
if(S->value <= 0){
wakeup(S->list);
}
}
```
https://blog.csdn.net/Ciellee/article/details/107331288


互斥锁和信号量的区别：
Mutex 相比信号量增加了所有权的概念，一只锁住的 Mutex 只能由给它上锁的线程解开，只有系铃人才能解铃。Mutex 的功能也就因而限制在了构造临界区上。



## 语言层面同步机制（协程/线程）

语言可以自己实现同步机制，部分同步机制是基于操作系统提供的同步机制进行封装的，部分是自己实现的控制。
不同语言支持的锁类型也存在不同。

### python

线程锁: threading.Lock() threading.Condition()
协程锁: asyncio.Lock() asyncio.Queue()


python条件变量也是基于操作系统的条件变量API实现。
```python
import threading

# 创建一个条件变量对象
condition = threading.Condition()

# 共享资源
data = []

def producer():
    with condition:
        # 生产数据
        data.append('some data')
        print("Produced data")
        # 通知等待的消费者线程
        condition.notify()

def consumer():
    with condition:
        # 等待直到有数据可消费
        while not data:
            condition.wait()
        # 消费数据
        print("Consumed", data.pop(0))

# 创建并启动生产者和消费者线程
t1 = threading.Thread(target=producer)
t2 = threading.Thread(target=consumer)
t1.start()
t2.start()
```

### golang
sync.Mutex
sync.RwMutex
sync.Cond
chan机制

### 乐观锁和悲观锁
这两个并不是实际的锁类型，而是锁的思想。
乐观锁的思想是，很少发生冲突，读不加锁，更新的时候，判断版本是否发生变化。
```golang
type OptimisticLock struct {
	value    int
	version  int
}

func (l *OptimisticLock) Update(newValue int) bool {
	currentVersion := l.version
	// 模拟一些耗时操作
	// ...
	if l.version == currentVersion {
		l.value = newValue
		l.version++
		return true
	} else {
		return false
	}
}
```

悲观锁的思想是，经常发生冲突，读写都需要加锁。



## 死锁

### 条件
- 互斥
- 持有 等待
- 非抢占
- 循环

### 死锁预防
打破上面的四个条件之一。

具体到语言中，主要是：

- 打破持有等待。为对象设置超时事件，到了超时时间后，释放资源。
- 打破循环。通过条件变脸，协调执行顺序，避免造成循环。


### 死锁避免
无法打破上面的条件，避免在使用时形成死锁。

主要采用银行家算法，在分配资源之前，检查是否可以依靠分配的资源完全执行并释放资源。最终形成一个可执行序列。




### 死锁检测

即使注意了 死锁预防。但是在程序中很少使用死锁避免手段。所以还有存在可能形成死锁。死锁检测思路：

- 语言层面调试能力，例如python-threading.debug、golang-race检测。
- 代码走读检测。
- 第三方库工具。pylock、github.com/sasha-s/go-deadlock等。



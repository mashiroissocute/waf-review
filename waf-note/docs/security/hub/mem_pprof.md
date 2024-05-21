
### gc原理
[垃圾回收](../../basic/golang/gc.md)

### 操作系统内存管理
[操作系统内存管理](../../basic/os/mem.md)

### Golang内存分配
[堆内存分配](../../basic/golang/heap.md)

[栈内存分配](../../basic/golang/stack.md)


### gc问题
Go的GC会在运行时暂停应用程序的执行，以便回收不再使用的内存。这些暂停时间可能会导致应用程序的延迟增加，从而影响性能。

HUB是一个对响应延时要求较高的应用，应当尽量避免GC对程序延时的影响。

#### GC对延时的影响: 
使用wrk压测，长尾延时较高，不满足最低时延要求。
通过统计gc次数，较为频繁。
![alt text](image-9.png)




### 调优方案
对性能有较高要求的系统面对自动垃圾回收型语言中的STW时或多或少会导致性能下降，对于gc调优主要包括以下两个部分：

- 1.减少临时堆对象的创建
- 2.优化gc触发策略


#### 减少临时堆对象的创建
##### 栈内存会随着函数的调用进行分配和回收

堆内存是由程序申请分配，需要gc回收，因此需要尽量避免堆内存的使用。当然，Go会尽可能的将内存分配到栈上，当分配到栈上可能导致譬如内存非法访问等问题时会使用堆内存，**通常**分配原则包括：

- 1.Sharing down typically stays on the stack!! 
在调用方函数内部创建的对象通过参数的形式传递给被调用方时，该变量会使用栈内存
- 2.Sharing up typically escapes to the heap!!
在被调用函数内部创建的对象通过指针形式返回给调用方时，该变量会使用堆内存


##### 对象复用：对象池｜syncPool

对于需要频繁创建同一类对象，且创建成本较高时，可以通过syncPool保存和复用堆上的对象，减少内存分配，降低 GC 压力
``` 
goos: linux
goarch: amd64
pkg: waf_bypass_hub/util/jsonParser
cpu: AMD EPYC 7K62 48-Core Processor
BenchmarkJsonMarshalNoPool-8      301642              4529 ns/op            1160 B/op         60 allocs/op
BenchmarkJsonMarshalPool-8        307179              3965 ns/op             200 B/op         40 allocs/op
```


#### 优化gc触发策略
##### GCPercent
##### Tuner
##### Ballast
- ballast分配的是虚拟内存还是物理内存？

![alt text](image-10.png)
profiling 是指在程序执行过程中，收集能够反映程序执行状态的数据。


一般常规分析内容：

- cpu：程序对cpu的使用情况 - 使用时长，占比等
- 内存：程序对内存的使用情况 - 使用占比，内存泄露等。如果在往里分，程序堆、栈使用情况
- I/O：IO的使用情况 - 哪个程序IO占用时间比较长
- goroutine：go的协程使用情况，调用链的情况
- goroutine leak：goroutine泄露检查
- go dead lock：死锁的检测分析
- data race detector：数据竞争分析，其实也与死锁分析有关


## 分析原理

CPU采样会记录所有的调用栈和它们的占用时间

在采样时，进程会每秒暂停一百次，每次会记录当前的调用栈信息。

汇总之后，根据调用栈在采样中出现的次数来推断函数的运行时间。 





https://debug-lixiwen.github.io/2021/07/18/shi-zhan/


https://blog.wolfogre.com/posts/go-ppof-practice/#%E5%89%8D%E8%A8%80
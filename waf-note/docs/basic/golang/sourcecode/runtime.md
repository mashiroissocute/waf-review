## runtime.getg()
从线程的TLS中获取当前线程正在执行的g。
在编译器中实现了函数具体逻辑。
`MOVQ (TLS), r`
https://imkira.com/Golang-runtime-getg/

## runtime.get_tls() / runtime.g()
FS 寄存器里面是 m.tls 的地址。
在 Golang 的汇编代码中，FS 寄存器实际上是 TLS 虚拟寄存器。
因此 get_tls(r) 就是将TLS 寄存器的地址赋给参数r。
g(r) 就是获取在 TLS 上存储的 g 结构体

``` 
#define get_tls(r)  MOVQ TLS, r
#define g(r)    0(r)(TLS*1)
```

```
get_tls(CX) // 获取m的TLS存储，TLS地址放在TLS寄存器中。
MOVQ g(CX), BX; // 根据TLS地址，计算g。
```

## runtime.systemstack()
每个 M 都有一个 g0，这个 g0 的栈空间使用的是 M 物理线程的栈空间，而不是其他 g 那样其实是 Golang 的堆空间分配的。这个 g0 唯一的作用就是利用它的栈空间来执行一些函数。
systemstack 会切换当前的 g 到 g0, 并且使用 g0 的栈空间, 然后调用传入的函数, 再切换回原来的g和原来的栈空间。 
切换到 g0 后会假装返回地址是 mstart, 这样 traceback 的时候可以在 mstart 停止。新建g0时，go.sched.pc和gp.sched.sp也是指向mstart。
``` 
systemstack(func() {
            newg.stack = stackalloc(uint32(stacksize))
        })
```
通过systemstack切换到m的g0栈去执行函数。g0栈位于操作系统分配的栈空间。其他g的栈位于操作系统分配的堆空间。

## runtime.mcall()
mcall 与 runtime·systemstack 非常相似，都是切换到g0，利用 g0 栈空间执行新的函数。
但是 执行完函数之后 runtime·systemstack会切回g。 mcall 保存上下文到g 后，没有切回g的逻辑。
```
func goexit1() {
    mcall(goexit0) // 切到g0调用退出函数，运行不再切回当前g
}
```

## runtime.gogo()
gogo 是专门从 g0 切换到 g 栈空间来执行代码的函数。
``` 
gogo(&gp.sched)
```


## doInit
``` 
doInit(&runtime_inittask) // 执行runtime包的init函数

doInit(&main_inittask) // 执行main包的init函数
```
函数内容： 
- 初始化 当前  `package`  所依赖的  `package` 
- 初始化 当前  `package`  中所有需要初始化的变量
- 执行 用户所写的  `func init()`

https://blog.leonard.wang/archives/inittask

## linkname
``` 
golang会为每个包都按需生产一个packagename..inittask 用于初始化函数调用

//go:linkname runtime_inittask runtime..inittask
var runtime_inittask initTask

//go:linkname main_inittask main..inittask
var main_inittask initTask

//go:linkname main_main main.main
func main_main()
```
可以用于使用其他包中的私有函数或者变量（小写）


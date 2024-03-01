[toc]
# Golang builtin func （内建函数）
### 使用要点
- 如果 panic 和 recover 发生在同一个协程，那么 recover 是可以捕获的，如果 panic 和 recover 发生在不同的协程，那么 recover 是不可以捕获的。只要一个协程就可以，就是放在不同的函数中。
- panic后，会停止后续程序执行
### 1. 一个使用场景

``` 
DBTx := models.CmdDB.Begin() //开启事物
defer func() { //捕获错误，回滚
    if r := recover(); r != nil || err != nil {
       DBTx.Rollback()
       logger.Errorf("Rollback err: %s recover: %s, stack:%s", err, r, string(debug.Stack()))
    }
}()
ngServerInfo, err := dc.AddDomainProtection(DBTx, domainID) //传入DBTx
if err != nil {
   panic(err)
}
DBTx.Commit() // 提交DB成功
```

### 2. defer
- 什么是defer：go语言注册延时调用的机制，使用defer修饰的函数会在当前函数执行完毕后执行。
- defer编译伪码![enter image description here](/tencent/api/attachments/s3/url?attachmentid=4873003)
- deferproc
	- 首先编译器把 `defer` 语句翻译成对应的 `deferproc` 函数的调用
	- 然后 `deferproc` 函数通过 `newdefer` 函数分配一个 `_defer` 结构体对象并放入当前的 `goroutine` 的 `_defer` 链表的表头；
	- 在 _defer 结构体对象中保存被延迟执行的函数 fn 的地址以及 fn 所需的参数
	- 返回到调用 deferproc 的函数继续执行后面的代码。
- deferreturn
	- 首先我们通过当前 `goroutine` 对应的 `g` 结构体对象的 `_defer` 链表判断是否有需要执行的 `defered` 函数，如果没有则返回；这里的没有是指g._defer== nil 或者 `defered` 函数不是在 `deferteturn` 的 `caller` 函数中注册的函数。
	- 然后我们在从 `_defer` 对象中把 `defered` 函数需要的参数拷贝到栈上，并释放 `_defer` 的结构体对象。
	- 最后调用 `jmpderfer` 函数调用 `defered` 函数，也就是 `defer` 关键字中传入的函数.
- defer内存模型![enter image description here](/tencent/api/attachments/s3/url?attachmentid=4873004)
- 多个defer的维护方式是链表，新出现的defer将放在链表头部。gorouting有一个指向defer链表的指针。该链表在堆上分配空间
- defer链表中的节点是defer结构体，其中有该defer是否执行的标志/函数的入口/参数/指向下一个defer结构体的link指针。
- 创建defer结构体时，会将涉及到的变量进行值拷贝放入堆中。在调用defer注册的函数时，又会将该值拷贝到栈中。

### 3. panic recover
- panic会停掉当前正在执行的程序（注意，不只是协程），但是与 `os.Exit(-1)` 直接退出程序不同，panic会先处理完当前goroutine已经defer挂上去的任务，执行完毕后，输出panic信息再退出整个程序。
-  `panic执行，且只执行，当前goroutine的defer`。当前goroutine中的所有defer都会执行，即使是在goroutine的函数调用中panic了：
``` 
go func() {
        defer fmt.Println("defer caller")  //会被执行
        func() {
            defer func() {
                fmt.Println("defer here")
            }()
            if user == "" {
                panic("should set user env.")
            }
        }()
}()
```

- `recover只会在defer中发挥作用`
- recover会将panic中的recovered字段设置为true。`panic流程中在调用完每个defer以后会检查recovered标记，如果为true则会退出panic流程，panic信息并不会输出。`
- 可以使用在当前`goroutine`的defer中使用recover来捕捉panic，虽然该goroutine的panic后的代码不会再继续执行，但是不会导致整个程序退出。
- 例子：
``` 
package main

import (
    "fmt"
    "time"
)

func main() {
    defer fmt.Println("defer main") // will this be called when panic?
    var user = ""
    go func() {
        defer func() {
            fmt.Println("defer caller")
            if err := recover(); err != nil {
                fmt.Println("recover success.")
            }
        }()
        func() {
            defer func() {
                fmt.Println("defer here")
            }()
            if user == "" {
                panic("should set user env.")
            }
            fmt.Println("after panic")
        }()
    }()
    time.Sleep(1 * time.Second)
    fmt.Println("get result")
}
```
将会输出：
``` 
defer here
defer caller
recover success.
get result
defer main
```
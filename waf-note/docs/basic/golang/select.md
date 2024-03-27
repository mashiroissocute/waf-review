## Select
Go select的用法与操作系统多路IO复用的select比较类似。但是其监听的是channel，而不是socket。
### 特性

- select的每个 case 必须是一个channel通信操作（I/O操作），要么是发送到channel要么是从channel中接收。
- 多个case同时可运行时，随机选择一个执行，此时将无法保证执行顺序
- 在没有default的情况下，如果没有 case 可运行，它将阻塞，直到有 case 可运行。在有default的情况下，如果没有 case 可运行，它将执行default。
- 对于case条件语句中，如果channel被设置为nil，则该分支将被永远阻塞。因为读写nil channel都将被永久阻塞。
- 如果有定时条件语句，那么判断逻辑为**如果在这个时间段内（定时时间）一直没有满足条件的case,则执行这个时间case**。**如果此段时间内出现了可操作的case,则直接执行这个case。一般用超时语句代替了default语句**
- 空的select{}，会直接阻塞当前goroutine

## 应用
### 多路复用

### 时间
`如果有定时条件语句，那么判断逻辑为如果在这个时间段内（定时时间）一直没有满足条件的case,则执行这个时间case。如果此段时间内出现了可操作的case,则直接执行这个case。`
因此，这并不是一个定时执行的概念。而是一个超时执行的概念。例如:
``` 
for {
        select {
        case <-time.After(time.Second * 2):
            fmt.Println("2 period")
        case <-time.After(time.Second * 5):
            fmt.Println("5 period")
        }
    }
```
5 period这个分支永远都执行不到，因为在5s内，一直有可满足的分支到来（2 period）

#### 定时器
定时器，就需要没有其他分支可以在时间段内可以执行，可以用阻塞的分支nil或者不要其他的分支。
``` 
for {
        select {
        case <-time.After(time.Second * time.Duration(intervalSec)):
           .. do something
		case <- donec: //结束定时器
			return

        }
    }
```


``` 
c := make(chan int, 10)
close(c) 

for {
	select {
	case _, ok := <-c:
		if !ok {
			c = nil
		}

	case <-time.After(time.Second * 5):
		fmt.Println("5 period")
	}
}
```



#### 超时器
```
ch1 := make(chan int, 10)
go func() {
        time.Sleep(10 * time.Second)
        ch1 <- 1
    }()

select {
    case str := <- ch1
        fmt.Println("receive str", str)
    case <- time.After(time.Second * 5): 
        fmt.Println("timeout!!")
}
```

### 判断缓冲区是够可读可写

``` 
buffer := make(chan int, 1) 读写轮流打印
// buffer := make(chan int) 永远不可读写

for {
	select {
	case <-buffer:
		fmt.Println("readable")
	case buffer <- 1:
		fmt.Println("writeable")
	default:
		fmt.Println("can not read or write")
	}
	time.Sleep(1 * time.Second)
}
```

### 阻塞go routine
```
func main()  {
    bufChan := make(chan int)
    
    go func() {
        for{
            bufChan <-1
            time.Sleep(time.Second)
        }
    }()


    go func() {
        for{
            fmt.Println(<-bufChan)
        }
    }()
     
    select{} //防止main退出
}
```












## 参考
https://www.jianshu.com/p/de4bc02e7c72
https://wudaijun.com/2017/10/go-select/
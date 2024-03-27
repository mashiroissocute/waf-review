## make
仅可用于slice、map、chanle的初始化
分配空间并初始化为0值
返回引用

## new
初始化任意类型
分配空间，且将内存清零，不会初始化。
因此值类型（int string array）都被赋予零值。
int 0 
string
array [n]int{0,0,0,0,0}
但是引用类型 （map slice chan）就会被赋予nil值
使用new分配引用类型的内存，因为返回的是应用类型的nil值，再操作该引用类型会panic


array 数组示例
```
    var a [5]int
    fmt.Printf("a: %p %#v \n", &amp;a, a)//a: 0xc04200a180 [5]int{0, 0, 0, 0, 0} 
    av := new([5]int)
    fmt.Printf("av: %p %#v \n", &amp;av, av)//av: 0xc000074018 &amp;[5]int{0, 0, 0, 0, 0}
    (*av)[1] = 8
    fmt.Printf("av: %p %#v \n", &amp;av, av)//av: 0xc000006028 &amp;[5]int{0, 8, 0, 0, 0}
```

silce 示例

```
    var a *[]int
    fmt.Printf("a: %p %#v \n", &amp;a, a) //a: 0xc042004028 (*[]int)(nil)
    av := new([]int)
    fmt.Printf("av: %p %#v \n", &amp;av, av) //av: 0xc000074018 &amp;[]int(nil)
    (*av)[0] = 8
    fmt.Printf("av: %p %#v \n", &amp;av, av) //panic: runtime error: index out of range
```

map 示例

```
    var m map[string]string
    fmt.Printf("m: %p %#v \n", &amp;m, m)//m: 0xc042068018 map[string]string(nil) 
    mv := new(map[string]string)
    fmt.Printf("mv: %p %#v \n", &amp;mv, mv)//mv: 0xc000006028 &amp;map[string]string(nil)
    (*mv)["a"] = "a"
    fmt.Printf("mv: %p %#v \n", &amp;mv, mv)//这里会报错panic: assignment to entry in nil map
```

channel示例

```
cv := new(chan string)
fmt.Printf("cv: %p %#v \n", &amp;cv, cv)//cv: 0xc000074018 (*chan string)(0xc000074020) 
//cv <- "good" //会报 invalid operation: cv <- "good" (send to non-chan type *chan string)
```




## 区别

- make和new都是golang用来分配内存的內建函数，且在堆上分配内存，make 即分配内存，也初始化内存。new只是将内存清零，并没有初始化内存。（都会分配内存，主要是内存清零和内存初始化的区别。）
- make返回的还是引用类型本身；而new返回的是指向类型的指针。
- make只能用来分配及初始化类型为slice，map，channel的数据；new可以分配任意类型的数据。

https://sanyuesha.com/2017/07/26/go-make-and-new/
https://cloud.tencent.com/developer/article/1706196
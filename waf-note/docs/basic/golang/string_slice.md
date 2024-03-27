## string
string是一片连续的内存空间，我们也可以将它理解成一个由字符组成的数组。一个字符串是只读的，因此我们并不会直接向字符串直接追加元素改变其本身的内存空间，所有在字符串上的写入操作都是通过拷贝实现的。
https://learnku.com/go/t/46493 （字符串不可变的例子）

## 底层
```
type StringHeader struct {
	Data uintptr
	Len  int
}
```

## 字符串拼接
因为字符串是只读的，因此字符串拼接。在正常情况下，运行时会调用  `copy`  将输入的多个字符串拷贝到目标字符串所在的内存空间。新的字符串是一片新的内存空间，与原来的字符串也没有任何关联。 所以，一旦需要拼接的字符串非常大，拷贝带来的性能损失是无法忽略的。
https://geektutu.com/post/hpg-string-concat.html （字符串拼接性能对比）

## 字符串到切片转换
当我们使用 Go 语言解析和序列化 JSON 等数据格式时，经常需要将数据在  `string`  和  `[]byte`  之间来回转换，类型转换的开销并没有想象的那么小。

- [`runtime.slicebytetostring`](https://draveness.me/golang/tree/runtime.slicebytetostring) 创建string结构体，并通过 [`runtime.memmove`](https://draveness.me/golang/tree/runtime.memmove) 将原  `[]byte`  中的字节全部复制到新的内存空间中。

- [`runtime.stringtoslicebyte`](https://draveness.me/golang/tree/runtime.stringtoslicebyte) 运行时会调用 [`runtime.rawbyteslice`](https://draveness.me/golang/tree/runtime.rawbyteslice) 创建新的字节切片并将字符串中的内容拷贝过去。

二者都存在内容拷贝。
字符串和  `[]byte`  中的内容虽然一样，但是字符串的内容是只读的，我们不能通过下标或者其他形式改变其中的数据，而  `[]byte`  中的内容是可以读写的。不过无论从哪种类型转换到另一种都需要拷贝数据，而内存拷贝的性能损耗会随着字符串和  `[]byte`  长度的增长而增长。

## 小结
字符串在做拼接和类型转换等操作时一定要注意性能的损耗，遇到需要极致性能的场景一定要尽量减少类型转换的次数。

## slice
slice和array
array是值类型，赋值和函数传参操作都会复制整个数组数据。
slice是引用类型.
array的长度固定，slice的长度可变。
使用下标初始化的slice不会拷贝原数组或者原切片中的数据，它只会创建一个指向原数组的切片结构体，所以修改新切片的数据也会修改原切片。

## slice底层
```
type SliceHeader struct {
	Data uintptr
	Len  int
	Cap  int
}
```
Pointer 是指向一个数组的指针，len 代表当前切片的长度，cap 是当前切片的容量。cap 总是大于等于 len 的。

如果想从 slice 中得到一块内存地址，可以这样做：
```
s := make([]byte, 200)
ptr := unsafe.Pointer(&s)
```

当然还有更加直接的方法，在 Go 的反射中就存在一个与之对应的数据结构 SliceHeader，我们可以用它来构造一个 slice
```
package main

import (
    "fmt"
    "reflect"
    "unsafe"
)

func main() {

    s := make([]byte, 10, 20)       //s为引用类型
    fmt.Println(unsafe.Pointer(&s)) //取出s在内存中的地址

    sh := (*reflect.SliceHeader)(unsafe.Pointer(&s)) //将改地址转换为sliceHeader类型的指针

    var o []byte
    sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&o)))
    sliceHeader.Data = sh.Data
    sliceHeader.Len = sh.Len
    sliceHeader.Cap = sh.Cap

    var e string = "aaass"
    stringHeader := (*reflect.StringHeader)(unsafe.Pointer(&e))
    sliceHeader.Data = stringHeader.Data
    sliceHeader.Len = stringHeader.Len
    sliceHeader.Cap = stringHeader.Len * 2
    fmt.Println(o)
    fmt.Println(e)

    var f string
    stringHeader = (*reflect.StringHeader)(unsafe.Pointer(&f))
    stringHeader.Data = sliceHeader.Data
    stringHeader.Len = sliceHeader.Len
    fmt.Println(f)

}
```

## slice创建
slice是有最大容量的，Slice 所允许申请的最大容量大小，与当前**值类型**和当前**平台位数**有直接关系
https://studygolang.com/articles/17484


## 扩容
当前slice的数组容量够的时候，直接在当前slice上增加元素。
```
array := [4]int{10, 20, 30, 40}
slice := array[0:2]
newSlice := append(slice, 50)
```
当前slice的数组容量不够的时候，计算新的容量，并复制内容过去。

## Copy
存在内存拷贝
![enter image description here](/tencent/api/attachments/s3/url?attachmentid=4872800)

## 空slice和nil slice
空slice data 指针指向一个固定的值
nil slice data 指针为nil

## 参考
https://halfrost.com/go_slice/#toc-1




## Slice 深坑
1. https://www.zhihu.com/question/27161493?sort=created 二维slice，append同一个在变化的一维slice对象

``` 
	data := []int{1, 2, 3, 4}
    res := [][]int{}
    for i := 100; i < 103; i++ {
        res = append(res, data)
        data[0] = i

    }
    fmt.Println(res) 
// [[102 2 3 4] [102 2 3 4] [102 2 3 4]] 
data没有发生扩容，data一直指向同一个底层数组。
res中的每一个元素都指向data指向的底层数组。
data变化，将影响res中每一个元素。
```
避免方法：凡是二维数组append，最好有使用 res = append(res, append([]int{}, data...))  这种形式

``` 
data := []int{1, 2, 3, 4}
    res := [][]int{}
    for i := 100; i < 103; i++ {
        res = append(res, append([]int{}, data...)) 
		// 每次新建一个slice，用于append。 
		// 1. 将data中的元素追加拷贝到新的[]int{}的底层数组后面。即使data改变，[]int{}也不会改变。
		// 2. res元素指向[】int{}的底层数组。只要不改变[]int{}, res的元素就不会变化
        data[0] = i

    }
    fmt.Println(res)
```

2. https://www.zhihu.com/question/27161493?sort=created 如果你 append 过后的结果没有 assign 回原来的 slice 变量，这种用法常常是错的。

```
	s := []int{5}
    s = append(s, 7)
    s = append(s, 9)
    x := append(s, 11)
    y := append(s, 12)
   fmt.Println(s, x, y) //`[5 7 9] [5 7 9 12] [5 7 9 12]`
   
   s x y 全部指向同一个底层数组
```
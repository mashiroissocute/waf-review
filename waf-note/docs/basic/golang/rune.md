## ASCII 、 Unicode 、Utf-8
**ASCII 和 Unicode都是字符集。**
所谓字符集，是指每个字符都对应唯一的ID。
ASCII中只有英文数字等字符，并不包含中文等其他国家的字符。
Unicode包含了中文，日文等其他国家的字符。并且兼容了ASCII。

ASCII和Unicode等字符集只是指定了字符和ID的对应关系。但是并没有指明如何对ID进行计算机编码。
Utf-8等编码格式，指定了编码方式。utf8是可变长度字符编码，不同的字符会对应不同大小的存储方式，**全存相同大小字节，浪费空间**。比如"a"字符(unicode值97)用1个字节，而"中"字符(unicode值20013）则用3个字节。字符的unicode值决定了字符需要用多少字节表示。


## byte 、rune 、string
``` 
type byte uint8 单字节
type rune int32 4个字节
string 是byte切片
```
byte其实是char的概念，string是byte切片，rune是4个字节。

go采用utf-8编码格式，最长的字符ID占四个字节，最短的字符ID占一个字节。中文的字符ID占三个字节。

因此使用rune（4个字节）来表示utf-8编码格式，是合适的选择。
``` 
func main(){
    str := "名称Tom"
    fmt.Println(len(str))              # 9
    fmt.Println(len([]rune(str)))      # 5
}
```
``` 
//string 转[]byte
b := []byte(str)
//[]byte转string
str = string(b)
//string 转 rune
r := []rune(str)
//rune 转 string
str = string(r)
```



## 转换原理
例如 str := "名称Tom"
即 什么时候该读3个字节以表示1个字符，什么时候该读1个字节以表示字符？
**根据Unicode符号范围**
```
Unicode符号范围      | UTF-8编码方式
(十六进制)           | （二进制）
--------------------+---------------------------------------------
0000 0000-0000 007F | 0xxxxxxx
0000 0080-0000 07FF | 110xxxxx 10xxxxxx
0000 0800-0000 FFFF | 1110xxxx 10xxxxxx 10xxxxxx
0001 0000-0010 FFFF | 11110xxx 10xxxxxx 10xxxxxx 10xxxxxx
```
例如 `张`  字，unicode编码 `5F20` ，对应的十六进制处于 `0000 0800-0000 FFFF` 中，也就是 `3` 个字节。


## 参考
http://www.randyfield.cn/post/2022-01-14-rune-unicode-utf8/
https://juejin.cn/post/6844903743524175879
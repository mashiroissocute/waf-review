
## list
![alt text](image.png)

### Element

``` 
type Element struct {
    next, prev *Element
    // The list to which this element belongs.
    list *List
    // The value stored with this element.
    Value interface{} //Value是inerface，因此可以存多种值，但是在取出来的时候要用switch判断该值的类型。
}

// Next returns the next list element or nil.
func (e *Element) Next() *Element {
    if p := e.next; e.list != nil &amp;&amp; p != &amp;e.list.root {
        return p
    }
    return nil
}

// Prev returns the previous list element or nil.
func (e *Element) Prev() *Element {
    if p := e.prev; e.list != nil &amp;&amp; p != &amp;e.list.root {
        return p
    }
    return nil
}
```
### List
``` 
type List struct {
    root Element // 头节点 sentinel list element, only &root, root.prev, and root.next are used
    len  int     // current list length excluding (this) sentinel element
}
```
###  Method
``` 
package main

import (
    "container/list"
    "fmt"
)

func main() {
    var l = list.New()
    //push back or front
    e1 := l.PushBack("33")
    e2 := l.PushFront(1)
    //insert back or front
    e3 := l.InsertBefore(make(map[int]interface{}), e2)
    e4 := l.InsertAfter(make(map[int]interface{}), e1)
    //remove
    v3 := l.Remove(e4)
    switch vv := v3.(type) {  //switch断言+转换类型
    case string:
        fmt.Println(vv)
    case int:
        fmt.Println(vv)
    default:
        fmt.Println("map")
    }
    //move back or front , before or after element
    l.MoveToBack(e3)
    l.MoveToFront(e2)
    l.MoveAfter(e2, e3)
    l.MoveBefore(e2, e3)
    //push front or back otherlist
    l.PushFrontList(l)
    l.PushBackList(l)

}
```
### 遍历
``` 
//这里只能把root（头节点）理解为nil
for e := l.Front(); e != nil; e = e.Next() {
        fmt.Println(e.Value)
    }
for e := l.Back(); e != nil; e = e.Prev() {
        fmt.Println(e.Value)
    }
```
### 资料
- https://ijayer.github.io/post/tech/code/golang/tutorial-go36-03/ （第三方）
- https://pkg.go.dev/container/list （官方）
- https://cs.opensource.google/go/go/+/refs/tags/go1.18:src/container/list/list.go （源码）


### Notice
switch 默认情况下 case 最后自带 break 语句，匹配成功后就不会执行其他 case，如果我们需要执行后面的 case，可以使用 fallthrough
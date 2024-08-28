## 值得做
- https://leetcode.cn/problems/longest-palindromic-substring/

## 快慢指针
一般看到有序数组 寻找target 使用双指针
### 删除有序数组重复项
```golang
func ...() {
    slow,fast := 0,0
    for fast != len(array) {
        if array[slow] != array[fast] {
            slow++
            array[slow],array[fast] = array[fast],array[slow] //swap 
            //or 
            array[slow]= array[fast] //revalue
        }
        fast++
    }
}
```

### 原地删除元素
```golang
func ...(){
    slow,fast := 0,0
    for fast != len(array){
        if array[fast] != val {
            array[slow] = arrat[fast]
            slow++
        }

        fast++
    }
    return slow // 返回数组长度
}
```


### 移动0
```golang
func ...(){
    slow,fast := 0,0
    for fast != len(array){
        if array[fast] != 0 {
            array[slow],array[fast] = array[fast],array[slow]
            slow++
        }
        fast++
    }

}
```


## 左右指针
一般看到有序数组 寻找target 使用双指针

### 普通左右指针
#### 两数之和
```golang
func ...(){
    left，right := 0, len(array) -1
    for left < right { //要求不重复，用小于
        if array[left] + array[right] == target{
            return left,right
        }else if array[left] + array[right] < target{
            left++
        }else if array[left] + array[right] > target{
            right--
        }
    }
    return -1,-1
}
```

#### 反转数组
```golang
func ...(){
    left，right := 0, len(array) -1
    for left < right { //要求不重复，用小于
       array[left],array[right]=array[right],array[left]
       left++
       right--
    }
}
```

#### 最长回文
```golang



func 寻找回文start和end(s string, left,right int) (left,right int){

    for left >= 0 && right <= len(array)-1 {
        if array[left] == array[right] {
            left--
            right++
        }else{
            beak
        }
    }
    return left, right
}

func ...(){
    longest := math.MinInt
    // 奇数
    for i:=0;i<len(array);i++{
        left,right := 寻找回文start和end(array,i,i)
        if right - left + 1 > longest {
            longest = right - left + 1 
        }
    }

    // 偶数
    for i:=0;i<len(array);i++{
        left,right := 寻找回文start和end(array,i,i+1)
        if right - left + 1 > longest {
            longest = right - left + 1 
        }
    }


}




```

### 二分左右指针

#### 二分查找
```golang

func ...(){
    left,right := 0,len(array)-1

    for left <= right {
        mid := left + (right-left) /2 
        // 不要用else 要罗列所有情况
        if array[mid] == target {
            return mid // 只要满足就return
        }else if array[mid] < target{
            left = mid + 1 
        }else if array[mid] > target{
            right = mid -1 
        }
    }

    // 固定结局
    return -1
}

```

#### 更多的二分算法
[二分算法](./binaryalgo.md)



## 滑动窗口
[滑动窗口](./slidingwindow.md)
## 值的做
- https://leetcode.cn/problems/minimum-window-substring/description/


滑动窗口一般用于搜索一个满足条件的连续区间

## 框架
```golang
// 滑动窗口算法框架
func slidingWindow(s string) {
    // 用合适的数据结构记录窗口中的数据
    window := make(map[byte]int)

    left, right := 0, 0

    // 不断地增加 right 指针扩大窗口 ，直到窗口中的数据符合要求
    for right < len(s) {
        c := s[right]
        right++

        // 进行窗口内数据的一系列更新
        window[c]++

    
        fmt.Printf("window: [%d, %d)\n", left, right)


        // 停止增加 right，转而不断增加 left 指针缩小窗口，直到窗口中的字符串不再符合要求
        for left < right && window needs shrink { // 判断什么时候开始收紧窗口，十分重要
            //在这里更新结果

            d := s[left]
            left++
            // 进行窗口内数据的一系列更新
            window[d]--
        }
    }
}
```

## 最小覆盖字符串
```golang
func ...(){
    need := make(map[byte]int)
    window := make(map[byte]int)
    vaild := 0
    start := 0
    length := math.MaxInt

    for _,c:=range t {
        need[c]++
    }
    left,right:=0,0
    
    for right < len(s) {
        //扩大窗口
        c := s[right]
        right++
        
        if _,ok := need[c]; ok{
            window[c]++
            if window[c] == need[c] {
                vaild++
            }
        }

        for valid == len(need) { //当满足条件的时候开始缩小窗口
            //更新长度
            if right - left + 1 < length {
                start = left
                length = right - left + 1
            }
            //缩小窗口
            c := s[left]
            left++

            if _,ok := needp[c]; ok {
                window[c]--
                if window[c] < need[c] {
                    valid--
                }
            }
        }
    }

    if length == math.MaxInt {
        return  -1,-1
    }

    return start, start+length-1
}
```


## 字符串的排列
```golang
func ...(){
    need := make(map[string]int)
    window := make(map[string]int)
    vaild := 0
    left,right := 0,0

    for _,c:=range t{
        need[c]++
    }

    for right < len(s) {
        c := s[right]
        right++

        if _,ok:=need[c]; ok {
            window[c]++
            if window[c] == need[c]{
                vaild++
            }
        }

        for right-left == len(t) { // 判断什么时候开始收紧窗口，十分重要
            if vaild == len(need){
                return true
            }
            c := s[left]
            left++
            if _,ok:=need[c]; ok {
                window[c]--
                if window[c] < need[c]{
                    vaild--
                }
            }
        }
    }
    return false
}
```

## 找所有字母异位词
和上一题一样
```golang
func ...(){
    need := make(map[string]int)
    window := make(map[string]int)
    vaild := 0
    left,right := 0,0
    pos := make([]int,4)

    for _,c:=range t{
        need[c]++
    }

    for right < len(s) {
        c := s[right]
        right++

        if _,ok:=need[c]; ok {
            window[c]++
            if window[c] == need[c]{
                vaild++
            }
        }

        for right-left == len(t) { // 判断什么时候开始收紧窗口，十分重要
            if vaild == len(need){
                pos = append(pos,left)
            }
            c := s[left]
            left++
            if _,ok:=need[c]; ok {
                window[c]--
                if window[c] < need[c]{
                    vaild--
                }
            }
        }
    }
    return pos
}
```


## 最长无重复子串
```golang
func ...(){
    window := make(map[byte]int)
    maxLength := math.MinInt

    left,right := 0,0
    for right < len(s) {
        c := s[right]
        right++
        window[c]++

        for window[c] > 1 {
            if right - left + 1 > maxLength {
                maxLength = right - left + 1
            }

            c := s[left]
            left--
            window[c]--
        }
    }

    if maxLength == math.MinInt {
        maxLength = len(s)
    }
    return maxLength
}

```
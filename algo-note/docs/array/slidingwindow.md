## 值的做
- https://leetcode.cn/problems/minimum-window-substring/description/


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
        for left < right && window needs shrink {
            d := s[left]
            left++

            // 进行窗口内数据的一系列更新
            window[d]--
        }
    }
}
```


## 值的做


## DP三要素

- 1.列出正确的「状态转移方程」,即描述问题结构的数学形式
- 2.判断算法问题是否具备「最优子结构」，是否能够通过子问题的最值得到原问题的最值。
- 3.如果动态规划问题存在「重叠子问题」，需要使用「DP table」来优化穷举过程，避免不必要的计算。

重叠子问题、最优子结构、状态转移方程就是动态规划三要素。

## 解题框架
```golang
// 自底向上迭代的动态规划
// 初始化 base case
dp[0][0][...] = base case
// 进行状态转移
for 状态1 in 状态1的所有取值：
    for 状态2 in 状态2的所有取值：
        for ...
            dp[状态1][状态2][...] = 求最值(选择1，选择2...)
```

## 零钱兑换
```golang

func coinChange(coins []int, amount int) int {
    dp := make([]int, amount+1)
    // 数组大小为 amount + 1，初始值也为 amount + 1
    for i := 0; i < len(dp); i++ {
        dp[i] = amount + 1
    }

    // base case
    dp[0] = 0
    // 外层 for 循环在遍历所有状态的所有取值
    for i := 0; i < len(dp); i++ {
        // 内层 for 循环在求所有选择的最小值
        for _, coin := range coins {
            // 子问题无解，跳过
            if i-coin < 0 {
                continue
            }
            dp[i] = min(dp[i], dp[i-coin]+1)

        }
    }

    if dp[amount] == amount+1 {
        return -1
    }
    return dp[amount]
}

```
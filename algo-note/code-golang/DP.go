package main

import "math"

func coinChange(coins []int, amount int) int {
	dp := make([]int, amount+1)
	dp[0] = 0

	for i := 1; i < len(dp); i++ {

		tempRes := math.MaxInt
		for _, coin := range coins {
			if i-coin < 0 {
				continue
			}

			if dp[i-coin] == math.MaxInt {
				continue
			}

			if tempRes > dp[i-coin]+1 {
				tempRes = dp[i-coin] + 1
			}
		}
		dp[i] = tempRes

	}

	if dp[amount] == math.MaxInt {
		return -1
	}
	return dp[amount]
}

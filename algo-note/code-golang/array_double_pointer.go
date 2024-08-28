package main

func longestPalindrome(s string) string {
	res := ""
	for i := 0; i < len(s); i++ {
		evenC := Palindrome(s, i, i) //奇数回文
		if len(evenC) > len(res) {
			res = evenC
		}
	}

	for i := 0; i < len(s)-1; i++ {
		oddC := Palindrome(s, i, i+1) //偶数回文
		if len(oddC) > len(res) {
			res = oddC
		}
	}
	return res
}

func Palindrome(s string, l, r int) string {
	for l >= 0 && r < len(s) {
		if s[l] == s[r] {
			l--
			r++
		} else {
			break
		}
	}
	return s[l+1 : r]
}

func twoSum(numbers []int, target int) []int {
	left, right := 0, len(numbers)-1
	for left < right {
		sum := numbers[left] + numbers[right]
		if sum == target {
			return []int{left + 1, right + 1}
		} else if sum < target {
			left++
		} else {
			right--
		}
	}
	return []int{-1, -1}
}

func moveZeroes(nums []int) {
	slow, fast := 0, 0
	for fast < len(nums) {
		if nums[fast] != 0 {
			nums[slow] = nums[fast]
			slow++
		}
		fast++
	}
	for slow < len(nums) {
		nums[slow] = 0
		slow++
	}
}





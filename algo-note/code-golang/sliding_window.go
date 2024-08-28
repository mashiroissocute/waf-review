package main

import "math"

func lengthOfLongestSubstring(s string) int {
	windowMap := make(map[byte]int)

	maxLen := math.MinInt
	left, right := 0, 0

	for right < len(s) {
		c := s[right]
		right++
		windowMap[c]++

		for windowMap[c] > 1 {
			windowMap[s[left]]--
			left++
		}

		if right-left > maxLen {
			maxLen = right - left
		}

	}

	if maxLen == math.MinInt {
		return 0
	}
	return maxLen
}

func findAnagrams(s string, p string) []int {
	pMap := make(map[byte]int)
	for i := range p {
		pMap[p[i]]++
	}

	hitCount := 0
	windownMap := make(map[byte]int)

	left, right := 0, 0
	res := make([]int, 0)

	for right < len(s) {
		c := s[right]
		if _, ok := pMap[c]; ok {
			windownMap[c]++
			if windownMap[c] == pMap[c] {
				hitCount++
			}
		}
		right++

		for right-left >= len(p) {
			if hitCount == len(pMap) {
				res = append(res, left)
			}

			c := s[left]
			if _, ok := pMap[c]; ok {
				if windownMap[c] == pMap[c] {
					hitCount--
				}
				windownMap[c]--
			}
			left++
		}
	}

	return res
}

func checkInclusion(s1 string, s2 string) bool {
	// fmt.Println(len(s1), len(s2))
	s1, s2 = s2, s1

	s2Map := make(map[byte]int)
	for i := range s2 {
		s2Map[s2[i]]++
	}

	left, right := 0, 0
	windowMap := make(map[byte]int)
	hitCount := 0

	for right < len(s1) {
		c := s1[right]
		if _, ok := s2Map[c]; ok {
			windowMap[c]++
			if windowMap[c] == s2Map[c] {
				hitCount++
			}
		}
		right++

		// fmt.Println("window [%d,%d)", left, right, s1[left:right])

		for right-left >= len(s2) {
			if hitCount == len(s2Map) {
				return true
			}

			c := s1[left]
			if _, ok := s2Map[c]; ok {
				if windowMap[c] == s2Map[c] {
					hitCount--
				}
				windowMap[c]--
			}
			left++

		}
	}
	return false
}

func minWindow(s string, t string) string {
	tMap := make(map[byte]int) // target Map
	for i := range t {
		tMap[t[i]]++
	}

	windowMap := make(map[byte]int) // sliding window Map
	hitCount := 0

	left, right := 0, 0
	minString := ""
	minLen := math.MaxInt
	for right < len(s) {
		// 扩张window的时候，干些什么
		c := s[right]
		if _, ok := tMap[c]; ok {
			windowMap[c]++
			if windowMap[c] == tMap[c] {
				hitCount++
			}
		}
		right++

		// debug的位置
		// fmt.Println("window [%d,%d), hitCount %d", left, right, hitCount)

		for hitCount == len(tMap) { // 设置开始收缩window的条件，不断去搜索window
			if right-left < minLen { // 每次收缩window时，都要更新结果
				minLen = right - left
				minString = s[left:right]
			}

			// 收缩window的时候，干些什么
			c := s[left]
			if _, ok := tMap[c]; ok {
				windowMap[c]--
				if windowMap[c] < tMap[c] {
					hitCount--
				}

			}
			left++
		}
	}

	return minString
}

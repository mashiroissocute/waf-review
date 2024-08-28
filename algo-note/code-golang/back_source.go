package main

func solveNQueens(n int) [][]string {

	res := make([][]string, 0)

	board := make([][]byte, n)

	for i := range board {
		board[i] = make([]byte, n)
	}

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			board[i][j] = '.'
		}
	}

	valid := func(row, col int) bool {
		// check  col
		for i := 0; i < row; i++ {
			if board[i][col] == 'Q' {
				return false
			}
		}

		// check upLeft
		for i, j := row-1, col-1; i >= 0 && j >= 0; i, j = i-1, j-1 {
			if board[i][j] == 'Q' {
				return false
			}
		}

		// check upRight
		for i, j := row-1, col+1; i >= 0 && j < n; i, j = i-1, j+1 {
			if board[i][j] == 'Q' {
				return false
			}
		}

		return true
	}

	var backSource func(int)
	backSource = func(row int) {
		if row == n {
			tempRes := make([]string, 0)

			for _, line := range board {
				tempRes = append(tempRes, string(line))
			}

			res = append(res, tempRes)
			return // 这个return非常重要
		}

		for col := 0; col < n; col++ {
			if !valid(row, col) {
				continue
			}

			board[row][col] = 'Q'

			backSource(row + 1)

			board[row][col] = '.'

		}
	}

	row := 0
	backSource(row)
	return res
}

func permute(nums []int) [][]int {

	res := make([][]int, 0)

	var permuteHelper func([]int, []int)
	permuteHelper = func(nums []int, path []int) {
		if len(nums) == 0 {
			res = append(res, append([]int{}, path...))
			return
		}

		for i := 0; i < len(nums); i++ {
			path = append(path, nums[i])

			newNums := make([]int, 0)
			newNums = append(newNums, nums[:i]...)
			newNums = append(newNums, nums[i+1:]...)
			permuteHelper(newNums, path)

			path = path[:len(path)-1]
		}
	}

	path := make([]int, 0)
	permuteHelper(nums, path)

	return res

}

func permute1(nums []int) [][]int {

	res := make([][]int, 0)
	used := make([]bool, len(nums))

	var permuteHelper func([]int)
	permuteHelper = func(path []int) {
		if len(path) == len(nums) {
			res = append(res, append([]int{}, path...))
			return
		}

		for i := 0; i < len(nums); i++ {
			if used[i] {
				continue
			}

			used[i] = true
			path = append(path, nums[i])

			permuteHelper(path)

			used[i] = false
			path = path[:len(path)-1]
		}
	}

	path := make([]int, 0)
	permuteHelper(path)

	return res

}



package main

// BFS
func minDepth1(root *TreeNode) int {
	q := make([]*TreeNode, 0)

	if root == nil {
		return 0
	}

	q = append(q, root)
	level := 1

	for len(q) != 0 {
		l := len(q)
		for i := 0; i < l; i++ {
			curNode := q[0]
			q = q[1:]

			if curNode.Left == nil && curNode.Right == nil {
				return level
			}

			if curNode.Left != nil {
				q = append(q, curNode.Left)
			}
			if curNode.Right != nil {
				q = append(q, curNode.Right)
			}

		}
		level++

	}
	return level

}

func minDepth(root *TreeNode) int {
	return minDepthHelper(root)
}

func minDepthHelper(root *TreeNode) int {
	if root == nil {
		return 0
	}

	leftLevel := minDepthHelper(root.Left)
	rightLevel := minDepthHelper(root.Right)

	if leftLevel == 0 {
		return rightLevel + 1
	}

	if rightLevel == 0 {
		return leftLevel + 1
	}

	if leftLevel > rightLevel {
		return rightLevel + 1
	}

	return leftLevel + 1

}

func openLock(deadends []string, target string) int {

	visited := make(map[string]bool)
	dead := make(map[string]bool)

	for _, v := range deadends {
		dead[v] = true
	}

	if dead["0000"] {
		return -1
	}

	q := make([]string, 0)
	q = append(q, "0000")
	visited["0000"] = true

	step := 0

	for len(q) != 0 {
		sz := len(q)
		// fmt.Println(q)

		for i := 0; i < sz; i++ {
			curLock := q[0] //这个一定是0
			q = q[1:]

			if curLock == target {
				return step
			}

			for j := 0; j < 4; j++ {
				nextLockUp := plusOne(curLock, j)
				if _, ok := visited[nextLockUp]; !ok {

					if _, ok := dead[nextLockUp]; !ok {
						q = append(q, nextLockUp)
						visited[nextLockUp] = true

					}

				}

				nextLockDown := minusOne(curLock, j)
				if _, ok := visited[nextLockDown]; !ok {

					if _, ok := dead[nextLockDown]; !ok {
						q = append(q, nextLockDown)
						visited[nextLockDown] = true

					}
				}

			}

		}
		step++

	}

	return -1

}

func plusOne(s string, j int) string {
	arr := []byte(s)
	if arr[j] == '9' {
		arr[j] = '0'
	} else {
		arr[j]++
	}
	return string(arr)
}

func minusOne(s string, j int) string {
	arr := []byte(s)
	if arr[j] == '0' {
		arr[j] = '9'
	} else {
		arr[j]--
	}
	return string(arr)
}

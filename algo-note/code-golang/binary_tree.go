package main

/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */
func diameterOfBinaryTree(root *TreeNode) int {

	diameter := 0

	var depthfunc func(*TreeNode) int
	depthfunc = func(root *TreeNode) int {
		if root == nil {
			return 0
		}

		leftDepth := depthfunc(root.Left)
		rightDepth := depthfunc(root.Right)

		if leftDepth+rightDepth+1 > diameter {
			diameter = leftDepth + rightDepth + 1
		}
		return max(leftDepth, rightDepth) + 1
	}

	depthfunc(root)
	return diameter
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func maxDepth(root *TreeNode) int {
	return maxDepthHelper(root)
}

func maxDepthHelper(root *TreeNode) int {
	if root == nil {
		return 0
	}

	leftLevel := maxDepthHelper(root.Left)
	rightLevel := maxDepthHelper(root.Right)

	if leftLevel > rightLevel {
		return leftLevel + 1
	}

	return rightLevel + 1

}

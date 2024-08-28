package main

func searchRange(nums []int, target int) []int {
	return []int{searchLeft(nums, target), searchRight(nums, target)}
}

func searchLeft(nums []int, target int) int {
	left, right := 0, len(nums)-1
	for left <= right {
		mid := (right + left) / 2
		if nums[mid] == target {
			if mid == 0 || nums[mid-1] != target {
				return mid
			} else {
				right = mid - 1
			}
		} else if nums[mid] < target {
			left = mid + 1
		} else if nums[mid] > target {
			right = mid - 1
		}
	}
	return -1
}

func searchRight(nums []int, target int) int {
	left, right := 0, len(nums)-1
	for left <= right {
		mid := (right + left) / 2
		if nums[mid] == target {
			if mid == len(nums)-1 || nums[mid+1] != target {
				return mid
			} else {
				left = mid + 1
			}
		} else if nums[mid] < target {
			left = mid + 1
		} else if nums[mid] > target {
			right = mid - 1
		}
	}
	return -1
}




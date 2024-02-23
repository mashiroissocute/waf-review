## 值的做
- https://leetcode.cn/problems/find-first-and-last-position-of-element-in-sorted-array/description/

## 二分查找
```golang
func binarySearch(nums []int, target int) int {
    left, right := 0, len(nums)-1
    for left <= right {
        mid := left + (right - left) / 2
        if nums[mid] < target {
            left = mid + 1
        } else if nums[mid] > target {
            right = mid - 1
        } else if nums[mid] == target {
            // 直接返回
            return mid
        }
    }
    // 直接返回
    return -1
}
```

## 二分查找左边界
```golang
func leftBound(nums []int, target int) int {
    left, right := 0, len(nums)-1
    for left <= right {
        mid := left + (right - left) / 2
        if nums[mid] < target {
            left = mid + 1
        } else if nums[mid] > target {
            right = mid - 1
        } else if nums[mid] == target {
            // 别返回，锁定左侧边界
            right = mid - 1
        }
    }
    // 判断 target 是否存在于 nums 中
    if left < 0 || left >= len(nums) {
        return -1
    }
    // 判断一下 nums[left] 是不是 target
    if nums[left] == target {
        return left
    }
    return -1
}
```

## 二分查找右边界
```golang
func rightBound(nums []int, target int) int {
    left, right := 0, len(nums)-1
    for left <= right {
        mid := left + (right - left) / 2
        if nums[mid] < target {
            left = mid + 1
        } else if nums[mid] > target {
            right = mid - 1
        } else if nums[mid] == target {
            left = mid + 1
        }
    }
    if right < 0 || right >= len(nums) {
        return -1
    }
    if nums[right] == target {
        return right
    }
    return -1
}
```


## 在排序数组中查找元素的第一个和最后一个位置
```golang
func searchRange(nums []int, target int) []int {
	return []int{searchLeft(nums, target), searchRight(nums, target)}
}


func searchLeft(nums []int, target int) int {
    left,right := 0, len(nums)-1


    for left <= right {
        mid := left + (right-left)/2
        if nums[mid] == target {
            right = mid - 1
        }else if nums[mid] > target {
            right = mid - 1
        }else if nums[mid] < target {
            left = mid + 1
        }
    }

    if left < 0 || left > len(nums)-1{
        return -1 
    }

    if nums[left] == target {
        return left
    }

    return -1

}


func searchRight(nums []int ,target int) int {

    left,right := 0,len(nums)-1

    for left <= right {
        mid := left + (right-left)/2
        if nums[mid] == target {
            left = mid +1
        }else if nums[mid] < taregt {
            left = mid +1
        }else if nums[mid] > target {
            right = mid -1 
        }
    }

    if right <0 || right >len(nums)-1{
        return -1 
    }

    if nums[right] == target {
        return right
    }
    
    return -1
}

```
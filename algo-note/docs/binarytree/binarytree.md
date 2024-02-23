## 值的做

## 二叉树解题的思维模式分两类：

- 1、是否可以通过遍历一遍二叉树得到答案？如果可以，用一个 traverse 函数配合外部变量来实现，这叫「遍历」的思维模式。

- 2、是否可以定义一个递归函数，通过子问题（子树）的答案推导出原问题的答案？如果可以，写出这个递归函数的定义，并充分利用这个函数的返回值，这叫「分解问题」的思维模式。

二叉树的所有问题，就是在前中后序位置注入逻辑，去达到自己的目的，你只需要单独思考每一个节点应该做什么，其他的不用你管，抛给二叉树遍历框架，递归会在所有节点上做相同的操作。


## 快速排序
二叉树前序遍历


```golang

func partition(nums []int, left,right int) int {
    povit := nums[left]

    for left < right {
        for left < right && nums[right] >= povit {
            right--
        }
        nums[right] = nums[left]
        for left < right && nums[left] <= povit {
            left++
        }
        nums[left] = nums[right]
    }

    nums[left] = povit
    return left
}

func quickSort(nums []int, left,right int) {

    if left == right {
        return 
    }

    pos := partition(nums, left ,right)
    
    quickSort(nums, left, pos - 1)
    quickSort(nums, pos+1, right)

}

```

## 归并排序
二叉树后序遍历

```golang
func merge(nums1 []int, nums2 []int) []nums{
    res := make([]int, len(num1)+len(nums2))
    
    p1 ,p2 := 0,0
    r1 := 0
    for p1 < len(nums1) && p2 < len(nums2) {
        if nums1[p1] <= nums2[p2] {
            res[r1] = nums[p1]
            r1++
            p1++
        }else{
            res[r1] = nums2[p2]
            r1++
            p2++    
        }
    }   


    for p1 <len(nums1) {
        res[r1] = nums[p1]
        r1++
        p1++
    }

    for p2 <len(nums2) {
        res[r1] = nums2[p2]
        r1++
        p2++  
    }
    
    return res
}


func mergeSort(nums []int) {
    if len(nums) == 1 {
        return nums
    }

    mid := len(nums)/2
    n1 := mergeSort(nums[:mid])
    n2 := mergeSort(nums[mid:])
    return merge(n1,n2)
}
```



## 数组前后序遍历
```golang
func preorderTravel(nums []int, i int) {
    if i == len(nums) {
        return 
    }

    do(nums[i])
    preorderTravel(nums, i+1)

}

func postorderTravel(nums []int, i int) {
    if i == len(nums) {
        return 
    }

    preorderTravel(nums, i+1)
    do(nums[i])

}
```

## 链表前后序遍历
```golang
func preorderTravel(root *listNode) {
    if root == nil {
        return
    }

    do(root)
    preorderTravel(root.next)
}

func postorderTravel(root *listNode){
    if root == nil {
        return 
    }

    postorderTravel(root.next)
    do(root)
}
```


## 二叉树前中后序遍历

```golang

type struct TreeNode{
    val int,
    left *TreeNode,
    right *TreeNode
}

func preorderTravel(root *TreeNode) {
    if root == nil {
        return 
    }

    do(root)
    preTravel(root.left)
    preTravel(root.right)
}

func inorderTravel(root *TreeNode) {
    if root == nil {
        return 
    }

    inorderTravel(root.left)
    do(root)
    inorderTravel(root.right)
}


func postorderTravel(root *TreeNode) {
    if root == nil {
        return 
    }
    postorderTravel(root.left)
    postorderTravel(root.right)
    do(root)
}



```
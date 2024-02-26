## 值的做

## 二叉树解题的思维模式分两类：

- 1、是否可以通过遍历一遍二叉树得到答案？如果可以，用一个 traverse 函数配合外部变量来实现，这叫「遍历」的思维模式。

- 2、是否可以定义一个递归函数，通过子问题（子树）的答案推导出原问题的答案？如果可以，写出这个递归函数的定义，并充分利用这个函数的返回值，这叫「分解问题」的思维模式。

二叉树的所有问题，就是在前中后序位置注入逻辑，去达到自己的目的，你只需要单独思考每一个节点应该做什么，其他的不用你管，抛给二叉树遍历框架，递归会在所有节点上做相同的操作。



## 函数写法
二叉树题目的递归解法可以分两类思路，第一类是遍历一遍二叉树得出答案，第二类是通过分解问题计算出答案，这两类思路分别对应着 回溯算法核心框架 和 动态规划核心框架。

二叉树中用遍历思路解题时函数签名一般是 void traverse(...)，没有返回值，靠更新外部变量来计算结果

而分解问题思路解题时函数名根据该函数具体功能而定，而且一般会有返回值，返回值是子问题的计算结果。

与此对应的，回溯算法核心框架 中给出的函数签名一般也是没有返回值的 void backtrack(...)，而在 动态规划核心框架 中给出的函数签名是带有返回值的 dp 函数。这也说明它俩和二叉树之间千丝万缕的联系。


## 前中后序的注意点
- 中序用于BST更多

- 前序位置的代码只能从函数参数中获取父节点传递来的数据。 例如`如果把根节点看做第 1 层，如何打印出每一个节点所在的层数？`

- 后序位置的代码不仅可以获取参数数据，还可以获取到子树通过函数返回值传递回来的数据。`如何打印出每个节点的左右子树各有多少节点？`




## 二叉树和DFS/BFX、回溯、动态规划
动归/DFS/回溯算法都可以看做二叉树问题的扩展，只是它们的关注点不同：

- 动态规划算法属于分解问题的思路，它的关注点在整棵「子树」。
- 回溯算法属于遍历的思路，它的关注点在节点间的「树枝」。
- DFS/BFS算法属于遍历的思路，它的关注点在单个「节点」。DFS对应前序遍历，BFS对应层序遍历。













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

## 二叉树的最大深度

```golang
func maxDeep(root *TreeNode) int {
    if root == nil {
        return 0
    }

    leftDeep := maxDeep(root.left)
    rightDeep := maxDeep(root.right)
    
    return max(leftDeep,rightDeep) + 1

}
```

## 打印每一个节点所在的层数
```golang

func printTree(root *TreeNode) {
    printLevel(root,1)
}

func printLevel(root *TreeNode, level int) {
    if root == nil {
        return 
    }

    print(root.val, level)
    printLevel(root.left, level+1)
    printLevel(root.right, level+1)
}

```


## 打印节点的左右子树各有多少节点
```golang

func printTree(root *TreeNode) {
   printNum(root)
}

func printNum(root *TreeNode) int{
    if root == nil {
        return 0
    }

    leftNum := printNum(root.left)
    rightNUm := printNum(root.right)

    print(root.val,leftNum)
    print(root.val,rightNum)
    return leftNum+rightNum+1
}
```


## 二叉树的直径
```golang
func ...() {
    maxDia := math.MinInt

    var diameterOfBinaryTree func (root *TreeNode) int
    diameterOfBinaryTree = (root *TreeNode) int{
        if root == nil {
            return 0 
        }

        leftLevel := diameterOfBinaryTree(root.left)
        rightLevel := diameterOfBinaryTree(root.right)
        if max(leftLevel,rightLevel) > maxDia {
            maxDia = max(leftLevel,rightLevel)
        }
        return max(leftLevel,rightLevel) + 1
    }

    diameterOfBinaryTree(root)
    return max
}
```

## 二叉树层序遍历

```golang
// list
import container/list
func ...() {
    l := list.Init()
    if root != nil {
        l.PushBack(root)
    }


    for l.Len() != 0 {
        cur := l.Front()
        print(cur.(*TreeNode))
        l.remove(cur)
        
        if cur.left !=nil {
            l.PushBack(cur.left)
        }

        if cur.right != nil {
            l.PushRight(cur.right)
        }

    }

}
```

```golang
// array 
func ...(){
    q := make([]*TreeNode,0)
    if root != nil{
        q = append(q, root)
    }

    for len(q) != 0 {
        sz := len(q)
        for i:=0;i<sz;i++{
            cur := q[0]
            q := q[1:]
            print(cur.val)

            if cur.left != nil {
                q = append(q, cur.left)
            }
            
            if cur.right != nil {
                q = append(q, cur.right)
            }
        }
    }
}

```
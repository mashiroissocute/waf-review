## 值的做

- https://leetcode.cn/problems/permutations/description/
- https://leetcode.cn/problems/n-queens/

## 框架
```golang
result = []
def backtrack(路径, 选择列表):
    if 满足结束条件:
        result.add(路径)
        return
    
    for 选择 in 选择列表:
        做选择 //if !isValid{ continue }
        backtrack(路径, 选择列表)
        撤销选择
```

- 与二叉树的遍历类似，但是因为二叉树不存在循环遍历的问题，所以不需要做选择和撤销选择。
在回溯问题中，一般存在循环遍历的问题，所以要做出选择和撤销选择。

- 与BFS对比，BFS也有做选择的步骤（visited），但是每次做出选择都是一种结果，所以不需要再撤销选择。
在回溯问题中，做出选择只是进入下一步，如果下一步没有拿到想要的结果，那就需要撤销选择。

## 全排列

```golang
func all(nums []int) [][]int{
    res := make([][]int,0)

    track := make([]int,0)
    used := make([]bool,len(nums))

    var backsource func() 
    backsource = func() {

        //如果满足条件
        if len(track) == len(nums) {
            res = append(res, append([]int{},track)) // 非常重要的一点！ track会变动，必须使用append([]int{},track)，这个append的实现有关。
            return
        }

        for i,val := range nums{

            // 剪枝
            if used[i] {
                continue
            }

            // 做选择
            track = append(track, val)
            used[i] = true
            
            backsource()

            //撤销选择
            track = track[:len(track)-1]
            used[i] = false
        }
        
    }

    backsource()
    return res
}
```


## N皇后
```golang

func isVaild(board []string, i int, j int) bool {
     n := len(board)
    // 检查列是否有皇后冲突
    for i := 0; i < n; i++ {
        if board[i][col] == 'Q' {
            return false
        }
    }
    // 检查右上方是否有皇后冲突
    for i, j := row-1, col+1; i >= 0 && j < n; i, j = i-1, j+1 {
        if board[i][j] == 'Q' {
            return false
        }
    }
    // 检查左上方是否有皇后冲突
    for i, j := row-1, col-1; i >= 0 && j >= 0; i, j = i-1, j-1 {
        if board[i][j] == 'Q' {
            return false
        }
    }
    return true
}


func Nqueue(n int) [][]string {
    res := make([][]string,0)

    board := make([]string,n)
    for i := range board {
        board[i] = strings.repeat(".", n)
    }

    var backsource func(i int) 
    backsource = func(i int) {

        if i == len(board) {
            res = append(res, append([]string{},board))
            return 
        }

        for j:=0;j<len(board);j++{
            if !isVaild(board,i,j) {
                continue
            }

            borad[i][j] = 'Q' //省略了[]byte和string之间的转换
            backsource(i+1)
            borad[i][j] = '.'
        }

    }
    
    backsource(0)
    return res
}
```

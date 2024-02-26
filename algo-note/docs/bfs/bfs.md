## 值的做
- https://leetcode.cn/problems/open-the-lock/



## 解决的问题
BFS一般用于解决最短路径问题。
问题的本质就是让你在一幅「图」中找到从起点 start 到终点 target 的最近距离

## BFS和DFS
从[二叉树最小深度](#二叉树最小深度)一题中可以得出体会：

- 1、为什么 BFS 可以找到最短距离，DFS 不行吗？

首先，BFS 的逻辑，depth 每增加一次，队列中的所有节点都向前迈一步，这保证了第一次到达终点的时候，走的步数是最少的。

DFS 不能找最短路径吗？其实也是可以的，但是时间复杂度相对高很多。DFS 实际上是靠递归的堆栈记录走过的路径，你要找到最短路径，肯定得把二叉树中所有树杈都探索完才能对比出最短的路径有多长。而 BFS 借助队列做到一次一步「齐头并进」，是可以在不遍历完整棵树的条件下找到最短距离的。

形象点说，DFS 是线，BFS 是面；

- 2、既然 BFS 那么好，为啥 DFS 还要存在？

BFS 可以找到最短距离，但是空间复杂度高，而 DFS 的空间复杂度较低。

假设给你的这个二叉树是满二叉树，节点数为 N，对于 DFS 算法来说，空间复杂度无非就是递归堆栈，最坏情况下顶多就是树的高度，也就是 O(logN)。

但是 BFS 算法，队列中每次都会储存着二叉树一层的节点，这样的话最坏情况下空间复杂度应该是树的最底层节点的数量，也就是 N/2，用 Big O 表示的话也就是 O(N)。

由此观之，BFS 还是有代价的，一般来说在找最短路径的时候使用 BFS，其他时候还是 DFS 使用得多一些（主要是递归代码好写）。


## 框架
```golang
// 计算从起点 start 到终点 target 的最近距离
func BFS(start Node, target Node) int {
    q := make([]Node, 0) // 核心数据结构
    visited := make(map[Node]bool) // 避免走回头路
    
    q = append(q, start) // 将起点加入队列
    visited[start] = true

    for len(q) > 0 {
        sz := len(q)
        /* 将当前队列中的所有节点向四周扩散 */
        for i := 0; i < sz; i++ {
            cur := q[0]
            q = q[1:]
            /* 划重点：这里判断是否到达终点 */
            if cur == target {
                return step
            }
            /* 将 cur 的相邻节点加入队列 */
            for _, x := range cur.adj() {
                if _, ok := visited[x]; !ok {
                    q = append(q, x)
                    visited[x] = true
                }
            }
        }
    }
    // 如果走到这里，说明在图中没有找到目标节点
}
```


## 二叉树最小深度
```golang
//dfs
func minDeep(root *TreeNode) int {
    if root == nil {
        return 0
    }

    minLeft := minDeep(root.left)
    minRight := minDeep(root.right)
    return min(minLeft,minRight) + 1
}
```

```golang
//bfs
func minDeep(root *TreeNode) int{
    q := make([]*TreeNode,0)
    deep := 0
    if root != nil {
        q = append(q, root)
    }

    for len(q) != 0 {
        sz := len(q)
        for i:=0;i<sz;i++{
            cur := q[0]
            q = q[1:]

            if cur.left == nil && cur.right==nil{
                return deep
            }

            if cur.left != nil {
                q = append(q, cur.left)
            }
            if cur.right != nil{
                q = append(q, cur.right)
            }
        }
        deep++
    }
    return deep
}
```



### 解开密码锁的最少次数
```golang
func up(s string, pos int) string{
    cur := []byte(s)
    if cur[i] == '9' {
        cur[i] = '0'
    }else{
        cur[i] = cur[i] + 1
    }
    return string(cur)
}

func down(s string, pos int) string{
    cur := []byte(s)
    if cur[i] == '0' {
        cur[i] = '9'
    }else{
        cur[i] = cur[i] - 1
    }
    return string(cur)
}

func ...(deadlocks []string, target string) int {
    
    visited := make(map[string]bool)
    deadlockMap := make(map[string]bool)
    for _, deadlock := range deadlocks{
        deadlockMap[deadlock] = true
    }

    q := make([]string,0)
    times := 0


    if target != "0000" {
        q = append(q, "0000")
    }

    for len(q) != 0 {
        sz := len(q)
        for i:=0;i<sz;i++{
            cur := q[0]
            q = q[1:]

            if cur == target {
                return times                
            }

            if _, ok := deadlockMap[cur]; ok{
                continue
            }

            //将当前状态的所有下一个状态加入到q
            for j:=0;j<4;j++{
                upper := up(cur,j)
                if _,ok := visited[upper]; !ok{
                    q = append(q,upper)
                    visited[q] = true
                }

                if _,ok := visited[downer]; !ok{
                    q = append(q,downer)
                    visited[q] = true
                }
            }
        }
        times++
    }
    return -1
}
```
package main

import "container/heap"

func removeElement(nums []int, val int) int {
	slow, fast := 0, 0

	for fast < len(nums) {
		if nums[fast] != val {
			nums[slow] = nums[fast]
			slow++
		}

		fast++
	}
	return slow
}

func removeDuplicates(nums []int) int {
	slow, fast := 0, 0

	for fast < len(nums) {
		if nums[slow] != nums[fast] {
			slow++
			nums[slow] = nums[fast]
		}
		fast++
	}

	return slow + 1
}

func deleteDuplicates(head *ListNode) *ListNode {
	if head == nil {
		return nil
	}
	slow, fast := head, head

	for fast != nil {
		if slow.Val != fast.Val {
			slow = slow.Next
			slow.Val = fast.Val
		}

		fast = fast.Next
	}

	slow.Next = nil
	return head
}

type ListNode struct {
	Val  int
	Next *ListNode
}

func mergeTwoLists(list1 *ListNode, list2 *ListNode) *ListNode {
	l := &ListNode{-1, nil}
	p := l

	l1 := list1
	l2 := list2

	for l1 != nil && l2 != nil {
		if l1.Val < l2.Val {
			p.Next = l1
			p = p.Next

			temp := l1.Next
			l1.Next = nil
			l1 = temp
		} else {
			p.Next = l2
			p = p.Next

			temp := l2.Next
			l2.Next = nil
			l2 = temp
		}
	}

	if l1 != nil {
		p.Next = l1
	}

	if l2 != nil {
		p.Next = l2
	}

	return l.Next
}

func partition(head *ListNode, x int) *ListNode {
	l1 := &ListNode{-1, nil}
	l2 := &ListNode{-1, nil}
	p := head
	p1 := l1
	p2 := l2

	for p != nil {
		if p.Val < x {
			p1.Next = p
			p1 = p1.Next
		} else {
			p2.Next = p
			p2 = p2.Next
		}

		temp := p.Next
		p.Next = nil
		p = temp
	}

	p1.Next = l2

	return l1.Next
}

func mergeKLists(lists []*ListNode) *ListNode {

	mergeList := &ListNode{-1, nil}
	p := mergeList

	mH := make(minHeap, 0)

	heap.Init(&mH)

	for _, l := range lists {
		if l != nil {
			heap.Push(&mH, l)
		}
	}

	for {
		if len(mH) == 0 {
			break
		}

		l := heap.Pop(&mH).(*ListNode)

		p.Next = l
		p = p.Next

		temp := l.Next
		l.Next = nil
		l = temp

		if l != nil {
			heap.Push(&mH, l)
		}

	}

	return mergeList.Next
}

type minHeap []*ListNode

func (mH *minHeap) Len() int {
	return len(*mH)
}

func (mH *minHeap) Less(i, j int) bool {
	return (*mH)[i].Val < (*mH)[j].Val
}

func (mH *minHeap) Swap(i, j int) {
	(*mH)[i], (*mH)[j] = (*mH)[j], (*mH)[i]
}

func (mH *minHeap) Push(x interface{}) {
	*mH = append(*mH, x.(*ListNode))
}

func (mH *minHeap) Pop() interface{} {
	old := *mH
	n := len(old)
	x := old[n-1]
	*mH = old[0 : n-1]
	return x
}

// 删除节点，最好使用一下虚拟节点
func removeNthFromEnd(head *ListNode, n int) *ListNode {
	// 找到倒数第n+1个点，可能是虚拟节点，该节点为slow
	slow, fast := &ListNode{-1, head}, &ListNode{-1, head}

	for i := 0; i < n+1; i++ {
		fast = fast.Next
	}

	for fast != nil {
		slow = slow.Next
		fast = fast.Next
	}

	// 如果要删除的第n个节点是head节点
	if slow.Next == head {
		return head.Next
	}

	// 删除倒数第n+1节点后面的倒数第n节点
	slow.Next = slow.Next.Next
	return head
}

func middleNode(head *ListNode) *ListNode {
	slow, fast := head, head
	for fast != nil && fast.Next != nil {
		fast = fast.Next.Next
		slow = slow.Next
	}
	return slow
}

func detectCycle(head *ListNode) *ListNode {
	slow, fast := head, head
	for fast != nil && fast.Next != nil {
		fast = fast.Next.Next
		slow = slow.Next

		if slow == fast {
			p := head
			for p != slow {
				p = p.Next
				slow = slow.Next
			}
			return p
		}
	}
	return nil
}

func getIntersectionNode(headA, headB *ListNode) *ListNode {

	p1, p2 := headA, headB
	if p1 == nil || p2 == nil {
		return nil
	}

	swithA, swithB := false, false

	for p1 != nil && p2 != nil {

		if p1 == p2 {
			return p1
		}

		p1 = p1.Next
		p2 = p2.Next

		if p1 == nil && !swithA {
			p1 = headB
			swithA = true
		}

		if p2 == nil && !swithB {
			p2 = headA
			swithB = true
		}

	}

	return nil

}










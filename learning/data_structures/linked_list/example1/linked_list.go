package example1

import "sort"

type Node struct {
	Value int
	Next  *Node
}

var root = new(Node)

func addNode(t *Node, v int) int {
	if root == nil {
		t = &Node{v, nil}
		root = t
		return 0
	}
	if v == t.Value {
		return -1
	}
	if t.Next == nil {
		t.Next = &Node{v, nil}
		return -2
	}
	return addNode(t.Next, v)
}

func sortList(node *Node) (sortedList *Node) {
	values := traverse(node)
	sort.Ints(values)
	var prevNode *Node
	for _, v := range values {
		currNode := &Node{
			Value: v,
		}
		if prevNode != nil {
			prevNode.Next = currNode
		}
		prevNode = currNode
		if sortedList == nil {
			sortedList = currNode
		}
	}
	return
}

func traverse(t *Node) (result []int) {
	if t == nil {
		return
	}
	for t != nil {
		result = append(result, t.Value)
		t = t.Next
	}
	return
}

func lookupNode(t *Node, v int) bool {
	if root == nil {
		t = &Node{v, nil}
		root = t
		return false
	}
	if v == t.Value {
		return true
	}
	if t.Next == nil {
		return false
	}
	return lookupNode(t.Next, v)
}

func size(t *Node) int {
	if t == nil {
		return 0
	}
	i := 0
	for t != nil {
		i++
		t = t.Next
	}
	return i
}

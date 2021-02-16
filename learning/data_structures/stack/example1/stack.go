package exmaple1

type NodeVal struct {
	Value,
	Number,
	Seed int
}

type Node struct {
	Value NodeVal
	Next  *Node
}

var size = 0
var stack = new(Node)

func Push(v NodeVal) bool {
	if stack == nil {
		stack = &Node{v, nil}
		size = 1
		return true
	}

	temp := &Node{v, nil}
	temp.Next = stack
	stack = temp
	size++
	return true
}

func Pop(t *Node) (NodeVal, bool) {
	if size == 0 {
		return NodeVal{}, false
	}

	if size == 1 {
		size = 0
		stack = nil
		return t.Value, true
	}

	stack = stack.Next
	size--
	return t.Value, true
}

func traverse(t *Node) (result []int) {
	if size == 0 {
		return
	}
	for t != nil {
		result = append(result, t.Value.Value)
		t = t.Next
	}
	return
}

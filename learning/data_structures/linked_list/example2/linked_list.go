package example2

type Node struct {
	Value    int
	Previous *Node
	Next     *Node
}

var root = new(Node)

func addNode(t *Node, v int) int {
	if root == nil {
		t = &Node{Value: v, Previous: nil, Next: nil}
		root = t
		return 0
	}
	if v == t.Value {
		return -1
	}
	if t.Next == nil {
		temp := t
		t.Next = &Node{Value: v, Previous: temp, Next: nil}
		return -2
	}
	return addNode(t.Next, v)
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

func reverse(t *Node) (result []int) {
	if t == nil {
		return
	}

	temp := t
	for t != nil {
		temp = t
		t = t.Next
	}

	for temp.Previous != nil {
		result = append(result, temp.Value)
		temp = temp.Previous
	}

	return
}

func lookupNode(t *Node, v int) bool {
	if root == nil {
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

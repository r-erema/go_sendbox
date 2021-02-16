package exmaple1

type Node struct {
	Value float64
	Next  *Node
}

var size = 0
var queue = new(Node)

func Push(t *Node, v float64) bool {
	if queue == nil {
		queue = &Node{v, nil}
		size++
		return true
	}

	t = &Node{Value: v, Next: nil}
	t.Next = queue
	queue = t
	size++

	return true
}

func Pop(t *Node) (float64, bool) {
	if size == 0 {
		return 0, false
	}

	if size == 1 {
		queue = nil
		size--
		return t.Value, true
	}

	temp := t
	for t.Next != nil {
		temp = t
		t = t.Next
	}

	v := (temp.Next).Value
	temp.Next = nil

	size--
	return v, true
}

func traverse(t *Node) (result []float64) {
	if size == 0 {
		return
	}
	for t != nil {
		result = append(result, t.Value)
		t = t.Next
	}
	return
}

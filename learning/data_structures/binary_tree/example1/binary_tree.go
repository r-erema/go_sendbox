package example1

import (
	"math/rand"
	"time"
)

type Tree struct {
	Left, Right *Tree
	Value       int
}

func traverse(t *Tree) (result []int) {
	if t == nil {
		return result
	}

	result = append(result, traverse(t.Left)...)
	result = append(result, t.Value)
	result = append(result, traverse(t.Right)...)
	return
}

func create(n int) (t *Tree) {
	rand.Seed(time.Now().Unix())
	for i := 0; i < 2*n; i++ {
		t = insert(t, rand.Intn(n*2))
	}
	return t
}

func insert(t *Tree, value int) *Tree {
	if t == nil {
		return &Tree{Left: nil, Right: nil, Value: value}
	}

	if value == t.Value {
		return t
	}

	if value < t.Value {
		t.Left = insert(t.Left, value)
		return t
	}

	t.Right = insert(t.Right, value)
	return t
}

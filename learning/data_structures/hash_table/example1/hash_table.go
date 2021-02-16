package example1

const (
	ModeChainingMethod = iota + 1
	ModeLinearProbingMethod
	Size = 15
)

type Node struct {
	Value int
	Next  *Node
}

type HashTable struct {
	Chain map[int]*Node
	Array []*int
	Size  int
	Mode  int
}

func hashFunction(i, size int) int {
	return i % size
}

func insert(table *HashTable, value int) int {
	index := hashFunction(value, table.Size)

	switch table.Mode {
	case ModeChainingMethod:
		element := &Node{Value: value, Next: table.Chain[index]}
		table.Chain[index] = element
	case ModeLinearProbingMethod:
		for table.Array[index] != nil && len(table.Array) > index+1 {
			index++
		}
		table.Array[index] = &value
	default:
		unknownMode()
	}

	return index
}

func traverse(table *HashTable) (result []int) {
	switch table.Mode {
	case ModeChainingMethod:
		for k := range table.Chain {
			t := table.Chain[k]
			for t != nil {
				result = append(result, t.Value)
				t = t.Next
			}
		}
	case ModeLinearProbingMethod:
		for _, v := range table.Array {
			if v != nil {
				result = append(result, *v)
			}
		}
		return result
	default:
		unknownMode()
	}
	return
}

func lookup(table *HashTable, value int) bool {
	index := hashFunction(value, table.Size)
	switch table.Mode {
	case ModeChainingMethod:
		t := table.Chain[index]
		for t != nil {
			if t.Value == value {
				return true
			}
			t = t.Next
		}
	case ModeLinearProbingMethod:

	default:
		unknownMode()
	}
	return false
}

func unknownMode() {
	panic("unknown method")
}

package example1

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_traverseChainingMethod(t *testing.T) {
	table := &HashTable{Chain: map[int]*Node{}, Array: []*int{}, Size: Size, Mode: ModeChainingMethod}
	for i := 0; i < 120; i++ {
		insert(table, i)
	}

	actual := traverse(table)
	expected := []int{
		107, 92, 77, 62, 47, 32, 17, 2, 111, 96, 81, 66, 51, 36, 21, 6, 113, 98, 83, 68, 53, 38, 23, 8, 119, 104, 89,
		74, 59, 44, 29, 14, 106, 91, 76, 61, 46, 31, 16, 1, 109, 94, 79, 64, 49, 34, 19, 4, 110, 95, 80, 65, 50, 35, 20,
		5, 117, 102, 87, 72, 57, 42, 27, 12, 116, 101, 86, 71, 56, 41, 26, 11, 118, 103, 88, 73, 58, 43, 28, 13, 105,
		90, 75, 60, 45, 30, 15, 0, 108, 93, 78, 63, 48, 33, 18, 3, 112, 97, 82, 67, 52, 37, 22, 7, 114, 99, 84, 69, 54,
		39, 24, 9, 115, 100, 85, 70, 55, 40, 25, 10,
	}
	assert.ElementsMatch(t, expected, actual)
}

func Test_traverseProbingMethod(t *testing.T) {
	table := &HashTable{Chain: map[int]*Node{}, Array: make([]*int, 100), Size: 10, Mode: ModeLinearProbingMethod}
	insert(table, 5)
	insert(table, 245)
	insert(table, 13)
	insert(table, 12)
	insert(table, 14)
	insert(table, 7)
	insert(table, 378)
	insert(table, 101)
	insert(table, 8)
	insert(table, 9)

	actual := traverse(table)
	expected := []int{101, 12, 13, 14, 5, 245, 7, 378, 8, 9}
	assert.Equal(t, expected, actual)
}

func Test_lookup(t *testing.T) {
	tests := []struct {
		number   int
		expected bool
		mode     int
	}{
		{number: 1, expected: true, mode: ModeChainingMethod},
		{number: 11, expected: true, mode: ModeChainingMethod},
		{number: 119, expected: true, mode: ModeChainingMethod},
		{number: 120, expected: false, mode: ModeChainingMethod},
		{number: 1000, expected: false, mode: ModeChainingMethod},
	}

	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("Test %d", i), func(t *testing.T) {

			table := &HashTable{Chain: map[int]*Node{}, Array: []*int{}, Size: Size, Mode: tt.mode}
			for i := 0; i < 120; i++ {
				insert(table, i)
			}

			actual := lookup(table, tt.number)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

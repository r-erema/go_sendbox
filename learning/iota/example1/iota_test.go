package example1

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_getDayNumber(t *testing.T) {
	tests := []struct {
		day         string
		expectedDay int
	}{
		{day: Monday, expectedDay: 1},
		{day: Tuesday, expectedDay: 2},
		{day: Wednesday, expectedDay: 3},
		{day: Thursday, expectedDay: 4},
		{day: Friday, expectedDay: 5},
		{day: Saturday, expectedDay: 6},
		{day: Sunday, expectedDay: 7},
		{day: "wrong_day", expectedDay: -1},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.day, func(t *testing.T) {
			number, err := dayNumber(tt.day)
			assert.Equal(t, tt.expectedDay, number)
			if number != -1 {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

package ip

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_wrongIPPart(t *testing.T) {
	tests := []struct {
		ip             string
		expectedResult []string
	}{
		{
			ip:             "123.123.222",
			expectedResult: []string{"123.123.222"},
		},
		{
			ip:             "192.168.a.7",
			expectedResult: []string{"a"},
		},
		{
			ip:             "500.168.256.7",
			expectedResult: []string{"500", "256"},
		},
		{
			ip:             "192.168.0.7",
			expectedResult: nil,
		},
	}
	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%d", i+1), func(t *testing.T) {
			wrongParts := wrongIPPart(tt.ip)
			assert.Equal(t, tt.expectedResult, wrongParts)
		})
	}
}

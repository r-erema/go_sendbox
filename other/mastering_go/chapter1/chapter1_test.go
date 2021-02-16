package chapter1

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func Test_sum(t *testing.T) {
	tests := []struct {
		userInput []string
		result    float64
	}{
		{userInput: []string{"8.7", "0", "14", "-2.2"}, result: 20.5},
		{userInput: []string{"4.11", "not_number", "9.32"}, result: 13.43},
	}
	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%d", i+1), func(t *testing.T) {
			os.Args = tt.userInput
			result := sum()
			assert.Equal(t, tt.result, result)
		})
	}
}

func Test_average(t *testing.T) {
	tests := []struct {
		userInput []string
		result    float64
	}{
		{userInput: []string{"35.7", "0.01", "abc", "-3"}, result: 10.903333333333334},
		{userInput: []string{"5.576", "75", "0.02"}, result: 26.86533333333333},
	}
	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%d", i+1), func(t *testing.T) {
			os.Args = tt.userInput
			result := average()
			assert.Equal(t, tt.result, result)
		})
	}
}

func Test_sumInt(t *testing.T) {
	tests := []struct {
		userInput []string
		result    int64
	}{
		{userInput: []string{"1", "6", "-74", StopWord, "3"}, result: -67},
		//{userInput: []string{StopWord, "3", "2"}, result: 0},
	}
	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%d", i+1), func(t *testing.T) {
			os.Args = tt.userInput
			result := sumInt()
			assert.Equal(t, tt.result, result)
		})
	}
}

func Test_customLog(t *testing.T) {
	fileName1, fileName2 := "test.log", "test2.log"
	f1, err := os.OpenFile(fileName1, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	require.NoError(t, err)
	f2, err := os.OpenFile(fileName2, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	require.NoError(t, err)
	customLog("Log data", f1, f2)
	f1Stat, _ := f1.Stat()
	f2Stat, _ := f2.Stat()
	assert.Equal(t, int64(49), f1Stat.Size())
	assert.Equal(t, int64(49), f2Stat.Size())
	err = os.Remove(fileName1)
	require.NoError(t, err)
	err = os.Remove(fileName2)
	require.NoError(t, err)
}

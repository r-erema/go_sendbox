package exmaple8

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/shirou/gopsutil/cpu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGoMaxProcs(t *testing.T) {

	tests := []struct {
		GoMaxProcsN                      int
		ExpectedCurrentActiveLogicalCPUs int
		ExpectedAllLogicalCPUs           int
		ExpectedAllPhysicalCPUs          int32
	}{
		{
			GoMaxProcsN:                      0,
			ExpectedCurrentActiveLogicalCPUs: 8,
			ExpectedAllLogicalCPUs:           8,
			ExpectedAllPhysicalCPUs:          4,
		},
		{
			GoMaxProcsN:                      4,
			ExpectedCurrentActiveLogicalCPUs: 4,
			ExpectedAllLogicalCPUs:           8,
			ExpectedAllPhysicalCPUs:          4,
		},
		{
			GoMaxProcsN:                      2,
			ExpectedCurrentActiveLogicalCPUs: 2,
			ExpectedAllLogicalCPUs:           8,
			ExpectedAllPhysicalCPUs:          4,
		},
		{
			GoMaxProcsN:                      5,
			ExpectedCurrentActiveLogicalCPUs: 5,
			ExpectedAllLogicalCPUs:           8,
			ExpectedAllPhysicalCPUs:          4,
		},
		{
			GoMaxProcsN:                      100,
			ExpectedCurrentActiveLogicalCPUs: 100,
			ExpectedAllLogicalCPUs:           8,
			ExpectedAllPhysicalCPUs:          4,
		},
	}

	for i, test := range tests {
		test := test
		t.Run(fmt.Sprintf("%d", i+1), func(t *testing.T) {
			_ = runtime.GOMAXPROCS(test.GoMaxProcsN)
			currentLogicalCPUs := runtime.GOMAXPROCS(0)
			assert.Equal(t, test.ExpectedCurrentActiveLogicalCPUs, currentLogicalCPUs)

			currentAllLogicalCPUs := runtime.NumCPU()
			assert.Equal(t, test.ExpectedAllLogicalCPUs, currentAllLogicalCPUs)

			info, err := cpu.Info()
			require.NoError(t, err)
			assert.Equal(t, test.ExpectedAllPhysicalCPUs, info[0].Cores)
		})
	}

}

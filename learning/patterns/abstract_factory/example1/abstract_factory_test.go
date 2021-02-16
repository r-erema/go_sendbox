package example1

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAbstractFactory(t *testing.T) {

	cocaColaFactory := CocaColaFactory{}
	liquid := cocaColaFactory.CreateLiquid(2.5)
	bottle := cocaColaFactory.CreateBottle(2.5)

	bottle.PourLiquid(liquid)

	assert.Equal(t, bottle.GetLiquidVolume(), bottle.GetVolume())
}

package example1

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEmbedding(t *testing.T) {
	c := Common{
		first:  first{},
		second: second{},
	}

	c.share()
	assert.Equal(t, "shared for 2nd", c.share())

}

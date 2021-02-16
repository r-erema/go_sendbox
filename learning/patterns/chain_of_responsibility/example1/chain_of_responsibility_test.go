package example1

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestChainOfResponsibility(t *testing.T) {

	s := Server{}

	r, _ := http.NewRequest("GET", "/?id=7", nil)
	response := s.handle(*r)
	assert.Equal(t, "400", response.Status)

	r, _ = http.NewRequest("GET", "/?id=8", nil)
	response = s.handle(*r)
	assert.Equal(t, "403", response.Status)

	r, _ = http.NewRequest("GET", "/?id=9", nil)
	response = s.handle(*r)
	assert.Equal(t, "200", response.Status)

}

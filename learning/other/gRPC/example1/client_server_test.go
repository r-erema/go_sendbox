package example1

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
	"time"
)

var cases = []struct {
	requestName, expected string
}{
	{"World", "Hello World"},
	{"John Doe", "Hello John Doe"},
	{"Roma", "Hello Roma"},
}

func TestClientServer(t *testing.T) {

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Errorf("failed to listen: %v", err)
	}

	s := CreateAndRunServer(listener)
	defer s.GracefulStop()

	c := CreateClient(listener.Addr().String())
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	for i, tt := range cases {
		t.Run(fmt.Sprintf("Test %d", i), func(t *testing.T) {
			r, err := c.SayHello(ctx, &HelloRequest{Name: tt.requestName})
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, tt.expected, r.GetMessage())
		})
	}

}

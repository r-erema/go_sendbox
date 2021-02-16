package example3

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type httpData struct {
	resp *http.Response
	err  error
}

func TestContextTimeout(t *testing.T) {

	tests := []struct {
		name               string
		clientTimeout      time.Duration
		serverResponseTime time.Duration
		expectedError      error
	}{
		{
			name:               "ok",
			clientTimeout:      time.Millisecond * 10,
			serverResponseTime: time.Millisecond * 5,
			expectedError:      nil,
		},
		{
			name:               "timeout",
			clientTimeout:      time.Millisecond * 5,
			serverResponseTime: time.Millisecond * 10,
			expectedError:      context.DeadlineExceeded,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				time.Sleep(tt.serverResponseTime)
			}))

			ctx, cancel := context.WithTimeout(context.Background(), tt.clientTimeout)
			defer cancel()
			err := connect(ctx, server.URL)
			assert.Equal(t, tt.expectedError, err)
		})
	}

}

func connect(ctx context.Context, url string) error {
	data := make(chan httpData, 1)
	t := &http.Transport{}
	httpClient := http.Client{Transport: t}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	go func() {
		resp, err := httpClient.Do(req)
		if err != nil {
			log.Println(err)
			data <- httpData{nil, err}
		} else {
			data <- httpData{resp, err}
		}
	}()

	select {
	case <-ctx.Done():
		<-data
		log.Println("The request was cancelled")
		return ctx.Err()
	case ok := <-data:
		err = ok.err
		resp := ok.resp
		if err != nil {
			log.Println("Error select:", err)
			return err
		}
		defer func() {
			err = resp.Body.Close()
			if err != nil {
				log.Println("Error body close:", err)
			}
		}()

		realHttpData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error select:", err)
			return err
		}
		fmt.Println("Server response:", realHttpData)
	}
	return nil
}

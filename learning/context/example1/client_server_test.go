package example1

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"testing"
	"time"
)

type requestIDKey int

type Key int

func TestClientServer(t *testing.T) {

	decorate := func(h http.HandlerFunc) http.HandlerFunc {
		return func(writer http.ResponseWriter, request *http.Request) {
			ctx := request.Context()
			id := rand.Int63()
			ctx = context.WithValue(ctx, requestIDKey(11), id)
			h(writer, request.WithContext(ctx))
		}
	}

	http.HandleFunc("/", decorate(func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		log.Println(ctx, "handler started")
		id, _ := ctx.Value(requestIDKey(11)).(int64)
		log.Printf("request id: %d", id)
		defer log.Println(ctx, "handler finished")

		log.Println("value for foo is:", ctx.Value("foo"))

		select {
		case <-time.After(5 * time.Second):
			_, _ = fmt.Fprintln(writer, "hello")
		case <-request.Context().Done():
			err := ctx.Err()
			log.Println(ctx, err.Error())
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
	}))
	url := "localhost:8139"

	go func() {
		t.Fatal(http.ListenAndServe(url, nil))
	}()

	time.Sleep(time.Millisecond * 100)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	ctx = context.WithValue(ctx, Key(7), "bar")
	defer cancel()
	req, err := http.NewRequest(http.MethodGet, "http://"+url, nil)
	if err != nil {
		t.Fatal(err)
	}
	req = req.WithContext(ctx)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusOK {
		t.Fatal(res.Status)
	}
	_, _ = io.Copy(os.Stdout, res.Body)
}

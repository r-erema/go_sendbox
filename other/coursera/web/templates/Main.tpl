package web

import (
    "encoding/json"
    "fmt"
    "github.com/gorilla/schema"
    "github.com/thoas/go-funk"
    "io/ioutil"
    "net/http"
    "net/url"
    "strconv"
    "strings"
)

var decoder  = schema.NewDecoder()

type Response map[string]interface{}

func jsonResponse(error string, response interface{}) ([]byte, error) {
    if response != nil {
        return json.Marshal(&Response{"error": error, "response": response})
    }
    return json.Marshal(&Response{"error": error})
}

{{ .ApisHTTPServe }}

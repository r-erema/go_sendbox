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

var decoder = schema.NewDecoder()

type Response map[string]interface{}

func jsonResponse(error string, response interface{}) ([]byte, error) {
	if response != nil {
		return json.Marshal(&Response{"error": error, "response": response})
	}
	return json.Marshal(&Response{"error": error})
}

func (srv *MyApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path == "/user/profile" {

		srv.ProfileHandler(w, r)
		return
	}

	if r.URL.Path == "/user/create" {

		if r.Method != "POST" {
			resp, _ := jsonResponse("bad method", nil)
			http.Error(w, string(resp), http.StatusNotAcceptable)
			return
		}

		if r.Header.Get("X-Auth") != "100500" {
			resp, _ := jsonResponse("unauthorized", nil)
			http.Error(w, string(resp), http.StatusForbidden)
			return
		}

		srv.CreateHandler(w, r)
		return
	}

	resp, _ := jsonResponse("unknown method", nil)
	http.Error(w, string(resp), http.StatusNotFound)
}

func (srv *MyApi) ProfileHandler(w http.ResponseWriter, r *http.Request) {
	in := ProfileParams{}
	var query url.Values
	var err error
	if r.Method == "GET" {
		query = r.URL.Query()
		_ = decoder.Decode(&in, r.URL.Query())
	} else {
		data, _ := ioutil.ReadAll(r.Body)
		query, err = url.ParseQuery(string(data))
		if err != nil {
			http.Error(w, "Can't parse body", http.StatusBadRequest)
			return
		}

	}
	_ = decoder.Decode(&in, query)

	if in.Login == "" {
		resp, _ := jsonResponse("login must me not empty", nil)
		http.Error(w, string(resp), http.StatusBadRequest)
		return
	}

	user, err := srv.Profile(r.Context(), in)
	if err != nil {
		if apiErr, ok := err.(ApiError); ok {
			resp, _ := jsonResponse(apiErr.Error(), nil)
			http.Error(w, string(resp), apiErr.HTTPStatus)
		} else {
			resp, _ := jsonResponse(err.Error(), nil)
			http.Error(w, string(resp), http.StatusInternalServerError)
		}
		return
	}
	jsonBytes, _ := json.Marshal(&Response{"error": "", "response": user})
	_, err = w.Write(jsonBytes)
	if err != nil {
		http.Error(w, "Can't encode json", http.StatusInternalServerError)
	}
}

func (srv *MyApi) CreateHandler(w http.ResponseWriter, r *http.Request) {
	in := CreateParams{}
	var query url.Values
	var err error
	if r.Method == "GET" {
		query = r.URL.Query()
		_ = decoder.Decode(&in, r.URL.Query())
	} else {
		data, _ := ioutil.ReadAll(r.Body)
		query, err = url.ParseQuery(string(data))
		if err != nil {
			http.Error(w, "Can't parse body", http.StatusBadRequest)
			return
		}

	}
	_ = decoder.Decode(&in, query)

	if in.Age < 0 {
		resp, _ := jsonResponse("age must be >= 0", nil)
		http.Error(w, string(resp), http.StatusBadRequest)
		return
	}

	if in.Age > 128 {
		resp, _ := jsonResponse("age must be <= 128", nil)
		http.Error(w, string(resp), http.StatusBadRequest)
		return
	}

	if _, err := strconv.Atoi(query.Get("age")); err != nil {
		resp, _ := jsonResponse("age must be int", nil)
		http.Error(w, string(resp), http.StatusBadRequest)
		return
	}

	if in.Login == "" {
		resp, _ := jsonResponse("login must me not empty", nil)
		http.Error(w, string(resp), http.StatusBadRequest)
		return
	}

	if len(in.Login) < 10 {
		resp, _ := jsonResponse("login len must be >= 10", nil)
		http.Error(w, string(resp), http.StatusBadRequest)
		return
	}

	in.Name = query.Get("full_name")

	if funk.IndexOf(strings.Split("user|moderator|admin", "|"), in.Status) == -1 && in.Status != "" {
		resp, _ := jsonResponse(fmt.Sprintf("status must be one of [%s]", strings.ReplaceAll("user|moderator|admin", "|", ", ")), nil)
		http.Error(w, string(resp), http.StatusBadRequest)
		return
	}

	if in.Status == "" {
		in.Status = "user"
	}

	user, err := srv.Create(r.Context(), in)
	if err != nil {
		if apiErr, ok := err.(ApiError); ok {
			resp, _ := jsonResponse(apiErr.Error(), nil)
			http.Error(w, string(resp), apiErr.HTTPStatus)
		} else {
			resp, _ := jsonResponse(err.Error(), nil)
			http.Error(w, string(resp), http.StatusInternalServerError)
		}
		return
	}
	jsonBytes, _ := json.Marshal(&Response{"error": "", "response": user})
	_, err = w.Write(jsonBytes)
	if err != nil {
		http.Error(w, "Can't encode json", http.StatusInternalServerError)
	}
}

func (srv *OtherApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path == "/user/create" {

		if r.Method != "POST" {
			resp, _ := jsonResponse("bad method", nil)
			http.Error(w, string(resp), http.StatusNotAcceptable)
			return
		}

		if r.Header.Get("X-Auth") != "100500" {
			resp, _ := jsonResponse("unauthorized", nil)
			http.Error(w, string(resp), http.StatusForbidden)
			return
		}

		srv.CreateHandler(w, r)
		return
	}

	resp, _ := jsonResponse("unknown method", nil)
	http.Error(w, string(resp), http.StatusNotFound)
}

func (srv *OtherApi) CreateHandler(w http.ResponseWriter, r *http.Request) {
	in := OtherCreateParams{}
	var query url.Values
	var err error
	if r.Method == "GET" {
		query = r.URL.Query()
		_ = decoder.Decode(&in, r.URL.Query())
	} else {
		data, _ := ioutil.ReadAll(r.Body)
		query, err = url.ParseQuery(string(data))
		if err != nil {
			http.Error(w, "Can't parse body", http.StatusBadRequest)
			return
		}

	}
	_ = decoder.Decode(&in, query)

	if funk.IndexOf(strings.Split("warrior|sorcerer|rouge", "|"), in.Class) == -1 && in.Class != "" {
		resp, _ := jsonResponse(fmt.Sprintf("class must be one of [%s]", strings.ReplaceAll("warrior|sorcerer|rouge", "|", ", ")), nil)
		http.Error(w, string(resp), http.StatusBadRequest)
		return
	}

	if in.Class == "" {
		in.Class = "warrior"
	}

	if in.Level < 1 {
		resp, _ := jsonResponse("level must be >= 1", nil)
		http.Error(w, string(resp), http.StatusBadRequest)
		return
	}

	if in.Level > 50 {
		resp, _ := jsonResponse("level must be <= 50", nil)
		http.Error(w, string(resp), http.StatusBadRequest)
		return
	}

	if _, err := strconv.Atoi(query.Get("level")); err != nil {
		resp, _ := jsonResponse("level must be int", nil)
		http.Error(w, string(resp), http.StatusBadRequest)
		return
	}

	in.Name = query.Get("account_name")

	if in.Username == "" {
		resp, _ := jsonResponse("username must me not empty", nil)
		http.Error(w, string(resp), http.StatusBadRequest)
		return
	}

	if len(in.Username) < 3 {
		resp, _ := jsonResponse("username len must be >= 3", nil)
		http.Error(w, string(resp), http.StatusBadRequest)
		return
	}

	user, err := srv.Create(r.Context(), in)
	if err != nil {
		if apiErr, ok := err.(ApiError); ok {
			resp, _ := jsonResponse(apiErr.Error(), nil)
			http.Error(w, string(resp), apiErr.HTTPStatus)
		} else {
			resp, _ := jsonResponse(err.Error(), nil)
			http.Error(w, string(resp), http.StatusInternalServerError)
		}
		return
	}
	jsonBytes, _ := json.Marshal(&Response{"error": "", "response": user})
	_, err = w.Write(jsonBytes)
	if err != nil {
		http.Error(w, "Can't encode json", http.StatusInternalServerError)
	}
}

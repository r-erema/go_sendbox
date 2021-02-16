package net

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/thoas/go-funk"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
)

const AccessToken = "_access_token_"
const PathSearchFile = "./dataset.xml"

func Search(query string, users []UserXML) []User {
	var result []User
	for _, u := range users {
		uReflection := reflect.ValueOf(u)
		for i := 0; i < uReflection.NumField(); i++ {
			str := uReflection.Field(i).String()
			_ = str
			if strings.Contains(uReflection.Field(i).String(), query) {
				id, _ := strconv.Atoi(u.ID)
				age, _ := strconv.Atoi(u.Age)
				result = append(result, User{
					Id:     id,
					Name:   u.FirstName + " " + u.LastName,
					Age:    age,
					About:  u.About,
					Gender: u.Gender,
				})
				break
			}
		}
	}
	return result
}

type UserXML struct {
	ID            string `xml:"id"`
	Guid          string `xml:"guid"`
	IsActive      string `xml:"isActive"`
	Balance       string `xml:"balance"`
	Picture       string `xml:"picture"`
	Age           string `xml:"age"`
	EyeColor      string `xml:"eyeColor"`
	FirstName     string `xml:"first_name"`
	LastName      string `xml:"last_name"`
	Gender        string `xml:"gender"`
	Company       string `xml:"company"`
	Email         string `xml:"email"`
	Phone         string `xml:"phone"`
	Address       string `xml:"address"`
	About         string `xml:"about"`
	Registered    string `xml:"registered"`
	FavoriteFruit string `xml:"favoriteFruit"`
}

type Data struct {
	XMLName xml.Name  `xml:"root"`
	Users   []UserXML `xml:"row"`
}

type SearchHandler struct {
	pathToFileForSearch string
}

func prepareBadRequestError(w http.ResponseWriter, text string) {
	w.WriteHeader(http.StatusBadRequest)
	errorResponse := SearchErrorResponse{text}
	jsonBytes, err := json.Marshal(errorResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(jsonBytes)
}

func (sh SearchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("AccessToken") != AccessToken {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	file, err := os.Open(sh.pathToFileForSearch)
	defer func() {
		_ = file.Close()
	}()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	query := r.URL.Query().Get("query")
	orderField := r.URL.Query().Get("order_field")
	orderBy, _ := strconv.Atoi(r.URL.Query().Get("order_by"))

	if query == "" {
		prepareBadRequestError(w, "ErrorEmptyQuery")
		return
	}

	if funk.IndexOf([]int{OrderByDesc, OrderByAsIs, OrderByAsc}, orderBy) == -1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if funk.IndexOf([]string{"", "Age"}, orderField) == -1 {
		prepareBadRequestError(w, "ErrorBadOrderField")
		return
	}

	data := Data{}
	b, _ := ioutil.ReadAll(file)
	err = xml.Unmarshal(b, &data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	users := Search(query, data.Users)
	foundUsersCount := len(users)
	if limit < foundUsersCount {
		users = users[offset:limit]
	}

	if foundUsersCount == 0 {
		_, _ = w.Write([]byte(""))
		return
	}

	usersBytes, err := json.Marshal(users)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(usersBytes)

}

var searchServerCommonMock = httptest.NewServer(SearchHandler{
	pathToFileForSearch: PathSearchFile,
})

var searchClientCommon = &SearchClient{
	AccessToken: AccessToken,
	URL:         searchServerCommonMock.URL,
}

func TestBadRequestParams(t *testing.T) {
	tests := []struct {
		req     SearchRequest
		wantErr bool
	}{
		{req: SearchRequest{Limit: -8}, wantErr: true},
		{req: SearchRequest{Limit: 26}, wantErr: true},
		{req: SearchRequest{Offset: -15}, wantErr: true},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("Test %d", i), func(t *testing.T) {
			_, err := searchClientCommon.FindUsers(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestServerTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second + time.Millisecond*100)
	}))
	client := &SearchClient{URL: server.URL}
	_, err := client.FindUsers(SearchRequest{})
	if err == nil {
		t.Errorf("Server shouldn't response in time")
		return
	}
}

func TestServerUnavailable(t *testing.T) {
	fakePort := "9999"
	client := &SearchClient{URL: searchServerCommonMock.URL + fakePort}
	_, err := client.FindUsers(SearchRequest{})
	if err == nil {
		t.Errorf("Search server must be unavalable")
		return
	}
}

func TestStatusUnauthorized(t *testing.T) {
	client := &SearchClient{URL: searchServerCommonMock.URL, AccessToken: "bad token"}
	_, err := client.FindUsers(SearchRequest{})
	if err == nil {
		t.Errorf("Search server must return Unauthorized status")
		return
	}
}

func TestInternalServerError(t *testing.T) {
	var server = httptest.NewServer(SearchHandler{
		pathToFileForSearch: "bad path",
	})
	client := &SearchClient{URL: server.URL, AccessToken: AccessToken}
	_, err := client.FindUsers(SearchRequest{})
	if err == nil {
		t.Errorf("Search server must return Internal Server Error status")
		return
	}
}

func TestBadRequest(t *testing.T) {
	tests := []struct {
		req     SearchRequest
		wantErr bool
	}{
		{req: SearchRequest{OrderBy: 10, Query: "test"}, wantErr: true},
		{req: SearchRequest{OrderField: "bad_field", Query: "test"}, wantErr: true},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("Test %d", i), func(t *testing.T) {
			_, err := searchClientCommon.FindUsers(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestResponseBadJSON(t *testing.T) {
	_, err := searchClientCommon.FindUsers(SearchRequest{Query: "no result query"})
	if err == nil {
		t.Errorf("Search server must return Internal Server Error status")
		return
	}
}

func TestResponseOK(t *testing.T) {
	tests := []struct {
		req SearchRequest
	}{
		{req: SearchRequest{Query: "green", Limit: 3}},
		{req: SearchRequest{Query: "1ec1fd0e-1151-482a-a791-c48fb324f519", Limit: 7}},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("Test %d", i), func(t *testing.T) {
			_, err := searchClientCommon.FindUsers(tt.req)
			if err != nil {
				t.Errorf("FindUsers() error = %v", err)
				return
			}
		})
	}
}

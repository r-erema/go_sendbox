package example1

import (
	"encoding/json"
	"net/http"
	"time"
)

type resourceServer struct {
	jwtTokens []JwtToken
}

type Resource struct {
	StatusOK    bool
	CurrentTime int64 `json:"current_time"`
}

func (r resourceServer) isTokenValid(token string) bool {
	for _, jwt := range r.jwtTokens {
		if jwt.AccessToken == token {
			return true
		}
	}
	return false
}

var rs = &resourceServer{}

type addTokenHandler struct {
	*resourceServer
}

func NewAddTokenHandler() *addTokenHandler {
	return &addTokenHandler{resourceServer: rs}
}

func (h addTokenHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	jwt := new(JwtToken)
	err := json.NewDecoder(request.Body).Decode(jwt)
	if err != nil {
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.jwtTokens = append(h.jwtTokens, *jwt)
}

type resourceHandler struct {
	*resourceServer
}

func NewResourceHandler() *resourceHandler {
	return &resourceHandler{resourceServer: rs}
}

func (h resourceHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	jwt := new(JwtToken)
	err := json.NewDecoder(request.Body).Decode(jwt)
	if err != nil {
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}

	if !h.isTokenValid(jwt.AccessToken) {
		http.Error(writer, "Invalid token", http.StatusForbidden)
		return
	}

	jsonBuf, err := json.Marshal(&Resource{
		StatusOK:    true,
		CurrentTime: time.Now().Unix(),
	})

	if err != nil {
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}

	_, err = writer.Write(jsonBuf)
	if err != nil {
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}
}

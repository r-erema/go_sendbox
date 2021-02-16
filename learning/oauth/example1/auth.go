package example1

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	AppIDQueryParam    = "app-id"
	AppLoginQueryParam = "login"
	AppPassQueryParam  = "pass"

	App1Token = "abc"
	App1Login = "App1_login"
	App1Pass  = "App1_pass"
	App2Token = "def"
	App2Login = "App2_login"
	App2Pass  = "App2_pass"
	App3Token = "ghi"
	App3Login = "App2_login"
	App3Pass  = "App2_pass"
)

type JwtToken struct {
	CreationTimestamp int64  `json:"creation_timestamp"`
	AccessToken       string `json:"access_token"`
	ResetToken        string `json:"reset_token"`
}

type authHandler struct {
	trustedApps       []app
	resourceServerUrl string
}

func NewAuthHandler(resourceServerUrl string) *authHandler {
	return &authHandler{
		trustedApps: []app{
			{token: App1Token, login: App1Login, password: App1Pass},
			{token: App2Token, login: App2Login, password: App2Pass},
			{token: App3Token, login: App3Login, password: App3Pass},
		},
		resourceServerUrl: resourceServerUrl,
	}
}

func (a authHandler) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	appID := request.URL.Query().Get(AppIDQueryParam)
	app, err := a.getApp(appID)
	if err != nil {
		http.Error(w, "Wrong app id", http.StatusBadRequest)
		return
	}

	if !app.auth(request.URL.Query().Get(AppLoginQueryParam), request.URL.Query().Get(AppPassQueryParam)) {
		http.Error(w, "Couldn't authenticated in app", http.StatusForbidden)
		return
	}

	jsonBuf, err := json.Marshal(&JwtToken{
		CreationTimestamp: time.Now().Unix(),
		AccessToken:       app.accessToken(),
		ResetToken:        app.refreshToken(),
	})

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	resp, err := http.Post(fmt.Sprintf("%s/add-token", a.resourceServerUrl), "application/json", bytes.NewBuffer(jsonBuf))
	if err != nil || resp.StatusCode != http.StatusOK {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonBuf)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

}

func (a authHandler) getApp(appID string) (*app, error) {
	for _, trustedApp := range a.trustedApps {
		if trustedApp.token == appID {
			return &trustedApp, nil
		}
	}
	return nil, errors.New("not found")
}

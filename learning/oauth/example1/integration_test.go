package example1_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	oa "go_sendbox/learning/oauth/example1"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuth(t *testing.T) {
	resServer := resourceServer()
	authServer := authServer(resServer.URL)

	client := http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/auth?%s", authServer.URL, appsQueryStrings()[0]), nil)
	require.NoError(t, err)
	r, err := client.Do(req)
	require.NoError(t, err)

	jwt := new(oa.JwtToken)
	err = json.NewDecoder(r.Body).Decode(&jwt)
	require.NoError(t, err)

	jsonBuf, err := json.Marshal(jwt)
	req, err = http.NewRequest("GET", fmt.Sprintf("%s/resource", resServer.URL), bytes.NewBuffer(jsonBuf))
	require.NoError(t, err)
	r, err = client.Do(req)
	require.NoError(t, err)

	resource := new(oa.Resource)
	err = json.NewDecoder(r.Body).Decode(resource)
	assert.True(t, resource.StatusOK)
	require.NoError(t, err)
}

func appsQueryStrings() []string {
	return []string{
		fmt.Sprintf("%s=%s&%s=%s&%s=%s", oa.AppIDQueryParam, oa.App1Token, oa.AppLoginQueryParam, oa.App1Login, oa.AppPassQueryParam, oa.App1Pass),
		fmt.Sprintf("%s=%s&%s=%s&%s=%s", oa.AppIDQueryParam, oa.App2Token, oa.AppLoginQueryParam, oa.App2Login, oa.AppPassQueryParam, oa.App2Pass),
		fmt.Sprintf("%s=%s&%s=%s&%s=%s", oa.AppIDQueryParam, oa.App3Token, oa.AppLoginQueryParam, oa.App3Login, oa.AppPassQueryParam, oa.App3Pass),
	}
}

func authServer(resourceServerURL string) *httptest.Server {
	return httptest.NewServer(oa.NewAuthHandler(resourceServerURL))
}

func resourceServer() *httptest.Server {
	mux := http.NewServeMux()
	mux.Handle("/add-token", oa.NewAddTokenHandler())
	mux.Handle("/resource", oa.NewResourceHandler())
	return httptest.NewServer(mux)
}

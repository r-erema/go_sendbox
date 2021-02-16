package example1

import (
	"crypto/sha256"
	"fmt"
	"time"
)

type app struct {
	token, login, password string
}

func (a app) auth(login, password string) bool {
	return a.login == login && a.password == password
}

func (a app) accessToken() string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%d%s%s%s",
		time.Now().Unix(),
		a.token,
		a.login,
		a.password,
	))))
}
func (a app) refreshToken() string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%s%d%s%s",
		a.password,
		time.Now().Unix(),
		a.token,
		a.login,
	))))
}

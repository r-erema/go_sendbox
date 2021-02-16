package password

import (
	"math/rand"
	"time"
)

func generatePassword(length int) (password string) {
	rand.Seed(int64(time.Now().Nanosecond()))

	var passwords [60]string
	for j := 0; j < 60; j++ {
		var tempPass string
		for i := 0; i < length; i++ {
			n := rand.Intn(128)
			tempPass += string([]byte{byte(n)})
		}
		passwords[j] = tempPass
	}
	password = passwords[time.Now().Second()]
	return
}

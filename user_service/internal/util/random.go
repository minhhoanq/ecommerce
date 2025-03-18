package util

import (
	"fmt"
	"math/rand"
	"strings"
)

const alphabet = "abcdfeghijklmnopqrstuvwxyz"

// RandomString generates a random string of length n
func RamdomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomEmail generates a random email
func RandomEmail() string {
	return fmt.Sprintf("%s@gmail.com", RamdomString(6))
}

// RandomUsername generates a random username
func RandomUsername() string {
	return RamdomString(6)
}

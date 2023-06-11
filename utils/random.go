package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyz")

	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int) int {
	return rand.Intn(max-min) + min
}

// RandomOwner generates a random owner name
func RandomOwner() string {
	return RandomString(6)
}

// Generate Random Email Address for user
func RandomEmail() string {
	return fmt.Sprintf("%s@payzonex.com", RandomString(6))
}

// Split Strings
func SplitStrings(s string) []string {
	var r []string
	for _, v := range s {
		r = append(r, string(v))
	}
	return r
}

// RandomCryptoPublicKeyAddress generates a random User Eth Address
func RandomCryptoPublicKeyAddress() string {
	return fmt.Sprintf("0x%s", RandomString(10))
}

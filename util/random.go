package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var random *rand.Rand

var chars = "abcdefghijklmnopqrstuvwxyz"

func init() {
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func RandomInt(min, max int) int {
	return min + int(random.Intn(max-min))
}

func RandomString(n int) string {
	var sb strings.Builder
	for i := 0; i < n; i++ {
		sb.WriteByte(chars[RandomInt(0, len(chars))])
	}
	return sb.String()
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomEmail() string {
	return fmt.Sprintf("%s@gmail.com", RandomString(6))
}

func RandomBalance() int64 {
	return int64(RandomInt(0, 1000))
}

func RandomCurrency() string {
	currencies := []string{"USD", "CAD", "INR"}
	return currencies[RandomInt(0, len(currencies))]
}

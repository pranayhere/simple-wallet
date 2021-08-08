package util

import (
    "fmt"
    "math/rand"
    "strings"
    "time"
)

const alphabets = "abcdefghijklmnopqrstuvwxyz"

func init() {
    rand.Seed(time.Now().UnixNano())
}

// RandomInt generate random integer between min and max
func RandomInt(min, max int64) int64 {
    return min + rand.Int63n(max-min+1)
}

// RandomString generates a random string of length n
func RandomString(n int) string {
    var sb strings.Builder
    k := len(alphabets)

    for i := 0; i < n; i++ {
        c := alphabets[rand.Intn(k)]
        sb.WriteByte(c)
    }

    return sb.String()
}

// RandomUser generates a random owner name
func RandomUser() string {
    return RandomString(6)
}

// RandomMoney generates random amount of money
func RandomMoney() int64 {
    return RandomInt(0, 1000)
}

// RandomEmail generate random email id
func RandomEmail() string {
    email := fmt.Sprintf("%s@email.com", RandomUser())
    return email
}

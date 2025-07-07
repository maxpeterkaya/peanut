package common

import "math/rand"

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var letterNumbers = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func GenerateUsername() string {
	return "admin"
}

func GeneratePassword(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func GenerateKey(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterNumbers[rand.Intn(len(letterNumbers))]
	}
	return string(b)
}

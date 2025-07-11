package common

import (
	"crypto/rand"
	"github.com/rs/zerolog/log"
	"math/big"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var letterNumbers = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func GenerateUsername() string {
	return "admin"
}

func GeneratePassword(n int) string {
	s := make([]rune, n)
	for i := range s {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			log.Fatal().Err(err).Msg("failed to generate password")
		}
		s[i] = letters[idx.Int64()]
	}
	return string(s)
}

func GenerateKey(n int) string {
	s := make([]rune, n)
	for i := range s {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(letterNumbers))))
		if err != nil {
			log.Fatal().Err(err).Msg("failed to generate key")
		}
		s[i] = letterNumbers[idx.Int64()]
	}
	return string(s)
}

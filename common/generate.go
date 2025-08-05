package common

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"math/big"
	ran "math/rand"
	"net/http"
	"strings"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var letterNumbers = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func GenerateUsername() string {
	p := ran.Intn(9999-1000) + 1000

	resp, err := http.Get("https://random-word-api.herokuapp.com/word")
	if err != nil {
		log.Err(err).Msg("Error generating username")
		return fmt.Sprintf("admin%d", p)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Err(err).Msg("error reading response")
		return fmt.Sprintf("admin%d", p)
	}

	var word []string
	err = json.Unmarshal(body, &word)
	if err != nil {
		log.Err(err).Msg("error unmarshalling response")
		return fmt.Sprintf("admin%d", p)
	}

	return fmt.Sprintf("%s%d", toTitle(word[0]), p)
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

func toTitle(str string) string {
	letters := strings.Split(str, "")
	return strings.ToUpper(letters[0]) + strings.Join(letters[1:], "")
}

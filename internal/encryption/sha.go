package encryption

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/rs/zerolog/log"
	"io"
	"os"
)

func SHAHash(input string) string {
	h := sha256.New()
	h.Write([]byte(input))
	return hex.EncodeToString(h.Sum(nil))
}

func SHAFileHash(filepath string) (string, error) {
	f, err := os.Open(filepath)
	if err != nil {
		log.Error().Err(err).Str("encryption", "sha256").Msg("failed to open file")
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Error().Err(err).Str("encryption", "sha256").Msg("failed to hash file")
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

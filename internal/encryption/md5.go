package encryption

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/rs/zerolog/log"
	"io"
	"os"
)

func MD5Hash(text string) string {
	h := md5.New()
	h.Write([]byte(text))
	return hex.EncodeToString(h.Sum(nil))
}

func MD5HashFile(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		log.Error().Err(err).Str("encryption", "md5").Msg("failed to open file")
		return "", err
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Error().Err(err).Str("encryption", "md5").Msg("failed to hash file")
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

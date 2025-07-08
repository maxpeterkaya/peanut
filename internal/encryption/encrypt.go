package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"github.com/rs/zerolog/log"
	"io"
	"peanut/internal/config"
)

func EncryptText(data string) (string, error) {
	encryptKey := []byte(config.Config.EncryptionKey)

	text := []byte(data)

	c, err := aes.NewCipher(encryptKey)
	if err != nil {
		log.Err(err).Msg("Error creating AES cipher")
		log.Trace()
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		log.Err(err).Msg("Error creating GCM")
		log.Trace()
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		log.Err(err).Msg("Error creating nonce")
		log.Trace()
		return "", err
	}

	encrypted := hex.EncodeToString(gcm.Seal(nonce, nonce, text, nil))

	return encrypted, nil
}

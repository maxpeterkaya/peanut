package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"github.com/rs/zerolog/log"
	"peanut/internal/config"
)

func DecryptText(text string) (string, error) {
	encryptKey := []byte(config.Config.EncryptionKey)

	ciphertext, err := hex.DecodeString(text)
	if err != nil {
		log.Err(err).Msg("Error decoding text")
	}

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

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		log.Err(nil).Msg("ciphertext too short")
		log.Trace()
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Err(err).Msg("Error decrypting text")
		log.Trace()
		return "", err
	}

	return string(plaintext), nil
}

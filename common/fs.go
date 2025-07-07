package common

import (
	"errors"
	"os"
)

func Exists(path string) bool {
	_, err := os.Stat(path)

	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

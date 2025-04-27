package utils

import (
	"crypto/md5"
	"fmt"
	"io"
)

func HashPassword(password, salt string) (string, error) {
	if len(password) == 0 || len(salt) == 0 {
		return "", fmt.Errorf("generate password with salt failed")
	}
	h := md5.New()
	_, err := io.WriteString(h, password+salt)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

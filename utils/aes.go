package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

var (
	key = []byte("rl2ChY4d1wzfWxYE")
)

// AesEncryptGCM encrypts plaintext using AES in GCM mode.
func AesEncryptGCM(plaintext []byte) ([]byte, error) {
	var (
		err   error
		gcm   cipher.AEAD
		block cipher.Block
	)
	if block, err = aes.NewCipher(key); err != nil {
		return nil, err
	}
	if gcm, err = cipher.NewGCM(block); err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return []byte(hex.EncodeToString(gcm.Seal(nonce, nonce, plaintext, nil))), nil
}

// AesDecryptGCM decrypts ciphertext using AES in GCM mode.
func AesDecryptGCM(encryptedHex []byte) ([]byte, error) {
	var (
		err               error
		nonce, ciphertext []byte
		gcm               cipher.AEAD
		block             cipher.Block
	)
	if ciphertext, err = hex.DecodeString(string(encryptedHex)); err != nil {
		return nil, err
	}
	if block, err = aes.NewCipher(key); err != nil {
		return nil, err
	}
	if gcm, err = cipher.NewGCM(block); err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	nonce, ciphertext = ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

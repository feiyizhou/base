package test

import (
	"fmt"
	"testing"

	"github.com/feiyizhou/base/utils"
)

func Test_RandCharStr(t *testing.T) {
	fmt.Println(utils.RandStr(22, utils.AllCharArr))
}

func Test_RandNumStr(t *testing.T) {
	fmt.Println(utils.RandStr(8, utils.NumCharArr))
}

func Test_AesGCM(t *testing.T) {
	plaintext := "Hello, World!"
	fmt.Println("Plaintext:", plaintext)

	ciphertext, err := utils.AesEncryptGCM([]byte(plaintext))
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}
	fmt.Println("Ciphertext:", string(ciphertext))

	decrypted, err := utils.AesDecryptGCM([]byte(ciphertext))
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}
	fmt.Println("Decrypted:", string(decrypted))

	if string(decrypted) != plaintext {
		t.Fatalf("Decrypted text does not match original plaintext")
	}
}

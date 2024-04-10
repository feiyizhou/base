package utils

import (
	"encoding/base64"

	"github.com/skip2/go-qrcode"
)

// ContentToQRCodeStr 根据content返回二维码图片的base64
func ContentToQRCodeStr(content string, level qrcode.RecoveryLevel, size int) (string, error) {
	var (
		pngBytes []byte
		err      error
	)
	pngBytes, err = qrcode.Encode(content, level, size)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(pngBytes), err
}

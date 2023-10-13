package utils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"strings"
)

func RandomBytes(nbBytes int) []byte {
	b := make([]byte, nbBytes)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}

func RandomHex(nbBytes int) string {
	b := RandomBytes(nbBytes)
	return strings.ToUpper(hex.EncodeToString(b))
}

func RandomBase64(nbBytes int) string {
	b := RandomBytes(nbBytes)
	return base64.RawURLEncoding.EncodeToString(b)
}

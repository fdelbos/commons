package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"

	"github.com/google/uuid"
)

type (
	APIKey struct {
		UserID      uuid.UUID
		Digest      []byte
		DisplayName string
	}
)

const (
	APIKeyLength = 32
)

func NewApiKey() (string, error) {
	key := new([APIKeyLength]byte)
	_, err := rand.Read(key[:])
	if err != nil {
		return "", err
	}

	keyb64 := base64.RawURLEncoding.EncodeToString(key[:])
	return keyb64, nil
}

func DigestFromAPIKey(keyb64 string) ([]byte, error) {
	key, err := base64.RawURLEncoding.DecodeString(keyb64)
	if err != nil {
		return nil, err
	}

	h := sha256.New()
	h.Write(key)
	return h.Sum(nil), nil
}

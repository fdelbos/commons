package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomBytes(t *testing.T) {
	b := RandomBytes(32)
	assert.Equal(t, 32, len(b))
	assert.NotEqual(t, b, RandomBytes(32))

}

func TestRandomHex(t *testing.T) {
	b := RandomHex(32)
	assert.Greater(t, len(b), 32)
	assert.NotEqual(t, b, RandomHex(32))
}

func TestRandomBase64(t *testing.T) {
	b := RandomBase64(32)
	assert.Greater(t, len(b), 32)
	assert.NotEqual(t, b, RandomBase64(32))
}

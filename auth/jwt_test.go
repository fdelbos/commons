package auth_test

import (
	"log"
	"testing"
	"time"

	. "github.com/fdelbos/commons/auth"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestJWT(t *testing.T) {
	pub, priv, err := NewJWTKeyPair()
	assert.NoError(t, err)

	subject := uuid.New().String()

	// test everything is ok
	{
		issuerJWT, err := NewJWTIssuer(priv)
		assert.NoError(t, err)
		token, err := issuerJWT.Issue(subject)
		assert.NoError(t, err)
		log.Print(token)

		validator, err := NewJWTValidator(pub)
		assert.NoError(t, err)

		sub, err := GetProvisionmalSubject(token)
		assert.NoError(t, err)
		assert.Equal(t, subject, sub)

		res, err := validator.Validate(token, time.Minute)
		assert.NoError(t, err)
		assert.Equal(t, subject, res)
	}

}

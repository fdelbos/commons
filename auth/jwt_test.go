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
	issuer := uuid.New().String()
	audience := uuid.New().String()

	// test everything is ok
	{
		issuerJWT, err := NewJWTIssuer(issuer, priv)
		assert.NoError(t, err)
		token, err := issuerJWT.Issue(time.Hour, subject, audience)
		assert.NoError(t, err)
		log.Print(token)

		audienceJWT, err := NewJWTAudience(pub, audience)
		assert.NoError(t, err)

		sub, err := GetProvisionmalSubject(token)
		assert.NoError(t, err)
		assert.Equal(t, subject, sub)

		res, err := audienceJWT.Validate(token)
		assert.NoError(t, err)
		assert.Equal(t, subject, res)
	}

}

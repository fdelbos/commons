package auth

import (
	"crypto"
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"time"

	"github.com/fdelbos/commons/utils"
	"github.com/golang-jwt/jwt/v5"
)

type (
	JWTIssuer struct {
		method  jwt.SigningMethod
		privKey crypto.PrivateKey
		issuer  string
	}

	JWTAudience struct {
		method   jwt.SigningMethod
		pubKey   crypto.PublicKey
		audience string
	}
)

const (
	MaxClockSkew = time.Minute
)

var (
	defaultMethod = jwt.SigningMethodEdDSA

	ErrInvalidSignature = errors.New("invalid signature")
)

func NewJWTKeyPair() ([]byte, []byte, error) {
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, nil, err
	}

	b, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return nil, nil, err
	}

	block := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: b,
	}

	privKey := pem.EncodeToMemory(block)

	b, err = x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil, nil, err
	}

	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: b,
	}
	pubKey := pem.EncodeToMemory(block)

	return pubKey, privKey, nil
}

// NewJWT returns a new JWT instance.
func NewJWTIssuer(issuer string, privKey []byte) (*JWTIssuer, error) {
	key, err := jwt.ParseEdPrivateKeyFromPEM(privKey)
	if err != nil {
		return nil, err
	}
	return &JWTIssuer{
		method:  jwt.SigningMethodEdDSA,
		issuer:  issuer,
		privKey: key,
	}, nil
}

func (j *JWTIssuer) Issue(ttl time.Duration, subject string, audiences ...string) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Issuer:    j.issuer,
		Subject:   subject,
		Audience:  audiences,
		ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		NotBefore: jwt.NewNumericDate(now.Add(-MaxClockSkew)),
		IssuedAt:  jwt.NewNumericDate(now),
		ID:        utils.RandomBase64(16),
	}
	token := jwt.NewWithClaims(j.method, claims)
	return token.SignedString(j.privKey)
}

func NewJWTAudience(pubKey []byte, audience string) (*JWTAudience, error) {
	key, err := jwt.ParseEdPublicKeyFromPEM(pubKey)
	if err != nil {
		return nil, err
	}
	return &JWTAudience{
		method:   jwt.SigningMethodEdDSA,
		pubKey:   key,
		audience: audience,
	}, nil
}

func (j *JWTAudience) Validate(token string) (string, error) {
	res, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if token.Method != j.method {
			return nil, ErrInvalidSignature
		}
		return j.pubKey, nil
	})
	if err != nil {
		return "", err
	}
	return res.Claims.GetSubject()
}

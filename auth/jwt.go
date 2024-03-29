package auth

import (
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type (
	JWTIssuer struct {
		method  jwt.SigningMethod
		privKey crypto.PrivateKey
	}

	JWTValidator struct {
		method jwt.SigningMethod
		pubKey crypto.PublicKey
	}
)

const (
	MaxClockSkew  = time.Minute
	RSA256KeySize = 2048
)

var (
	ErrInvalid = errors.New("invalid or expired token")
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

func NewRSAJWTKeyPair() ([]byte, []byte, error) {
	privKey, err := rsa.GenerateKey(rand.Reader, RSA256KeySize)
	if err != nil {
		return nil, nil, err
	}
	privKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privKey),
		},
	)

	pubKey := privKey.PublicKey
	pubKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(&pubKey),
		},
	)

	return privKeyPEM, pubKeyPEM, nil
}

// NewJWT returns a new JWT instance.
func NewJWTIssuer(privKey []byte) (*JWTIssuer, error) {
	key, err := jwt.ParseEdPrivateKeyFromPEM(privKey)
	if err != nil {
		return nil, err
	}
	return &JWTIssuer{
		method:  jwt.SigningMethodEdDSA,
		privKey: key,
	}, nil
}

func (j *JWTIssuer) Issue(subject string) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:  subject,
		IssuedAt: jwt.NewNumericDate(now),
	}
	token := jwt.NewWithClaims(j.method, claims)
	return token.SignedString(j.privKey)
}

func NewJWTValidator(pubKey []byte) (*JWTValidator, error) {
	key, err := jwt.ParseEdPublicKeyFromPEM(pubKey)
	if err != nil {
		return nil, err
	}
	return &JWTValidator{
		method: jwt.SigningMethodEdDSA,
		pubKey: key,
	}, nil
}

func NewJWTValidatorWithMethod(pubKey []byte, method jwt.SigningMethod) (*JWTValidator, error) {
	key, err := jwt.ParseEdPublicKeyFromPEM(pubKey)
	if err != nil {
		return nil, err
	}
	return &JWTValidator{
		method: method,
		pubKey: key,
	}, nil
}

func cleanToken(token string) string {
	token = strings.TrimPrefix(token, "Bearer ")
	return strings.TrimSpace(token)
}

func GetProvisionalSubject(token string) (string, error) {
	token = cleanToken(token)
	if token == "" {
		return "", ErrInvalid
	}

	dest := jwt.RegisteredClaims{}
	res, parts, err := jwt.NewParser().ParseUnverified(token, &dest)
	if err != nil {
		return "", err
	}
	if len(parts) != 3 {
		return "", ErrInvalid
	}
	sub, err := res.Claims.GetSubject()
	if err != nil {
		return "", ErrInvalid
	}
	return sub, nil
}

func (j *JWTValidator) Validate(token string, ttl time.Duration) (string, error) {
	token = cleanToken(token)

	res, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if token.Method != j.method {
			return nil, ErrInvalid
		}
		return j.pubKey, nil
	})
	if err != nil {
		return "", ErrInvalid
	}
	subject, err := res.Claims.GetSubject()
	if err != nil {
		return "", ErrInvalid
	}
	issuedAt, err := res.Claims.GetIssuedAt()
	if err != nil {
		return "", ErrInvalid
	}
	if issuedAt == nil {
		return "", ErrInvalid
	}

	if time.Now().After(issuedAt.Add(ttl + MaxClockSkew)) {
		return "", ErrInvalid
	}
	if time.Now().Before(issuedAt.Add(-MaxClockSkew)) {
		return "", ErrInvalid
	}

	return subject, nil
}

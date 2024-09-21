package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	minLengthSecretKey = 6
)

var (
	signingMethod     = jwt.SigningMethodHS256
	ErrSecretKeyLen   = fmt.Errorf("secret Key length should be at least %d", minLengthSecretKey)
	ErrClaimsTampered = fmt.Errorf("claims from token don't match, possibly token tampering")
)

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minLengthSecretKey {
		return nil, ErrSecretKeyLen
	}
	return &JWTMaker{secretKey}, nil
}

func (jwtMaker *JWTMaker) Create(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}
	t := jwt.NewWithClaims(signingMethod, payload)
	s, err := t.SignedString([]byte(jwtMaker.secretKey))
	return s, err
}

func (jwtMaker *JWTMaker) Verify(token string) (*Payload, error) {

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if token.Method != signingMethod {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtMaker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)

	if err != nil {
		return nil, err
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrClaimsTampered
	}

	return payload, nil

}

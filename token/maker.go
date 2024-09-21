package token

import (
	"errors"
	"time"
)

var (
	ErrTokenExpired        = errors.New("token expired")
	ErrKeySizeNotSupported = errors.New("unsupported secret key size")
	ErrInvalidToken        = errors.New("invalid token")
)

type Maker interface {
	Create(username string, duration time.Duration) (string, error)
	Verify(string) (*Payload, error)
}

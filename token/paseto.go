package token

import (
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey string
}

// Expiry implements Maker.
func (maker *PasetoMaker) Expiry() time.Duration {
	panic("unimplemented")
}

func NewPaseto(symmetricKey string) (Maker, error) {

	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, ErrKeySizeNotSupported
	}

	return &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: symmetricKey,
	}, nil
}

func (maker *PasetoMaker) Create(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}
	return paseto.NewV2().Encrypt([]byte(maker.symmetricKey), payload, nil)
}

func (maker *PasetoMaker) Verify(token string) (*Payload, error) {
	payload := &Payload{}
	err := paseto.NewV2().Decrypt(token, []byte(maker.symmetricKey), payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}
	err = payload.Valid()
	if err != nil {
		return nil, err
	}
	return payload, err
}

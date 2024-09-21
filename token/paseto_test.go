package token

import (
	"testing"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/stretchr/testify/require"
	"himavisoft.simple_bank/util"
)

func TestPasetoCreate(t *testing.T) {

	// Create the token

	paseto, err := NewPaseto(util.RandomString(chacha20poly1305.KeySize))
	require.NoError(t, err)

	randomUsername := util.RandomString(10) + "__username"
	randomDuration := time.Duration(util.RandomInt(100, 1000)) * time.Second
	token, err := paseto.Create(randomUsername, randomDuration)
	require.NoError(t, err)

	payload, err := paseto.Verify(token)
	require.NoError(t, err)

	require.Equal(t, randomUsername, payload.Username)
	require.Equal(t, payload.ExpiredAt, payload.IssuedAt.Add(randomDuration))
	require.True(t, payload.ExpiredAt.After(time.Now()))

}

func TestPasetoVerifyWithEmptyToken(t *testing.T) {

	// Create the token

	paseto, err := NewPaseto(util.RandomString(chacha20poly1305.KeySize))
	require.NoError(t, err)

	randomUsername := ""
	randomDuration := time.Duration(util.RandomInt(100, 1000)) * time.Second
	token, err := paseto.Create(randomUsername, randomDuration)
	require.NoError(t, err)

	payload, err := paseto.Verify(token)
	require.NoError(t, err)

	require.Equal(t, randomUsername, payload.Username)
	require.Equal(t, payload.ExpiredAt, payload.IssuedAt.Add(randomDuration))
	require.True(t, payload.ExpiredAt.After(time.Now()))

}

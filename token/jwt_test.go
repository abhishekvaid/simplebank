package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"himavisoft.simple_bank/util"
)

func TestJWTCreate(t *testing.T) {

	// Create the token

	jwtMaker, err := NewJWTMaker("abcdefghijklmnopqrstuvwxyz")
	require.NoError(t, err)

	randomUsername := util.RandomString(10) + "__username"
	randomDuration := time.Duration(util.RandomInt(100, 1000)) * time.Second
	token, err := jwtMaker.Create(randomUsername, randomDuration)
	require.NoError(t, err)

	payload, err := jwtMaker.Verify(token)
	require.NoError(t, err)

	require.Equal(t, randomUsername, payload.Username)
	require.Equal(t, payload.ExpiredAt, payload.IssuedAt.Add(randomDuration))
	require.True(t, payload.ExpiredAt.After(time.Now()))

}

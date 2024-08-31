package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"himavisoft.simple_bank/util"
)

func createTestUser(t *testing.T) User {

	arg := GetUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: "secret",
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.NotEmpty(t, user.Username)

	require.Equal(t, user.Username, arg.Username)
	require.Equal(t, user.HashedPassword, arg.HashedPassword)
	require.Equal(t, user.FullName, arg.FullName)
	require.Equal(t, user.Email, arg.Email)

	require.Zero(t, user.PasswordChangedAt)
	require.NotZero(t, user.CreatedAt)

	return user

}

func TestCreateUser(t *testing.T) {
	createTestUser(t)
}

func TestGetUser(t *testing.T) {

	newUser := createTestUser(t)

	user, err := testQueries.GetUser(context.Background(), newUser.Username)

	require.NoError(t, err)
	require.Equal(t, newUser, user)

}

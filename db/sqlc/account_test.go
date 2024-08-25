package db

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"

	"himavisoft.simple_bank/util"
)

func createTestAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  100, // util.RandomBalance(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.NotEmpty(t, account.ID)

	require.Equal(t, account.Owner, arg.Owner)
	require.Equal(t, account.Balance, arg.Balance)
	require.Equal(t, account.Currency, arg.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account

}

func TestCreateAccount(t *testing.T) {
	createTestAccount(t)
}

func TestGetAccount(t *testing.T) {
	acc1 := createTestAccount(t)

	acc2, err := testQueries.GetAccount(context.Background(), acc1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, acc2)

	require.EqualValues(t, acc1, acc2)

}

func TestGetAccounts(t *testing.T) {
	for i := 0; i < 5; i++ {
		createTestAccount(t)
	}

	arg := GetAccountsParams{
		Limit:  4,
		Offset: 0,
	}

	accounts, err := testQueries.GetAccounts(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, accounts, 4)

}

func TestDeleteAccount(t *testing.T) {
	acc1 := createTestAccount(t)

	err := testQueries.DeleteAccount(context.Background(), acc1.ID)
	require.NoError(t, err)

	acc2, err := testQueries.GetAccount(context.Background(), acc1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, acc2)
}

func TestUpdateAccount(t *testing.T) {
	acc1 := createTestAccount(t)

	arg := UpdateAccountParams{
		ID:      acc1.ID,
		Balance: acc1.Balance + 100,
	}

	acc2, err := testQueries.UpdateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, acc2)

	require.Equal(t, acc1.Owner, acc2.Owner)
	require.Equal(t, acc1.Currency, acc2.Currency)
	require.Equal(t, acc1.Balance+100, acc2.Balance)
	require.Equal(t, acc1.CreatedAt, acc2.CreatedAt)
	require.Equal(t, acc1.ID, acc2.ID)

}

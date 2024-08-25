package db

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {

	s := NewStore(testDB)

	var n int64 = 10
	var amount int64 = 10
	fromAccountInitial := createTestAccount(t)
	toAccountInitial := createTestAccount(t)

	errs := make(chan error)
	results := make(chan TransferTxResults)

	for i := 0; i < int(n); i++ {
		go func() {
			txId := fmt.Sprintf("%sTX:%d", strings.Repeat(" ", i*4), i+1)
			r, err := s.TransferTx(context.WithValue(context.Background(), txKey, txId), TransferTxParams{
				FromAccountId: fromAccountInitial.ID,
				ToAccountId:   toAccountInitial.ID,
				Amount:        int64(amount),
			})

			errs <- err
			results <- r
		}()
	}

	for i := 0; i < int(n); i++ {

		err := <-errs
		require.NoError(t, err)

		r := <-results

		transfer := r.Transfer

		// Test Transfer Object
		require.NotEmpty(t, transfer)
		require.NotEmpty(t, transfer.ID)
		require.NotEmpty(t, transfer.CreatedAt)
		require.Equal(t, transfer.FromAccountID, fromAccountInitial.ID)
		require.Equal(t, transfer.ToAccountID, toAccountInitial.ID)
		require.Equal(t, transfer.Amount, amount)

		// Test Entry Objects
		fromEntry := r.FromEntry
		require.NotEmpty(t, fromEntry)
		require.NotEmpty(t, fromEntry.ID)
		require.NotEmpty(t, fromEntry.CreatedAt)
		require.Equal(t, fromEntry.AccountID, fromAccountInitial.ID)
		require.Equal(t, fromEntry.Amount, -amount)

		fromEntryFromDB, err := s.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)
		require.EqualValues(t, fromEntryFromDB, fromEntry)

		toEntry := r.ToEntry
		require.NotEmpty(t, toEntry)
		require.NotEmpty(t, toEntry.ID)
		require.NotEmpty(t, toEntry.CreatedAt)
		require.Equal(t, toEntry.AccountID, toAccountInitial.ID)
		require.Equal(t, toEntry.Amount, amount)

		toEntryFromDB, err := s.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)
		require.EqualValues(t, toEntryFromDB, toEntry)

		// Test Account Objects

		require.NotEmpty(t, r.FromAccount)
		require.NotEmpty(t, r.ToAccount)

		require.Greater(t, fromAccountInitial.Balance, r.FromAccount.Balance)
		require.Equal(t, (fromAccountInitial.Balance-r.FromAccount.Balance)%amount, int64(0))
		require.GreaterOrEqual(t, (fromAccountInitial.Balance-r.FromAccount.Balance)/amount, int64(1))
		require.LessOrEqual(t, (fromAccountInitial.Balance-r.FromAccount.Balance)/amount, n)

		require.Less(t, toAccountInitial.Balance, r.ToAccount.Balance)
		require.Equal(t, (r.ToAccount.Balance-toAccountInitial.Balance)%amount, int64(0))
		require.GreaterOrEqual(t, (r.ToAccount.Balance-toAccountInitial.Balance)/amount, int64(1))
		require.LessOrEqual(t, (r.ToAccount.Balance-toAccountInitial.Balance)/amount, n)

	}

	fromAccountLive, err := s.GetAccount(context.Background(), fromAccountInitial.ID)
	require.NoError(t, err)
	require.Equal(t, -int64(n)*amount, fromAccountLive.Balance-fromAccountInitial.Balance)

	toAccountLive, err := s.GetAccount(context.Background(), toAccountInitial.ID)
	require.NoError(t, err)
	require.Equal(t, int64(n)*amount, toAccountLive.Balance-toAccountInitial.Balance)

}

func TestTransferTxWithDeadlock(t *testing.T) {
	s := NewStore(testDB)

	a1 := createTestAccount(t)
	a2 := createTestAccount(t)

	n := 10
	amount := int64(10)
	var err error

	errs := make(chan error)

	for i := 0; i < n; i++ {
		params := TransferTxParams{
			FromAccountId: a1.ID,
			ToAccountId:   a2.ID,
			Amount:        amount,
		}
		if i&1 == 0 {
			params = TransferTxParams{
				FromAccountId: a2.ID,
				ToAccountId:   a1.ID,
				Amount:        amount,
			}
		}
		go func() {
			_, err := s.TransferTx(context.Background(), params)
			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	a1Final, err := s.GetAccount(context.Background(), a1.ID)
	require.NoError(t, err)
	require.Equal(t, a1, a1Final)

	a2Final, err := s.GetAccount(context.Background(), a2.ID)
	require.NoError(t, err)
	require.Equal(t, a2, a2Final)

}

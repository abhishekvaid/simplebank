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

	n := 10
	var amount int64 = 10
	fromAccount := createTestAccount(t)
	toAccount := createTestAccount(t)

	errs := make(chan error)
	results := make(chan TransferTxResults)

	for i := 0; i < n; i++ {
		go func() {
			txId := fmt.Sprintf("%sTX:%d", strings.Repeat(" ", i*4), i+1)
			r, err := s.TransferTx(context.WithValue(context.Background(), txKey, txId), TransferTxParams{
				FromAccountId: fromAccount.ID,
				ToAccountId:   toAccount.ID,
				Amount:        int64(amount),
			})

			errs <- err
			results <- r
		}()
	}

	for i := 0; i < n; i++ {

		err := <-errs
		require.NoError(t, err)

		result := <-results

		transfer := result.Transfer

		// Test Transfer Object
		require.NotEmpty(t, transfer)
		require.NotEmpty(t, transfer.ID)
		require.NotEmpty(t, transfer.CreatedAt)
		require.Equal(t, transfer.FromAccountID, fromAccount.ID)
		require.Equal(t, transfer.ToAccountID, toAccount.ID)
		require.Equal(t, transfer.Amount, amount)

		// Test Entry Object
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.NotEmpty(t, fromEntry.ID)
		require.NotEmpty(t, fromEntry.CreatedAt)
		require.Equal(t, fromEntry.AccountID, fromAccount.ID)
		require.Equal(t, fromEntry.Amount, -amount)

		fromEntryFromDB, err := s.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)
		require.EqualValues(t, fromEntryFromDB, fromEntry)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.NotEmpty(t, toEntry.ID)
		require.NotEmpty(t, toEntry.CreatedAt)
		require.Equal(t, toEntry.AccountID, toAccount.ID)
		require.Equal(t, toEntry.Amount, amount)

		toEntryFromDB, err := s.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)
		require.EqualValues(t, toEntryFromDB, toEntry)

		// Test Account Object

		fromAccount_ := result.FromAccount
		require.NotEmpty(t, fromAccount_)

		// toAccount_ := result.ToAccount
		// require.NotEmpty(t, toAccount_)

		require.True(t, fromAccount.Balance > fromAccount_.Balance)
		require.Equal(t, (fromAccount.Balance-fromAccount_.Balance)%amount, int64(0))
		require.LessOrEqual(t, 1, int((fromAccount.Balance-fromAccount_.Balance)/amount))
		require.LessOrEqual(t, int((fromAccount.Balance-fromAccount_.Balance)/amount), n)

		// require.True(t, toAccount_.Balance > toAccount.Balance)
		// require.Equal(t, (toAccount_.Balance-toAccount.Balance)%amount, int64(0))
		// require.LessOrEqual(t, 1, int((toAccount_.Balance-toAccount.Balance)/amount))
		// require.LessOrEqual(t, int((toAccount_.Balance-toAccount.Balance)/amount), n)

	}

	fromAccountFromDB, err := s.GetAccount(context.Background(), fromAccount.ID)
	require.NoError(t, err)
	require.Equal(t, -int64(n)*amount, fromAccountFromDB.Balance-fromAccount.Balance)
	// toAccountFromDB, err := store.GetAccount(context.Background(), toAccount.ID)
	// require.NoError(t, err)
	// require.Equal(t, int64(n)*amount, toAccountFromDB.Balance-toAccount.Balance)
	// require.Equal(t, int64(2*n)*amount, toAccountFromDB.Balance-fromAccountFromDB.Balance)

}

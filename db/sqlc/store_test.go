package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	n := 5
	var amount int64 = 63
	fromAccount := createTestAccount(t)
	toAccount := createTestAccount(t)

	errs := make(chan error)
	results := make(chan TransferTxResults)

	for i := 0; i < n; i++ {
		go func() {
			r, err := store.TransferTx(context.Background(), TransferTxParams{
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

		require.NotEmpty(t, transfer)
		require.NotEmpty(t, transfer.ID)
		require.NotEmpty(t, transfer.CreatedAt)
		require.Equal(t, transfer.FromAccountID, fromAccount.ID)
		require.Equal(t, transfer.ToAccountID, toAccount.ID)
		require.Equal(t, transfer.Amount, amount)

		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.NotEmpty(t, fromEntry.ID)
		require.NotEmpty(t, fromEntry.CreatedAt)
		require.Equal(t, fromEntry.AccountID, fromAccount.ID)
		require.Equal(t, fromEntry.Amount, -amount)

		fromEntryFromDB, err := store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)
		require.EqualValues(t, fromEntryFromDB, fromEntry)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.NotEmpty(t, toEntry.ID)
		require.NotEmpty(t, toEntry.CreatedAt)
		require.Equal(t, toEntry.AccountID, toAccount.ID)
		require.Equal(t, toEntry.Amount, amount)

		toEntryFromDB, err := store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)
		require.EqualValues(t, toEntryFromDB, toEntry)

	}

}

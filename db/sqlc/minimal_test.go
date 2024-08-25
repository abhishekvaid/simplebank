package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeductTx(t *testing.T) {
	n := 50
	var amount int64 = 10
	account := createTestAccount(t)

	errs := make(chan error)
	results := make(chan DeductTxResult)

	for i := 0; i < n; i++ {
		go func() {
			txId := fmt.Sprintf("TX:%d", i+1)
			result, err := DeductTxMinimal(context.WithValue(context.Background(), txKey, txId), account.ID, amount)
			errs <- err
			results <- result
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
		result := <-results
		require.Equal(t, result.AfterBalance-result.BeforeBalance, -amount)
		fmt.Printf("[%s] %d --> %d\n", result.TxId, result.BeforeBalance, result.AfterBalance)
	}

}

func TestTransferTxMinimal(t *testing.T) {

	s := NewStore(testDB)

	n := 50
	var amount int64 = 10
	fromAccount := createTestAccount(t)
	toAccount := createTestAccount(t)

	errs := make(chan error)
	results := make(chan TransferTxMinimalResult)

	for i := 0; i < n; i++ {
		go func() {
			txId := fmt.Sprintf("TX:%d", i+1)
			result, err := s.TransferTxMinimal(context.WithValue(context.Background(), txKey, txId), fromAccount.ID, toAccount.ID, amount)
			errs <- err
			results <- result
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results

		fromBefore, fromAfter, toBefore, toAfter := result.FromAccountBefore, result.FromAccountAfter, result.ToAccountBefore, result.ToAccountAfter

		require.Greater(t, fromBefore.Balance, fromAfter.Balance)
		require.Equal(t, fromBefore.Balance-fromAfter.Balance, amount)

		require.Greater(t, toAfter.Balance, toBefore.Balance)
		require.Equal(t, toAfter.Balance-toBefore.Balance, amount)
	}

	fromAccountFromDB, err := s.GetAccount(context.Background(), fromAccount.ID)
	require.NoError(t, err)
	require.Equal(t, -int64(n)*amount, fromAccountFromDB.Balance-fromAccount.Balance)
	toAccountFromDB, err := s.GetAccount(context.Background(), toAccount.ID)
	require.NoError(t, err)
	require.Equal(t, int64(n)*amount, toAccountFromDB.Balance-toAccount.Balance)
	require.Equal(t, int64(2*n)*amount, toAccountFromDB.Balance-fromAccountFromDB.Balance)

}

package db

import (
	"context"
	"database/sql"
	"fmt"
)

type DeductTxResult struct {
	BeforeBalance int64  `json:"before_balance"`
	AfterBalance  int64  `json:"after_balance"`
	TxId          string `json:"tx_id"`
}

func DeductTxMinimal(ctx context.Context, accountId, amount int64) (DeductTxResult, error) {

	var driverName string = "postgres"
	var dataSource string = "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable"

	testDB, _ := sql.Open(driverName, dataSource)

	tx, _ := testDB.BeginTx(ctx, nil)

	q := New(tx)

	r := DeductTxResult{}

	r.TxId = fmt.Sprintf("%v", ctx.Value(txKey))
	accountBeforeUpdate, err := q.GetAccountForUpdate(ctx, accountId)
	if err != nil {
		return r, err
	}
	// fmt.Printf("[%s]: balance BEFORE update : (%d)  \n", r.TxId, accountBeforeUpdate.Balance)

	accountAfterUpdate, err := q.UpdateAccount(ctx, UpdateAccountParams{
		ID:      accountId,
		Balance: accountBeforeUpdate.Balance - amount,
	})
	if err != nil {
		return r, err
	}

	// fmt.Printf("[%s]: balance AFTER update : (%d)  \n", r.TxId, accountAfterUpdate.Balance)

	r.BeforeBalance = accountBeforeUpdate.Balance
	r.AfterBalance = accountAfterUpdate.Balance

	tx.Commit()

	return r, err

}

type TransferTxMinimalResult struct {
	FromAccountBefore Account `json:"from_account_before"`
	FromAccountAfter  Account `json:"from_account_after"`
	ToAccountBefore   Account `json:"to_account_before"`
	ToAccountAfter    Account `json:"to_account_after"`
	TxId              string  `json:"tx_id"`
}

func (s *Store) TransferTxMinimal(ctx context.Context, fromAccountId, toAccountId, amount int64) (TransferTxMinimalResult, error) {

	// var driverName string = "postgres"
	// var dataSource string = "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable"

	// testDB, _ := sql.Open(driverName, dataSource)

	tx, _ := s.db.BeginTx(ctx, nil)

	q := New(tx)

	r := TransferTxMinimalResult{}

	r.TxId = fmt.Sprintf("%v", ctx.Value(txKey))

	var err error

	r.FromAccountBefore, err = q.GetAccountForUpdate(ctx, fromAccountId)
	if err != nil {
		return r, err
	}

	r.FromAccountAfter, err = q.UpdateAccount(ctx, UpdateAccountParams{
		ID:      r.FromAccountBefore.ID,
		Balance: r.FromAccountBefore.Balance - amount,
	})
	if err != nil {
		return r, err
	}

	r.ToAccountBefore, err = q.GetAccountForUpdate(ctx, toAccountId)
	if err != nil {
		return r, err
	}

	r.ToAccountAfter, err = q.UpdateAccount(ctx, UpdateAccountParams{
		ID:      r.ToAccountBefore.ID,
		Balance: r.ToAccountBefore.Balance + amount,
	})
	if err != nil {
		return r, err
	}

	tx.Commit()

	return r, err

}

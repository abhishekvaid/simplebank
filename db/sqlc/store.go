package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		Queries: New(db),
		db:      db,
	}
}

var txKey = struct{}{}

// This function creates exactly 1 transaction and passes that into the callback function
func (s *Store) execTx(ctx context.Context, callback func(*Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	qErr := callback(New(tx))

	if qErr != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			fmt.Errorf("query err = %v | rollback error = %v", qErr, rbErr)
		}
		return qErr
	}

	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountId int64 `json:"from_account_id"`
	ToAccountId   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResults struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx which does following 1.) creates 1 transfer record 2.) Two individual account entries 3.) deduct / add money in account records
func (s *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResults, error) {

	var r TransferTxResults

	tx, err := s.db.BeginTx(ctx, nil)

	if err != nil {
		tx.Rollback()
		return r, err
	}

	q := New(tx)

	txName := ctx.Value(txKey)

	r.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
		FromAccountID: arg.FromAccountId,
		ToAccountID:   arg.ToAccountId,
		Amount:        arg.Amount,
	})
	if err != nil {
		tx.Rollback()
		return r, err
	}

	r.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
		AccountID: arg.FromAccountId,
		Amount:    -1.0 * arg.Amount,
	})
	if err != nil {
		return r, err
	}

	r.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
		AccountID: arg.ToAccountId,
		Amount:    arg.Amount,
	})
	if err != nil {
		tx.Rollback()
		return r, err
	}

	fromAccount, err := q.GetAccountForUpdate(ctx, arg.FromAccountId)
	if err != nil {
		tx.Rollback()
		return r, err
	}
	fmt.Printf("[%s]: From GetAccount : (%d)  \n", txName, fromAccount.Balance)

	r.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
		ID:      fromAccount.ID,
		Balance: fromAccount.Balance - arg.Amount,
	})

	fmt.Printf("[%s]: From UpdateAccount : (%d)  \n", txName, r.FromAccount.Balance)

	if err != nil {
		tx.Rollback()
		return r, err
	}

	toAccount, err := q.GetAccountForUpdate(ctx, arg.ToAccountId)
	if err != nil {
		return r, err
	}

	fmt.Printf("[%s]: To GetAccount : (%d)  \n", txName, toAccount.Balance)

	r.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
		ID:      toAccount.ID,
		Balance: toAccount.Balance + arg.Amount,
	})

	fmt.Printf("[%s]: To UpdateAccount : (%d)  \n", txName, r.ToAccount.Balance)

	if err != nil {
		tx.Rollback()
		return r, err
	}

	tx.Commit()

	return r, err

}

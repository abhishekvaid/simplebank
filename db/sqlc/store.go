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

	q := New(tx)

	qErr := callback(q)

	if qErr != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("query err = %v | rollback error = %v", qErr, rbErr)
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

	err := s.execTx(ctx, func(q *Queries) error {

		var err error

		r.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountId,
			ToAccountID:   arg.ToAccountId,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		r.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountId,
			Amount:    -1.0 * arg.Amount,
		})
		if err != nil {
			return err
		}

		r.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountId,
			Amount:    arg.Amount,
		})
		if err != nil {

			return err
		}

		r.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			Amount: -arg.Amount,
			ID:     arg.FromAccountId,
		})

		if err != nil {
			return err
		}

		r.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			Amount: arg.Amount,
			ID:     arg.ToAccountId,
		})

		if err != nil {
			return err
		}

		return nil
	})

	return r, err

}

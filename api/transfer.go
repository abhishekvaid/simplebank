package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "himavisoft.simple_bank/db/sqlc"
)

type transferAccountRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required"`
	ToAccountID   int64  `json:"to_account_id" binding:"required"`
	Amount        int64  `json:"amount" binding:"required,min=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (s *Server) TransferAmount(ctx *gin.Context) {

	body := transferAccountRequest{}
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fromAccount, ok := s.ValidateAccountCurrency(ctx, body.FromAccountID, body.Currency)
	if !ok {
		return
	}

	toAccount, ok := s.ValidateAccountCurrency(ctx, body.ToAccountID, body.Currency)
	if !ok {
		return
	}

	arg := db.TransferTxParams{
		FromAccountId:    body.FromAccountID,
		FromAccountOwner: fromAccount.Owner,
		ToAccountId:      body.ToAccountID,
		ToAccountOwner:   toAccount.Owner,
		Amount:           body.Amount,
	}

	transfer, err := s.store.TransferTx(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, transfer)

}

func (s *Server) ValidateAccountCurrency(ctx *gin.Context, accountID int64, currency string) (*db.Account, bool) {

	arg := db.GetAccountParams{
		ID:    accountID,
		Owner: ctx.GetHeader(authorizedUserId),
	}

	account, err := s.store.GetAccount(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return nil, false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return nil, false
	}

	if account.Currency != currency {
		err = fmt.Errorf("[Account Currency Mismatch] account with id=(%v) doesn't have required currency (%v)", accountID, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return nil, false
	}
	return &account, true
}

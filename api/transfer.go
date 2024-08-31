package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "himavisoft.simple_bank/db/sqlc"
)

func (s *Server) TransferAmount(ctx *gin.Context) {

	type TransferAccountRequest struct {
		FromAccountID int64  `json:"from_account_id" binding:"required"`
		ToAccountID   int64  `json:"to_account_id" binding:"required"`
		Amount        int64  `json:"amount" binding:"required,min=0"`
		Currency      string `json:"currency" binding:"required,currency"`
	}

	body := TransferAccountRequest{}
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if ok := s.ValidateAccountCurrency(ctx, body.FromAccountID, body.Currency); !ok {
		return
	}

	if ok := s.ValidateAccountCurrency(ctx, body.ToAccountID, body.Currency); !ok {
		return
	}

	arg := db.TransferTxParams{
		FromAccountId: body.FromAccountID,
		ToAccountId:   body.ToAccountID,
		Amount:        body.Amount,
	}

	transfer, err := s.store.TransferTx(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, transfer)

}

func (s *Server) ValidateAccountCurrency(ctx *gin.Context, accountID int64, currency string) bool {
	account, err := s.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	if account.Currency != currency {
		err = fmt.Errorf("[Account Currency Mismatch] account with id=(%v) doesn't have required currency (%v)", accountID, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}
	return true
}

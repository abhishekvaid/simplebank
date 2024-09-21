package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "himavisoft.simple_bank/db/sqlc"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type getAccountRequest struct {
	ID int32 `uri:"id" binding:"required"`
}

type listAccountRequest struct {
	PageOffset int32 `form:"page_offset" binding:"required,min=1"`
	PageSize   int32 `form:"page_size" binding:"required,min=1,max=20"`
}

type createAccountRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
}

func (s *Server) GetAccount(ctx *gin.Context) {

	req := getAccountRequest{}
	if err := ctx.BindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.GetAccountParams{
		ID:    int64(req.ID),
		Owner: ctx.GetString(authorizedUserId),
	}

	account, err := s.store.GetAccount(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)

}

func (s *Server) ListAccounts(ctx *gin.Context) {

	req := listAccountRequest{}
	if err := ctx.BindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.ListAccountsParams{
		Owner:  ctx.GetHeader(authorizedUserId),
		Limit:  req.PageSize,
		Offset: req.PageOffset,
	}
	account, err := s.store.ListAccounts(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)

}

func (s *Server) CreateAccount(ctx *gin.Context) {

	req := createAccountRequest{}

	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateAccountParams{
		Owner:    ctx.GetHeader(authorizedUserId),
		Balance:  0,
		Currency: req.Currency,
	}

	account, err := s.store.CreateAccount(ctx, arg)
	if err != nil {
		// Make changes here since accounts now depend on user
		switch err.(*pq.Error).Code.Name() {
		case "foreign_key_violation":
			ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("account with owner=%v can't be created since owner with username doesn't exist", arg.Owner)))
			return
		case "unique_violation":
			ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("account with owner=%v already exists with currency=%v", arg.Owner, arg.Currency)))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)

}

package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "himavisoft.simple_bank/db/sqlc"

	_ "github.com/lib/pq"
)

func (s *Server) GeAccount(ctx *gin.Context) {
	type GetAccountParams struct {
		ID int32 `uri:"id" binding:"required"`
	}

	q := GetAccountParams{}

	if err := ctx.BindUri(&q); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := s.store.GetAccount(ctx, int64(q.ID))
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

func (s *Server) GeAccounts(ctx *gin.Context) {
	type GetAccountsParams struct {
		PageOffset int32 `form:"page_offset" binding:"required,min=1"`
		PageSize   int32 `form:"page_size" binding:"required,min=1,max=20"`
	}

	params := GetAccountsParams{}

	if err := ctx.BindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.GetAccountsParams{
		Limit:  params.PageSize,
		Offset: params.PageOffset,
	}

	account, err := s.store.GetAccounts(ctx, arg)
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

	type CreateContextParams struct {
		Owner    string `json:"owner" binding:"required"`
		Currency string `json:"currency" binding:"required,oneof=USD CAD INR"`
	}

	var reqDto CreateContextParams

	if err := ctx.ShouldBindBodyWithJSON(&reqDto); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateAccountParams{
		Owner:    reqDto.Owner,
		Balance:  0,
		Currency: reqDto.Currency,
	}

	account, err := s.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)

}

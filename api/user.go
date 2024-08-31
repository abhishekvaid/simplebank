package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "himavisoft.simple_bank/db/sqlc"

	_ "github.com/lib/pq"
)

func (s *Server) CreateUser(ctx *gin.Context) {

	type CreateUserReqDTO struct {
		Username       string `json:"username" binding:"required"`
		HashedPassword string `json:"hashed_password" binding:"required"`
		FullName       string `json:"full_name" binding:"required"`
		Email          string `json:"email" binding:"required"`
	}

	var reqDTO CreateUserReqDTO

	if err := ctx.ShouldBindBodyWithJSON(&reqDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.GetUserParams{
		Username:       reqDTO.Username,
		HashedPassword: reqDTO.HashedPassword,
		FullName:       reqDTO.FullName,
		Email:          reqDTO.Email,
	}

	user, err := s.store.CreateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)

}

func (s *Server) GetUser(ctx *gin.Context) {

	type GetUserDTO struct {
		Username string `json:"username" binding:"required"`
	}

	var reqDTO GetUserDTO

	if err := ctx.ShouldBindBodyWithJSON(&reqDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.GetUserParams{
		Username: reqDTO.Username,
	}

	user, err := s.store.CreateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)

}

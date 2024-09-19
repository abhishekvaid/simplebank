package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "himavisoft.simple_bank/db/sqlc"
	"himavisoft.simple_bank/util"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type createUserResponse struct {
	Username          string    `json:"username"`
	HashedPassword    string    `json:"-"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func (s *Server) CreateUser(ctx *gin.Context) {

	var reqDTO createUserRequest

	if err := ctx.ShouldBindBodyWithJSON(&reqDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(reqDTO.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("cannot generate hash of the password provided")))
	}

	arg := db.CreateUserParams{
		Username:       reqDTO.Username,
		HashedPassword: hashedPassword,
		FullName:       reqDTO.FullName,
		Email:          reqDTO.Email,
	}

	user, err := s.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, createUserResponse(user))

}

func (s *Server) GetUser(ctx *gin.Context) {

	type GetUserDTO struct {
		Username string `uri:"username" binding:"required"`
	}

	var reqDTO GetUserDTO

	if err := ctx.ShouldBindUri(&reqDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := s.store.GetUser(ctx, reqDTO.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)

}

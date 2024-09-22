package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "himavisoft.simple_bank/db/sqlc"
	"himavisoft.simple_bank/util"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

var (
	ErrLoginCreds    = errors.New("wrong login creds")
	ErrUserNotFound  = errors.New("username doesn't exist")
	ErrWrongPassword = errors.New("password doesn't match")
)

type createUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

// func createUserResponse(user db.User)

type userResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

type loginParams struct {
	Username string `json:"username" binding:"required,min=1"`
	Password string `json:"password" binding:"required,min=1"`
}

type loginResponse struct {
	Token string       `json:"token"`
	User  userResponse `json:"user"`
}

func createUserResponse(user db.User) *userResponse {
	return &userResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

func createLoginResponse(user db.User, token string) *loginResponse {
	return &loginResponse{
		Token: token,
		User:  *createUserResponse(user),
	}
}

func (s *Server) Login(ctx *gin.Context) {

	loginParams := loginParams{}

	if err := ctx.ShouldBindBodyWithJSON(&loginParams); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := s.store.GetUser(ctx, loginParams.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(ErrUserNotFound))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.CheckPassword(loginParams.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(ErrWrongPassword))
		return
	}

	token, err := s.tokenMaker.Create(loginParams.Username, s.config.TokenExpiry)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, createLoginResponse(user, token))

}

func (s *Server) CreateUser(ctx *gin.Context) {

	var req createUserRequest

	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("cannot generate hash of the password provided")))
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
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

	username := ctx.GetString(authorizedUserId)
	user, err := s.store.GetUser(ctx, username)
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

package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	db "himavisoft.simple_bank/db/sqlc"
	"himavisoft.simple_bank/util"

	_ "github.com/lib/pq"
)

var (
	ErrLoginCreds    = errors.New("wrong login creds")
	ErrUserNotFound  = errors.New("username doesn't exist")
	ErrWrongPassword = errors.New("password doesn't match")
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type userResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

type loginParams struct {
	Username string `json:"username" binding:"required,min=1,alphanum"`
	Password string `json:"password" binding:"required,min=1"`
}

type loginResponse struct {
	SessionID          uuid.UUID    `json:"session_id"`
	AccessToken        string       `json:"access_token"`
	AccessTokenExpiry  time.Time    `json:"access_token_expiry"`
	RefreshToken       string       `json:"refresh_token"`
	RefreshTokenExpiry time.Time    `json:"refresh_token_expiry"`
	User               userResponse `json:"user"`
}

func (s *Server) Login(ctx *gin.Context) {

	loginParams := loginParams{}

	if err := ctx.ShouldBindBodyWithJSON(&loginParams); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := s.store.GetUser(ctx, loginParams.Username)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
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

	accessToken, accessTokenPayload, err := s.tokenMaker.Create(loginParams.Username, s.config.TokenExpiry)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshTokenPayload, err := s.tokenMaker.Create(loginParams.Username, s.config.RefreshTokenExpiry)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	session, err := s.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshTokenPayload.ID,
		Username:     accessTokenPayload.Username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.Request.RemoteAddr,
		IsBlocked:    false,
		ExpireAt:     time.Time{},
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, loginResponse{
		SessionID:          session.ID,
		AccessToken:        accessToken,
		AccessTokenExpiry:  accessTokenPayload.ExpiredAt,
		RefreshToken:       refreshToken,
		RefreshTokenExpiry: refreshTokenPayload.ExpiredAt,
		User: userResponse{
			Username:          user.Username,
			FullName:          user.FullName,
			Email:             user.Email,
			PasswordChangedAt: user.PasswordChangedAt,
			CreatedAt:         user.CreatedAt,
		},
	})

}

// func (s *Server) RefreshToken(ctx *gin.Context) {

// 	loginParams := loginParams{}

// 	if err := ctx.ShouldBindBodyWithJSON(&loginParams); err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}

// 	user, err := s.store.GetUser(ctx, loginParams.Username)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			ctx.JSON(http.StatusNotFound, errorResponse(ErrUserNotFound))
// 			return
// 		}
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	err = util.CheckPassword(loginParams.Password, user.HashedPassword)
// 	if err != nil {
// 		ctx.JSON(http.StatusUnauthorized, errorResponse(ErrWrongPassword))
// 		return
// 	}

// 	accessToken, accessTokenPayload, err := s.tokenMaker.Create(loginParams.Username, s.config.TokenExpiry)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	refreshToken, refreshTokenPayload, err := s.tokenMaker.Create(loginParams.Username, s.config.RefreshTokenExpiry)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	s.store.CreateSession(ctx, db.CreateSessionParams{
// 		ID:           refreshTokenPayload.ID,
// 		Username:     accessTokenPayload.Username,
// 		RefreshToken: refreshToken,
// 		UserAgent:    ctx.Request.UserAgent(),
// 		ClientIp:     ctx.Request.RemoteAddr,
// 		IsBlocked:    false,
// 		ExpireAt:     time.Time{},
// 	})

// 	ctx.JSON(http.StatusOK, createLoginResponse(user, accessToken))

// }

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

	// user, err := s.store.CreateUser(ctx, arg)
	// if err != nil {
	// 	if pqErr, ok := err.(*pq.Error); ok {
	// 		switch pqErr.Code.Name() {
	// 		case "unique_violation":
	// 			ctx.JSON(http.StatusForbidden, errorResponse(err))
	// 			return
	// 		}
	// 	}
	// 	ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	// 	return
	// }

	user, err := s.store.CreateUser(ctx, arg)
	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, userResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	})

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

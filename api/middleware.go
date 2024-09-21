package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"himavisoft.simple_bank/token"
)

const (
	authorizationHeaderKey = "Authorization"
	authorizationType      = "Bearer"
	authorizedUserId       = "AuthorizedUserId"
)

var (
	ErrAuthNoHeader       = errors.New("authorization header not present in request")
	ErrAuthValWrongFormat = errors.New("authorization header value not in correct format")
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {

	return func(ctx *gin.Context) {

		header := ctx.GetHeader(authorizationHeaderKey)
		if header == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(ErrAuthNoHeader))
			return
		}

		headerTokens := strings.Fields(header)

		if len(headerTokens) != 2 || headerTokens[0] != authorizationType {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(ErrAuthValWrongFormat))
			return
		}

		payload, err := tokenMaker.Verify(headerTokens[1])
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		username := payload.Username

		ctx.Set(authorizedUserId, username)

		ctx.Next()

	}
}

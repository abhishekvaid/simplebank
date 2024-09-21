package api

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	db "himavisoft.simple_bank/db/sqlc"
	"himavisoft.simple_bank/util"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	m.Run()
}

func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		TokenSecret: util.RandomString(32),
		TokenExpiry: time.Minute,
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
}

func setupToken(t *testing.T, server *Server, req *http.Request, username string, expiry time.Duration) {
	token, err := server.tokenMaker.Create(username, expiry)
	require.NoError(t, err)
	req.Header.Set(authorizationHeaderKey, fmt.Sprintf("%s %s", authorizationType, token))
}

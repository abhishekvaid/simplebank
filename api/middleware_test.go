package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"himavisoft.simple_bank/token"
	"himavisoft.simple_bank/util"
)

func setupAuthInReq(t *testing.T, maker token.Maker, expiryDuration time.Duration, request *http.Request, username, authorizationHeaderKey string) {

	token, err := maker.Create(username, expiryDuration)
	require.NoError(t, err)
	request.Header.Set(authorizationHeaderKey, fmt.Sprintf("%s %s", authorizationType, token))
}

func TestAuthMiddleware(t *testing.T) {

	testcases := []struct {
		name          string
		setupTest     func(t *testing.T, maker token.Maker, req *http.Request)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupTest: func(t *testing.T, maker token.Maker, req *http.Request) {
				setupAuthInReq(t, maker, time.Minute, req, util.RandomUsername(10), authorizationHeaderKey)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusOK)
			},
		},
		{
			name: "No Authorization",
			setupTest: func(t *testing.T, maker token.Maker, req *http.Request) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusUnauthorized)
			},
		},
		{
			name: "UnSupported Authorization",
			setupTest: func(t *testing.T, maker token.Maker, req *http.Request) {
				setupAuthInReq(t, maker, time.Minute, req, util.RandomUsername(10), "unsupported")
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusUnauthorized)
			},
		},
		{
			name: "Invalid Authorization Type",
			setupTest: func(t *testing.T, maker token.Maker, req *http.Request) {
				setupAuthInReq(t, maker, time.Minute, req, util.RandomUsername(10), "")
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusUnauthorized)
			},
		},
	}

	for _, testcase := range testcases {

		t.Run(testcase.name, func(t *testing.T) {

			server := newTestServer(t, nil)
			authURL := "/auth"
			server.router.GET(authURL, authMiddleware(server.tokenMaker), func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{})
			})

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, authURL, nil)

			require.NoError(t, err)

			testcase.setupTest(t, server.tokenMaker, request)

			server.router.ServeHTTP(recorder, request)

			testcase.checkResponse(t, recorder)
		})
	}

}

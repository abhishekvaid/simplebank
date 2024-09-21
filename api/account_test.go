package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockdb "himavisoft.simple_bank/db/mock"
	db "himavisoft.simple_bank/db/sqlc"
	"himavisoft.simple_bank/util"
)

func TestGetAccount(t *testing.T) {

	user := randomUser(t)
	account := randomAccount(user)

	testCases := []struct {
		Name           string
		ID             int64
		prepareRequest func(t *testing.T, server *Server, req *http.Request)
		buildStubs     func(mockStore *mockdb.MockStore)
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			Name: "200 Ok",
			ID:   account.ID,
			prepareRequest: func(t *testing.T, server *Server, req *http.Request) {
				setupToken(t, server, req, user.Username, time.Minute)
			},
			buildStubs: func(mockStore *mockdb.MockStore) {
				getAccountParams := db.GetAccountParams{
					ID:    account.ID,
					Owner: user.Username,
				}
				mockStore.EXPECT().GetAccount(gomock.Any(), getAccountParams).Times(1).Return(*account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				util.CompareResponseBodyJSON(t, recorder, account)
			},
		},
		{
			Name: "NotFound",
			ID:   account.ID * 10,
			prepareRequest: func(t *testing.T, server *Server, req *http.Request) {
				setupToken(t, server, req, user.Username, time.Minute)
			},
			buildStubs: func(mockStore *mockdb.MockStore) {
				getAccountParams := db.GetAccountParams{
					ID:    account.ID * 10,
					Owner: user.Username,
				}
				mockStore.EXPECT().GetAccount(gomock.Any(), getAccountParams).Times(1).Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, rec.Code)
			},
		},
		{
			Name: "BadRequest",
			ID:   0,
			prepareRequest: func(t *testing.T, server *Server, req *http.Request) {
				setupToken(t, server, req, user.Username, time.Minute)
			},
			buildStubs: func(mockStore *mockdb.MockStore) {
				getAccountParams := db.GetAccountParams{
					ID:    0,
					Owner: user.Username,
				}
				mockStore.EXPECT().GetAccount(gomock.Any(), getAccountParams).Times(0).Return(db.Account{}, nil)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			Name: "401 No Token Sent",
			ID:   0,
			prepareRequest: func(t *testing.T, server *Server, req *http.Request) {

			},
			buildStubs: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().GetAccount(gomock.Any(), 0).Times(0).Return(db.Account{}, nil)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, rec.Code)
				util.CompareResponseBodyJSON(t, rec, errorResponse(ErrAuthNoHeader))
			},
		},
		{
			Name: "500 InternalServerError",
			ID:   account.ID,
			prepareRequest: func(t *testing.T, server *Server, req *http.Request) {
				setupToken(t, server, req, user.Username, time.Minute)
			},
			buildStubs: func(mockStore *mockdb.MockStore) {
				getAccountParams := db.GetAccountParams{
					ID:    account.ID,
					Owner: user.Username,
				}
				mockStore.EXPECT().GetAccount(gomock.Any(), getAccountParams).Times(1).Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, rec.Code)
			},
		},
	}

	for _, testCase := range testCases {

		t.Run(testCase.Name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			testCase.buildStubs(store)

			url := fmt.Sprintf("/accounts/%d", testCase.ID)

			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			testCase.prepareRequest(t, server, request)

			server.router.ServeHTTP(recorder, request)

			testCase.checkResponse(t, recorder)

		})

	}
}

func requireBodyMatchAccount(t *testing.T, buf *bytes.Buffer, want db.Account) {
	t.Helper()
	bytes, err := io.ReadAll(buf)
	require.NoError(t, err)
	got := db.Account{}
	err = json.Unmarshal(bytes, &got)
	require.NoError(t, err)
	require.Equal(t, got, want)
}

func randomAccount(user *db.User) *db.Account {
	return &db.Account{
		ID:       int64(util.RandomInt(10, 10000)),
		Owner:    user.Username,
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	}
}

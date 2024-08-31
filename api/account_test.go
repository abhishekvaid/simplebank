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

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockdb "himavisoft.simple_bank/db/mock"
	db "himavisoft.simple_bank/db/sqlc"
	"himavisoft.simple_bank/util"
)

func TestGetAccount(t *testing.T) {

	account := randomAccount()

	testCases := []struct {
		Name          string
		ID            int64
		buildStubs    func(mockStore *mockdb.MockStore)
		checkResponse func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			Name: "StatusOk",
			ID:   account.ID,
			buildStubs: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().GetAccount(gomock.Any(), account.ID).Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, rec.Code)
				requireBodyMatchAccount(t, rec.Body, account)
			},
		},
		{
			Name: "NotFound",
			ID:   account.ID,
			buildStubs: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().GetAccount(gomock.Any(), account.ID).Times(1).Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, rec.Code)
			},
		},
		{
			Name: "BadRequest",
			ID:   0,
			buildStubs: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().GetAccount(gomock.Any(), 0).Times(0).Return(db.Account{}, nil)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			Name: "BadRequest",
			ID:   0,
			buildStubs: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().GetAccount(gomock.Any(), 0).Times(0).Return(db.Account{}, nil)
			},
			checkResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, rec.Code)
			},
		},
		{
			Name: "InternalServerError",
			ID:   account.ID,
			buildStubs: func(mockStore *mockdb.MockStore) {
				mockStore.EXPECT().GetAccount(gomock.Any(), account.ID).Times(1).Return(db.Account{}, sql.ErrConnDone)
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
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			testCase.buildStubs(store)

			url := fmt.Sprintf("/accounts/%d", testCase.ID)

			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

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

func randomAccount() db.Account {
	return db.Account{
		ID:       int64(util.RandomInt(10, 10000)),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
	}
}

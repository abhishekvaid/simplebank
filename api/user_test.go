package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	mockdb "himavisoft.simple_bank/db/mock"
	db "himavisoft.simple_bank/db/sqlc"
	"himavisoft.simple_bank/util"
)

type eqCreateUsersParamMatcher struct {
	createUserRequest *createUserRequest
}

func (matcher eqCreateUsersParamMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}
	if util.CheckPassword(matcher.createUserRequest.Password, arg.HashedPassword) != nil {
		return false
	}

	arg2 := db.CreateUserParams{
		Username:       matcher.createUserRequest.Username,
		HashedPassword: arg.HashedPassword,
		FullName:       matcher.createUserRequest.FullName,
		Email:          matcher.createUserRequest.Email,
	}
	fmt.Println(arg2)
	return reflect.DeepEqual(arg2, arg)
}

func (matcher eqCreateUsersParamMatcher) String() string {
	return fmt.Sprintf("checks if \"%s\" plain text password hashes to passed in password", matcher.createUserRequest.Password)
}

func TestCreateUser(t *testing.T) {

	constraintErr := &pq.Error{
		Code: "23505", // PostgreSQL code for unique violation
	}

	randomErr := fmt.Errorf("Some random error")

	aCreateUserRequest := randomCreateUserRequest()
	anEmptyUser := db.User{}

	matcher := eqCreateUsersParamMatcher{
		createUserRequest: aCreateUserRequest,
	}

	testcases := []struct {
		name            string
		setupStub       func(*mockdb.MockStore)
		checkConditions func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "202",
			setupStub: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					CreateUser(gomock.Any(), matcher).
					Times(1).
					Return(anEmptyUser, nil)
			},
			checkConditions: func(t *testing.T, respRecorder *httptest.ResponseRecorder) {
				require.Equal(t, respRecorder.Code, http.StatusOK)
				util.MatchResponseBodyWith(t, respRecorder.Body, db.User{})
			},
		},
		{
			name: "400",
			setupStub: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					CreateUser(gomock.Any(), matcher).
					Times(1).
					Return(anEmptyUser, constraintErr)
			},
			checkConditions: func(t *testing.T, respRecorder *httptest.ResponseRecorder) {
				require.Equal(t, respRecorder.Code, http.StatusForbidden)
				util.MatchResponseBodyWith(t, respRecorder.Body, errorResponse(constraintErr))
			},
		},
		{
			name: "500",
			setupStub: func(store *mockdb.MockStore) {
				store.
					EXPECT().
					CreateUser(gomock.Any(), matcher).
					Times(1).
					Return(anEmptyUser, randomErr)
			},
			checkConditions: func(t *testing.T, respRecorder *httptest.ResponseRecorder) {
				require.Equal(t, respRecorder.Code, http.StatusInternalServerError)
				util.MatchResponseBodyWith(t, respRecorder.Body, errorResponse(randomErr))
			},
		},
	}

	for _, testcase := range testcases {

		t.Run(testcase.name, func(t *testing.T) {

			t.Helper()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStore := mockdb.NewMockStore(ctrl)
			server := newTestServer(t, mockStore)

			testcase.setupStub(mockStore)

			byteArr, err := json.Marshal(aCreateUserRequest)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewReader(byteArr))
			require.NoError(t, err)
			responseRecorder := httptest.NewRecorder()

			server.router.ServeHTTP(responseRecorder, req)

			testcase.checkConditions(t, responseRecorder)

		})
	}

}

func TestGetUser(t *testing.T) {

	user := randomUser(t)

	testcases := []struct {
		name           string
		prepareRequest func(t *testing.T, server *Server, req *http.Request)
		setupStubs     func(t *testing.T, m *mockdb.MockStore)
		checkResponse  func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "401 No Token Sent",
			prepareRequest: func(t *testing.T, server *Server, req *http.Request) {
				// don't setup any token
			},
			setupStubs: func(t *testing.T, m *mockdb.MockStore) {
				m.EXPECT().GetUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusUnauthorized)
				util.CompareResponseBodyJSON(t, recorder, errorResponse(ErrAuthNoHeader))
			},
		},
		{
			name: "200 Ok",
			prepareRequest: func(t *testing.T, server *Server, req *http.Request) {
				setupToken(t, server, req, user.Username, time.Minute)
			},
			setupStubs: func(t *testing.T, m *mockdb.MockStore) {
				m.EXPECT().GetUser(gomock.Any(), user.Username).Times(1).Return(*user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				util.CompareResponseBodyJSON(t, recorder, user)
			},
		},
	}

	for _, testcase := range testcases {

		t.Run(testcase.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStore := mockdb.NewMockStore(ctrl)
			server := newTestServer(t, mockStore)

			req, err := http.NewRequest(http.MethodGet, "/users", nil)
			require.NoError(t, err)

			responseRecorder := httptest.NewRecorder()

			testcase.prepareRequest(t, server, req)
			testcase.setupStubs(t, mockStore)

			server.router.ServeHTTP(responseRecorder, req)

			testcase.checkResponse(t, responseRecorder)

		})
	}

}

func randomUser(t *testing.T) *db.User {
	hashedPassword, err := util.HashPassword("secret")
	require.NoError(t, err)
	return &db.User{
		Username:          util.RandomUsername(10),
		HashedPassword:    hashedPassword,
		FullName:          util.RandomString(10),
		Email:             util.RandomEmail(),
		PasswordChangedAt: time.Now(),
		CreatedAt:         time.Now(),
	}
}

func randomCreateUserRequest() *createUserRequest {
	return &createUserRequest{
		Username: util.RandomString(10) + "__username",
		Password: util.RandomString(10) + "__password",
		FullName: util.RandomString(5) + "__fullName",
		Email:    util.RandomEmail(),
	}
}

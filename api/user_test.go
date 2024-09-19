package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

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
				util.RequireBodyMatch(t, respRecorder.Body, db.User{})
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
				util.RequireBodyMatch(t, respRecorder.Body, errorResponse(constraintErr))
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
				util.RequireBodyMatch(t, respRecorder.Body, errorResponse(randomErr))
			},
		},
	}

	for _, testcase := range testcases {

		t.Run(testcase.name, func(t *testing.T) {

			t.Helper()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStore := mockdb.NewMockStore(ctrl)
			server := NewServer(mockStore)

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

func randomCreateUserRequest() *createUserRequest {
	return &createUserRequest{
		Username: util.RandomString(10) + "__username",
		Password: util.RandomString(10) + "__password",
		FullName: util.RandomString(5) + "__fullName",
		Email:    util.RandomEmail(),
	}
}

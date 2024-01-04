package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	mockdb "github.com/broemp/red_card/db/mock"
	db "github.com/broemp/red_card/db/sqlc"
	"github.com/broemp/red_card/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := util.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func TestCreateUserAPI(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username": user.Username,
				"password": password,
				"name":     user.Name,
			},
			buildStubs: func(store *mockdb.MockStore) {
				// arg := db.CreateUserParams{
				// 	Username: user.Username,
				// 	Name:     user.Name,
				// }
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"username": user.Username,
				"password": password,
				"name":     user.Name,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidUsername",
			body: gin.H{
				"username": "invalid-user#1",
				"password": password,
				"name":     user.Name,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "TooShortPassword",
			body: gin.H{
				"username": user.Username,
				"password": "123",
				"name":     user.Name,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)
			fmt.Println(string(data))

			url := "/users/register"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

// TODO: Fix test
//	func TestLoginUserAPI(t *testing.T) {
//		user, password := randomUser(t)
//
//		testCases := []struct {
//			name          string
//			body          gin.H
//			buildStubs    func(store *mockdb.MockStore)
//			checkResponse func(recorder *httptest.ResponseRecorder)
//		}{
//			{
//				name: "OK",
//				body: gin.H{
//					"username": user.Username,
//					"password": password,
//				},
//				buildStubs: func(store *mockdb.MockStore) {
//					store.EXPECT().
//						GetUserAuth(gomock.Any(), gomock.Any()).
//						Times(1).
//						Return(user, nil)
//					store.EXPECT().CreateSession(gomock.Any(), gomock.Any()).
//						Times(1)
//				},
//				checkResponse: func(recorder *httptest.ResponseRecorder) {
//					require.Equal(t, http.StatusOK, recorder.Code)
//				},
//			},
//			{
//				name: "UserNotFound",
//				body: gin.H{
//					"username": "NotFound",
//					"password": password,
//				},
//				buildStubs: func(store *mockdb.MockStore) {
//					store.EXPECT().
//						GetUserAuth(gomock.Any(), gomock.Any()).
//						Times(1).
//						Return(db.User{}, sql.ErrNoRows)
//				},
//				checkResponse: func(recorder *httptest.ResponseRecorder) {
//					require.Equal(t, http.StatusNotFound, recorder.Code)
//				},
//			},
//		}
//		for i := range testCases {
//			tc := testCases[i]
//
//			t.Run(tc.name, func(t *testing.T) {
//				ctrl := gomock.NewController(t)
//				defer ctrl.Finish()
//
//				store := mockdb.NewMockStore(ctrl)
//				tc.buildStubs(store)
//
//				server := newTestServer(t, store)
//				recorder := httptest.NewRecorder()
//
//				data, err := json.Marshal(tc.body)
//				require.NoError(t, err)
//
//				url := "/users/login"
//				request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
//				require.NoError(t, err)
//
//				server.router.ServeHTTP(recorder, request)
//				tc.checkResponse(recorder)
//			})
//		}
//	}

func randomUser(t *testing.T) (user db.User, password string) {
	password = util.RandomString(8)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		Username:       util.RandomUsername(),
		HashedPassword: hashedPassword,
		Name:           util.RandomUsername(),
	}
	return
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)

	require.NoError(t, err)
	require.Equal(t, user.Username, gotUser.Username)
	require.Equal(t, user.Name, gotUser.Name)
	require.Empty(t, gotUser.HashedPassword)
}

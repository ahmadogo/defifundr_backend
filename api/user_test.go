package api

// import (
// 	"bytes"
// 	"database/sql"
// 	"encoding/json"
// 	"fmt"

// 	"net/http"
// 	"net/http/httptest"
// 	"reflect"
// 	"testing"

// 	mockdb "github.com/demola234/defiraise/db/mock"
// 	db "github.com/demola234/defiraise/db/sqlc"
// 	"github.com/demola234/defiraise/utils"
// 	"github.com/gin-gonic/gin"
// 	"github.com/golang/mock/gomock"
// 	"github.com/stretchr/testify/require"
// )

// type eqCreateUserParamsMatcher struct {
// 	arg      db.CreateUserParams
// 	password string
// }

// func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
// 	arg, ok := x.(db.CreateUserParams)
// 	if !ok {
// 		return false
// 	}

// 	err := utils.CheckPassword(e.password, arg.HashedPassword)
// 	if err != nil {
// 		return false
// 	}

// 	e.arg.HashedPassword = arg.HashedPassword
// 	return reflect.DeepEqual(e.arg, arg)
// }

// func (e eqCreateUserParamsMatcher) String() string {
// 	return fmt.Sprintf("%v %v", e.arg, e.password)
// }

// func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
// 	return eqCreateUserParamsMatcher{arg, password}
// }

// func randomUser(t *testing.T) (user db.Users, password string) {
// 	password = utils.RandomString(6)
// 	Username := utils.RandomOwner()

// 	hashedPassword, err := utils.HashPassword(password)
// 	require.NoError(t, err)

// 	user = db.Users{
// 		Username:       Username,
// 		Avatar:         utils.RandomString(6),
// 		Email:          utils.RandomEmail(),
// 		HashedPassword: hashedPassword,
// 	}
// 	return user, password
// }

// func TestCreateUser(t *testing.T) {
// 	users, password := randomUser(t)

// 	testCases := []struct {
// 		name          string
// 		body          gin.H
// 		buildStubs    func(store *mockdb.MockStore)
// 		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
// 	}{
// 		{
// 			name: "OK",
// 			body: gin.H{
// 				"username": users.Username,
// 				"email":    users.Email,
// 				"password": password,
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				arg := db.CreateUserParams{
// 					Username: users.Username,
// 					Email:    users.Email,
// 				}

// 				store.EXPECT().
// 					CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).
// 					Times(1).
// 					Return(users, nil)
// 			},
// 			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recoder.Code)
// 			},
// 		},
// 	}

// 	for i := range testCases {
// 		tc := testCases[i]
// 		t.Run(tc.name, func(t *testing.T) {
// 			tc := testCases[i]
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			store := mockdb.NewMockStore(ctrl)
// 			tc.buildStubs(store)

// 			server := newTestServer(t, store)
// 			recorder := httptest.NewRecorder()

// 			data, err := json.Marshal(tc.body)
// 			require.NoError(t, err)

// 			url := "/users"
// 			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
// 			require.NoError(t, err)

// 			server.router.ServeHTTP(recorder, request)
// 			tc.checkResponse(t, recorder)
// 		})
// 	}
// }

// // func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.Users) {
// // 	data, err := os.ReadFile(body.String())
// // 	require.NoError(t, err)

// // 	var gotUser db.Users
// // 	err = json.Unmarshal(data, &gotUser)

// // 	require.NoError(t, err)
// // 	require.Equal(t, user.Username, gotUser.Username)
// // 	require.Equal(t, user.Email, gotUser.Email)
// // 	require.Empty(t, gotUser.HashedPassword)
// // }

// func TestLoginUser(t *testing.T) {
// 	users, password := randomUser(t)

// 	testCases := []struct {
// 		name          string
// 		buildStubs    func(store *mockdb.MockStore)
// 		body          gin.H
// 		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
// 	}{
// 		{
// 			name: "OK",
// 			body: gin.H{
// 				"username": users.Username,
// 				"password": password,
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {

// 				store.EXPECT().
// 					GetUser(gomock.Any(), gomock.Eq(users.Username)).
// 					Times(1).
// 					Return(users, nil)

// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "NoUser",
// 			body: gin.H{
// 				"username": users.Username,
// 				"password": password,
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {

// 				store.EXPECT().
// 					GetUser(gomock.Any(), gomock.Eq(users.Username)).
// 					Times(1).
// 					Return(db.Users{}, sql.ErrNoRows)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusNotFound, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "InvalidPassword",
// 			body: gin.H{
// 				"username": users.Username,
// 				"password": "invalid_password",
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {

// 				store.EXPECT().
// 					GetUser(gomock.Any(), gomock.Eq(users.Username)).
// 					Times(1).
// 					Return(users, nil)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusUnauthorized, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "EmailLogin",
// 			body: gin.H{
// 				"username": users.Email,
// 				"password": password,
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {

// 				store.EXPECT().
// 					GetUser(gomock.Any(), gomock.Eq(users.Email)).
// 					Times(1).
// 					Return(users, nil)
// 				store.EXPECT().
// 					CreateSession(gomock.Any(), gomock.Any()).
// 					Times(1)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 			},
// 		},
// 	}
// 	for i := range testCases {
// 		tc := testCases[i]
// 		t.Run(tc.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			store := mockdb.NewMockStore(ctrl)
// 			tc.buildStubs(store)

// 			server := newTestServer(t, store)
// 			recorder := httptest.NewRecorder()

// 			data, err := json.Marshal(tc.body)
// 			require.NoError(t, err)

// 			url := "/users/login"
// 			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
// 			require.NoError(t, err)

// 			server.router.ServeHTTP(recorder, request)
// 			tc.checkResponse(t, recorder)

// 		})
// 	}

// }

// func TestVerifyUser(t *testing.T) {
// 	users, password := randomUser(t)

// 	testCases := []struct {
// 		name          string
// 		buildStubs    func(store *mockdb.MockStore)
// 		body          gin.H
// 		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
// 	}{
// 		{
// 			name: "OK",
// 			body: gin.H{
// 				"username": users.Username,
// 				"password": password,
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {

// 				store.EXPECT().
// 					GetUser(gomock.Any(), gomock.Eq(users.Username)).
// 					Times(1).
// 					Return(users, nil)

// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "NoUser",
// 			body: gin.H{
// 				"username": users.Username,
// 				"password": password,
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {

// 				store.EXPECT().
// 					GetUser(gomock.Any(), gomock.Eq(users.Username)).
// 					Times(1).
// 					Return(db.Users{}, sql.ErrNoRows)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusNotFound, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "InvalidPassword",
// 			body: gin.H{
// 				"username": users.Username,
// 				"password": "invalid_password",
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {

// 				store.EXPECT().
// 					GetUser(gomock.Any(), gomock.Eq(users.Username)).
// 					Times(1).
// 					Return(users, nil)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusUnauthorized, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "EmailLogin",
// 			body: gin.H{
// 				"username": users.Email,
// 				"password": password,
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {

// 				store.EXPECT().
// 					GetUser(gomock.Any(), gomock.Eq(users.Email)).
// 					Times(1).
// 					Return(users, nil)
// 				store.EXPECT().
// 					CreateSession(gomock.Any(), gomock.Any()).
// 					Times(1)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 			},
// 		},
// 	}
// 	for i := range testCases {
// 		tc := testCases[i]
// 		t.Run(tc.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			store := mockdb.NewMockStore(ctrl)
// 			tc.buildStubs(store)

// 			server := newTestServer(t, store)
// 			recorder := httptest.NewRecorder()

// 			data, err := json.Marshal(tc.body)
// 			require.NoError(t, err)

// 			url := "/users/verify"
// 			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
// 			require.NoError(t, err)

// 			server.router.ServeHTTP(recorder, request)
// 			tc.checkResponse(t, recorder)

// 		})
// 	}
// }

// func TestResendVerificationCode(t *testing.T) {
// 	users, password := randomUser(t)

// 	testCases := []struct {
// 		name          string
// 		buildStubs    func(store *mockdb.MockStore)
// 		body          gin.H
// 		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
// 	}{
// 		{
// 			name: "OK",
// 			body: gin.H{
// 				"username": users.Username,
// 				"password": password,
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {

// 				store.EXPECT().
// 					GetUser(gomock.Any(), gomock.Eq(users.Username)).
// 					Times(1).
// 					Return(users, nil)

// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "NoUser",
// 			body: gin.H{
// 				"username": users.Username,
// 				"password": password,
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {

// 				store.EXPECT().
// 					GetUser(gomock.Any(), gomock.Eq(users.Username)).
// 					Times(1).
// 					Return(db.Users{}, sql.ErrNoRows)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusNotFound, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "InvalidPassword",
// 			body: gin.H{
// 				"username": users.Username,
// 				"password": "invalid_password",
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {

// 				store.EXPECT().
// 					GetUser(gomock.Any(), gomock.Eq(users.Username)).
// 					Times(1).
// 					Return(users, nil)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusUnauthorized, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "EmailLogin",
// 			body: gin.H{
// 				"username": users.Email,
// 				"password": password,
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {

// 				store.EXPECT().
// 					GetUser(gomock.Any(), gomock.Eq(users.Email)).
// 					Times(1).
// 					Return(users, nil)
// 				store.EXPECT().
// 					CreateSession(gomock.Any(), gomock.Any()).
// 					Times(1)

// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 			},
// 		},
// 	}
// 	for i := range testCases {
// 		tc := testCases[i]
// 		t.Run(tc.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			store := mockdb.NewMockStore(ctrl)
// 			tc.buildStubs(store)

// 			server := newTestServer(t, store)
// 			recorder := httptest.NewRecorder()

// 			data, err := json.Marshal(tc.body)
// 			require.NoError(t, err)

// 			url := "/users/verify/resend"
// 			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
// 			require.NoError(t, err)

// 			server.router.ServeHTTP(recorder, request)
// 			tc.checkResponse(t, recorder)

// 		})
// 	}
// }

// func TestResetPassword(t *testing.T) {
// 	users, password := randomUser(t)

// 	testCases := []struct {
// 		name          string
// 		buildStubs    func(store *mockdb.MockStore)
// 		body          gin.H
// 		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
// 	}{
// 		{
// 			name: "OK",
// 			body: gin.H{
// 				"username": users.Username,
// 				"password": password,
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {

// 				store.EXPECT().
// 					GetUser(gomock.Any(), gomock.Eq(users.Username)).
// 					Times(1).
// 					Return(users, nil)

// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "NoUser",
// 			body: gin.H{
// 				"username": users.Username,
// 				"password": password,
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {

// 				store.EXPECT().
// 					GetUser(gomock.Any(), gomock.Eq(users.Username)).
// 					Times(1).
// 					Return(db.Users{}, sql.ErrNoRows)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusNotFound, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "InvalidPassword",
// 			body: gin.H{
// 				"username": users.Username,
// 				"password": "invalid_password",
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {

// 				store.EXPECT().
// 					GetUser(gomock.Any(), gomock.Eq(users.Username)).
// 					Times(1).
// 					Return(users, nil)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusUnauthorized, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "EmailLogin",
// 			body: gin.H{
// 				"username": users.Email,
// 				"password": password,
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {

// 				store.EXPECT().
// 					GetUser(gomock.Any(), gomock.Eq(users.Email)).
// 					Times(1).
// 					Return(users, nil)
// 				store.EXPECT().
// 					CreateSession(gomock.Any(), gomock.Any()).
// 					Times(1)

// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 			},
// 		},
// 	}
// 	for i := range testCases {
// 		tc := testCases[i]
// 		t.Run(tc.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			store := mockdb.NewMockStore(ctrl)
// 			tc.buildStubs(store)

// 			server := newTestServer(t, store)
// 			recorder := httptest.NewRecorder()

// 			data, err := json.Marshal(tc.body)
// 			require.NoError(t, err)

// 			url := "/users/reset-password"
// 			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
// 			require.NoError(t, err)

// 			server.router.ServeHTTP(recorder, request)
// 			tc.checkResponse(t, recorder)

// 		})
// 	}
// }

// func TestCheckUsernameExists(t *testing.T) {
// 	users, _ := randomUser(t)

// 	testCases := []struct {
// 		name          string
// 		buildStubs    func(store *mockdb.MockStore)
// 		body          gin.H
// 		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
// 	}{
// 		{
// 			name: "OK",
// 			body: gin.H{
// 				"username": users.Username,
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {

// 				store.EXPECT().
// 					GetUser(gomock.Any(), gomock.Eq(users.Username)).
// 					Times(1).
// 					Return(users, nil)

// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 				require.Equal(t, "{\"exists\":true}", recorder.Body.String())
// 			},
// 		},
// 		{
// 			name: "NoUser",
// 			body: gin.H{
// 				"username": users.Username,
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {

// 				store.EXPECT().
// 					GetUser(gomock.Any(), gomock.Eq(users.Username)).
// 					Times(1).
// 					Return(db.Users{}, sql.ErrNoRows)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 				require.Equal(t, "{\"exists\":false}", recorder.Body.String())
// 			},
// 		},
// 		{
// 			name: "InvalidPassword",
// 			body: gin.H{
// 				"username": users.Username,
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {

// 				store.EXPECT().
// 					GetUser(gomock.Any(), gomock.Eq(users.Username)).
// 					Times(1).
// 					Return(users, nil)
// 			},
// 			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 				require.Equal(t, "{\"exists\":true}", recorder.Body.String())
// 			},
// 		},
// 	}
// 	for i := range testCases {
// 		tc := testCases[i]
// 		t.Run(tc.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()

// 			store := mockdb.NewMockStore(ctrl)
// 			tc.buildStubs(store)

// 			server := newTestServer(t, store)
// 			recorder := httptest.NewRecorder()

// 			data, err := json.Marshal(tc.body)
// 			require.NoError(t, err)

// 			url := "/users/checkUsername/:username"
// 			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
// 			require.NoError(t, err)

// 			server.router.ServeHTTP(recorder, request)
// 			tc.checkResponse(t, recorder)

// 		})
// 	}
// }

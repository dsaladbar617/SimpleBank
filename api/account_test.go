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

	mockdb "github.com/dsaladbar617/SimpleBank/db/mock"
	db "github.com/dsaladbar617/SimpleBank/db/sqlc"
	"github.com/dsaladbar617/SimpleBank/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetAccountAPI(t *testing.T) {
	account := randomAccount()

	testCases := map[string]struct{
		accountID int64
		buildStubs func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{

			"OK": {
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
			GetAccount(gomock.Any(), gomock.Eq(account.ID)).
			Times(1).
			Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			} ,
		},
			"NotFound": {
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},

			"InternalError": {
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
			"InvalidID": {
			accountID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		}

	for name, test := range testCases{

		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			// defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			test.buildStubs(store)


			server :=newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", test.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			test.checkResponse(t, recorder)
		})
	}
}

func TestCreateAccountAPI(t *testing.T) {
	account := randomAccount()
	testCases := map[string]struct{
		body gin.H
		buildStubs func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	} {

			"OK": {
				body: gin.H{
					"owner": account.Owner,
					"currency": account.Currency,
				},
				buildStubs: func(store *mockdb.MockStore) {
					arg := db.CreateAccountParams{
						Owner: account.Owner,
						Currency: account.Currency,
						Balance: 0,
					}

					store.EXPECT().
							CreateAccount(gomock.Any(), gomock.Eq(arg)).
							Times(1).
							Return(account, nil)
				},
				checkResponse: func(recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusOK, recorder.Code)
					requireBodyMatchAccount(t, recorder.Body, account)
				},

			},
			"InternalError" : {
				body: gin.H{
					"owner":    account.Owner,
					"currency": account.Currency,
				},
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						CreateAccount(gomock.Any(), gomock.Any()).
						Times(1).
						Return(db.Account{}, sql.ErrConnDone)
				},
				checkResponse: func(recorder *httptest.ResponseRecorder){
					require.Equal(t, http.StatusInternalServerError, recorder.Code)
				},
			},
			"InvalidCurrency": {
				body: gin.H{
					"owner":   	account.Owner,
					"currency": "invalid",
				},
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						CreateAccount(gomock.Any(), gomock.Any()).
						Times(0)
				},
				checkResponse: func(recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusBadRequest, recorder.Code)
				},
			},

			"InvalidOwner": {
				body: gin.H{
					"owner":    "",
					"currency": account.Currency,
				},
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						CreateAccount(gomock.Any(), gomock.Any()).
						Times(0)
				},
				checkResponse: func(recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusBadRequest, recorder.Code)
				},
			},
		}


// 	testCases := []struct{
// 		name string
// 		body gin.H
// 		buildStubs func(store *mockdb.MockStore)
// 		checkResponse func(recorder *httptest.ResponseRecorder)
// 	} {
// 		{
// 			name: "OK",
// 			body: gin.H{
// 				"owner": account.Owner,
// 				"currency": account.Currency,
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				arg := db.CreateAccountParams{
// 					Owner: account.Owner,
// 					Currency: account.Currency,
// 					Balance: 0,
// 				}

// 				store.EXPECT().
// 						CreateAccount(gomock.Any(), gomock.Eq(arg)).
// 						Times(1).
// 						Return(account, nil)
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 				requireBodyMatchAccount(t, recorder.Body, account)
// 			},
// 		},
// 		{
// 			name: "InternalError",
// 			body: gin.H{
// 				"owner":    account.Owner,
// 				"currency": account.Currency,
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					CreateAccount(gomock.Any(), gomock.Any()).
// 					Times(1).
// 					Return(db.Account{}, sql.ErrConnDone)
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder){
// 				require.Equal(t, http.StatusInternalServerError, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "InvalidCurrency",
// 			body: gin.H{
// 				"owner":   	account.Owner,
// 				"currency": "invalid",
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					CreateAccount(gomock.Any(), gomock.Any()).
// 					Times(0)
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusBadRequest, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "InvalidOwner",
// 			body: gin.H{
// 				"owner":    "",
// 				"currency": account.Currency,
// 			},
// 			buildStubs: func(store *mockdb.MockStore) {
// 				store.EXPECT().
// 					CreateAccount(gomock.Any(), gomock.Any()).
// 					Times(0)
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusBadRequest, recorder.Code)
// 		},
// 	},
// }
	for name, test := range testCases {

		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			test.buildStubs(store)

			server :=newTestServer(t, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(test.body)
			require.NoError(t, err)

			url := "/accounts"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			test.checkResponse(recorder)
		})
	}
	// for i := range testCases {
	// 	tc := testCases[i]

	// 	t.Run(tc.name, func(t *testing.T) {
	// 		ctrl := gomock.NewController(t)
	// 		defer ctrl.Finish()

	// 		store := mockdb.NewMockStore(ctrl)
	// 		tc.buildStubs(store)

	// 		server := NewServer(store)
	// 		recorder := httptest.NewRecorder()

	// 		data, err := json.Marshal(tc.body)
	// 		require.NoError(t, err)

	// 		url := "/accounts"
	// 		request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	// 		require.NoError(t, err)

	// 		server.router.ServeHTTP(recorder, request)
	// 		tc.checkResponse(recorder)
	// 	})
	// }
}

func randomAccount() db.Account {
	return db.Account{
		ID: util.RandomInt(1,1000),
		Owner: util.RandomOwner(),
		Balance: util.RandomMoney(),
		Currency: util.RandomCurrency(),
		CreatedAt: time.Now().Truncate(time.Second),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}
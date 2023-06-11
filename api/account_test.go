package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/Richd0tcom/SafeX-Pay/db/mock"
	db "github.com/Richd0tcom/SafeX-Pay/db/sqlc"
	"github.com/Richd0tcom/SafeX-Pay/utils"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type CaseGetAccount struct {
	name  string
	accountID interface{}
	buildStubs func(store *mockdb.MockStore)
	checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
}
func TestGetAccountAPI(t *testing.T){
	account := randomAccount()
	

	testCases := []CaseGetAccount{
		{
			name: "OK",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account,nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t,http.StatusOK, recorder.Code)//make sure the status code is th same as the main controlers

				resBody, err:= io.ReadAll(recorder.Body) 
				require.NoError(t, err)
			
				var tempAccount db.Account
			
				json.Unmarshal(resBody, &tempAccount)
				
				require.Equal(t, account, tempAccount)
			},
		},
		{
			name: "Not found",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{},sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t,http.StatusNotFound, recorder.Code)//make sure the status code is th same as the main controlers
			},
		},
		{
			name: "Bad Request(invalid id)",
			accountID: "yiuyouiou",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t,http.StatusBadRequest, recorder.Code)//make sure the status code is th same as the main controlers
			},
		},
		{
			name: "Internal Server Error",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{},sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t,http.StatusInternalServerError, recorder.Code)//make sure the status code is th same as the main controlers
			},
		},

	}

	for _, tc:= range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl:= gomock.NewController(t)
		// defer ctrl.Finish() //this is not needed in latest version of gomock. Register a Cleanup function instead
	
		store:= mockdb.NewMockStore(ctrl)
		//build stubs
		tc.buildStubs(store)
	
	
		//create a mock http seerver
		server:= NewServer(store)
	
		//use the http recorder to record the responses of the mock server
		recorder:= httptest.NewRecorder()
	
		url:= fmt.Sprintf("/account/%v", tc.accountID)
		request, err:= http.NewRequest(http.MethodGet, url, nil)
		require.NoError(t, err)
	
		server.serverRouter.ServeHTTP(recorder, request)

		tc.checkResponse(t, recorder)
		})
	}



	
}

//create a random account to make tests independent
func randomAccount() db.Account{
	//instantiates a new account object(struct)

	return db.Account{
		ID: uuid.New(),
		Owner: utils.RandomOwner(),
		Balance: utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}
}
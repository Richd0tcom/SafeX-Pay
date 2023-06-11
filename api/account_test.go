package api

import (
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

func TestGetAccount(t *testing.T){
	account := randomAccount()

	ctrl:= gomock.NewController(t)
	defer ctrl.Finish() //this is not needed in latest version of gomock. Register a Cleanup function instead

	fmt.Println(account.ID)
	
	store:= mockdb.NewMockStore(ctrl)
	//build stubs
	store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account,nil)


	//create a mock http seerver
	server:= NewServer(store)

	//use the http recorder to record the responses of the mock server
	recorder:= httptest.NewRecorder()

	url:= fmt.Sprintf("/account/%v", account.ID)
	request, err:= http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	server.serverRouter.ServeHTTP(recorder, request)

	require.Equal(t,http.StatusOK, recorder.Code)//make sure the status code is th same as the main controlers

	resBody, err:= io.ReadAll(recorder.Body) 
	require.NoError(t, err)

	var tempAccount db.Account

	json.Unmarshal(resBody, &tempAccount)
	
	require.Equal(t, account, tempAccount)
	
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
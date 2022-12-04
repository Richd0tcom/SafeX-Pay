package db

import (
	"context"
	"testing"
	"time"
	"database/sql"
	"github.com/Richd0tcom/SafeX-Pay/utils"
	"github.com/stretchr/testify/require"
)

//create a random account to make tests independent
func createRandomAccount(t *testing.T) Account{
	//instantiates a new account object(struct)
	args := CreateAccountParams{
		Owner: utils.RandomOwner(),
		Balance: utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}

	//creates a new account
	account, err := testQueries.CreateAccount(context.Background(), args)

	//checks that the CreateAccount function returns NO error. will fail the test if error exists
	require.NoError(t, err)
	//checks that the new account created should not be empty.
	require.NotEmpty(t,account)

	//checks that the values of the account object are of the sam type.
	require.Equal(t, args.Owner, account.Owner)
	require.Equal(t, args.Balance, account.Balance)
	require.Equal(t, args.Currency, account.Currency)

	//checks that the values of the Timestamp is not its zero value i.e it has been created
	require.NotZero(t,account.ID)
	require.NotZero(t,account.CreatedAt)

	return account
}
func TestCreateAccounts(t *testing.T){
	createRandomAccount(t)
}
func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)

	//check that there are no errors
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	//check that the returned data is valid
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}
func TestListAccounts(t *testing.T){
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}

}
func TestUpdateAccount(t *testing.T){
	account1 := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: utils.RandomMoney(),
	}

	account2, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}
func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	//try to retrieve deleted account. should throw an error
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}
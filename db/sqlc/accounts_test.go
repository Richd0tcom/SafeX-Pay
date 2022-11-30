package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/Richd0tcom/SafeX-Pay/utils"
)
 
func TestCreateAccounts(t *testing.T){
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

}
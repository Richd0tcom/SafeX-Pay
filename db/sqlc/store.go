package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// provides all functions to execute queries and transactions.
// we will use this interface to mock our tests by providing all the queries (method set) used by the SQLStore struct
type Store interface {
	Querier
	TransferTx(ctx context.Context, args CreateTransferParams) (TransferTxResults, error)
}


// provides all functions to execute queries and transactions.
//
// we will use composition instead of inheritance by mbedding the *Queries struct in the Store.
//
// SQLStore implements the Store interface.
type SQLStore struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new Store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// Takes a context and a callback function as input, starts a new database transaction, creat a new Queries object and with that transaction and call the callback function with the created Queries object and finally commit or rollback the transaction based on the error returned by the callback function.
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	querysObject := New(tx)
	txErr := fn(querysObject)
	if txErr != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			//this means that the rollback failed
			return fmt.Errorf("transaction Error: %v, Rollback Error: %v", txErr, rbErr)
		}
		return txErr //IF the rollback is successful, return the originall transaction error
	}
	//If all operations are successful, commit the transaction
	return tx.Commit()
}

// TransferTx perform a money transfer from one account to another account
func (store *SQLStore) TransferTx(ctx context.Context, args CreateTransferParams) (TransferTxResults, error) {
	var result TransferTxResults

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: args.FromAccountID,
			ToAccountID:   args.ToAccountID,
			Amount:        args.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.FromAccountID,
			Amount:    -args.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.ToAccountID,
			Amount:    args.Amount,
		})
		if err != nil {
			return err
		}

		//TODO: UPDATE account balance info (changed to use UUID)
		if args.FromAccountID.Time() < args.ToAccountID.Time(){
			result.FromAccount,result.ToAccount,err= addMoney(ctx, q, args.FromAccountID, -args.Amount, args.ToAccountID, args.Amount)
			if err != nil {
				return err
			}
		} else {
			result.ToAccount,result.FromAccount,err =addMoney(ctx, q, args.ToAccountID, args.Amount, args.FromAccountID, -args.Amount)
			if err != nil {
				return err
			}
		}
		
		//faulty implementation was fixed in 

		return nil
	})

	return result, err
}

type TransferTxParams struct {
	FromAccountID uuid.UUID `json:"from_account_id"`
	ToAccountID   uuid.UUID `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResults struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func addMoney(
	ctx context.Context,
	q  *Queries,
	accountID1 uuid.UUID,
	amount1 int64,
	accountID2 uuid.UUID,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err=q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return 
	}
	account2, err=q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	return
}
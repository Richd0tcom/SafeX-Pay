package db

import (
	"context"
	"database/sql"
	"fmt"
)

//provides all functions to execute queries and transactions.
//we will use composition instead of inheritance by mbedding the *Queries struct in the Store
type Store struct {
	*Queries
	db *sql.DB
}

//NewStore creates a new Store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
		Queries: New(db),
	}
}
var txKey = struct{}{}
//Takes a context and a callback function as input, starts a new database transaction, creat a new Queries object and with that transaction and call the callback function with the created Queries object and finally commit or rollback the transaction based on the error returned by the callback function.
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error{
	tx, err :=store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	querysObject := New(tx)
	txErr := fn(querysObject)
	if txErr != nil {
		rbErr:= tx.Rollback()
		if rbErr != nil {
			//this means that the rollback failed 
			return fmt.Errorf("transaction Error: %v, Rollback Error: %v", txErr, rbErr)
		}
		return txErr //IF the rollback is successful, return the originall transaction error
	}
	//If all operations are successful, commit the transaction
	return tx.Commit()
}

//TransferTx perform a money transfer from one account to another account
func (store Store) TransferTx(ctx context.Context, args TransferTxParams) (TransferTxResults, error) {
	var result TransferTxResults

	err:= store.execTx(ctx, func(q *Queries) error {
		var err error

		
		result.Transfer, err = q.CreateTransfer(ctx,CreateTransferParams{
			FromAccountID: args.FromAccountID,
			ToAccountID: args.ToAccountID,
			Amount: args.Amount,
		})
		if err != nil {
			return err
		}

		
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.FromAccountID,
			Amount: -args.Amount,
		})
		if err != nil {
			return err
		}

		
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.ToAccountID,
			Amount: args.Amount,
		})
		if err != nil {
			return err
		}

		//TODO: UPDATE account balance info
		
		acc1,err:= q.GetAccountForUpdate(ctx, args.FromAccountID)
		if err != nil {
			return err
		}
		
		result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID: args.FromAccountID,
			Balance: acc1.Balance - args.Amount,
		})
		if err != nil {
			return err
		}

		
		acc2,err:= q.GetAccountForUpdate(ctx, args.ToAccountID)
		if err != nil {
			return err
		}
	
		result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID: args.ToAccountID,
			Balance: acc2.Balance + args.Amount,
		})
		if err != nil {
			return err
		}
		//the above implementation fails beacause the two goroutines are fetching the same account at the same time and it is non blocking

		return nil
	})

	return result, err
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID int64 `json:"to_account_id"`
	Amount int64 `json:"amount"`
}

type TransferTxResults struct {
	Transfer Transfer `json:"transfer"`
	FromAccount Account `json:"from_account"`
	ToAccount Account `json:"to_account"`
	FromEntry Entry `json:"from_entry"`
	ToEntry Entry `json:"to_entry"`
}
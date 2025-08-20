package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Store provides all functions to execute db queries and transactions
type Store struct {
	*Queries
	db *pgxpool.Pool
}

// NewStore creates a new store
func NewStore(db *pgxpool.Pool) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {

	// start a transaction
	tx, err := store.db.Begin(ctx)
	if err != nil {
		return err
	}

	// execute the function
	if err := fn(New(tx)); err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	// commit the transaction
	return tx.Commit(ctx)
}

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx performs a money transfer from one account to the other
// It creates a transfer record, add account entries, and update accounts balance within a single database transaction
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// create a transfer record
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        pgtype.Int8{Int64: arg.Amount, Valid: true},
		})
		if err != nil {
			return err
		}

		// create a from account entry
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    pgtype.Int8{Int64: -arg.Amount, Valid: true},
		})

		if err != nil {
			return err
		}

		// create a to account entry
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    pgtype.Int8{Int64: arg.Amount, Valid: true},
		})

		if err != nil {
			return err
		}

		// get current accounts with locking to prevent race conditions
		// lock accounts in a consistent order to prevent deadlocks
		// if arg.FromAccountID < arg.ToAccountID {
		// 	result.FromAccount, err = q.GetAccountForUpdate(ctx, arg.FromAccountID)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	result.ToAccount, err = q.GetAccountForUpdate(ctx, arg.ToAccountID)
		// 	if err != nil {
		// 		return err
		// 	}
		// } else {
		// 	result.ToAccount, err = q.GetAccountForUpdate(ctx, arg.ToAccountID)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	result.FromAccount, err = q.GetAccountForUpdate(ctx, arg.FromAccountID)
		// 	if err != nil {
		// 		return err
		// 	}
		// }

		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}

		return err
	})

	return result, err
}

func addMoney(ctx context.Context, q *Queries, accountID1 int64, amount1 int64, accountID2 int64, amount2 int64) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: pgtype.Int8{Int64: amount1, Valid: true},
	})

	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: pgtype.Int8{Int64: amount2, Valid: true},
	})
	if err != nil {
		return
	}

	return
}

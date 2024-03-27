package database

import (
	"context"
	"database/sql"
	"fmt"

	"log"
)

type Store struct {
	*Queries
	db *sql.DB
}

var transactionKey = struct{}{}

func NewPostgresStore(connectionString string) (*Store, error) {

	conn, err := sql.Open("postgres", connectionString)

	if err != nil {
		log.Fatal(err)
	}

	err = conn.Ping()
	if err != nil {
		return nil, err
	}
	return &Store{
		Queries: New(conn),
		db:      conn,
	}, nil
}

// lol didn't know where else to put this
// exec fn  within a db transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	// get a tx for making transaction requests
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	//exec the callback fn and return an error if it fails
	err = fn(q)
	//rollback the transaction in case of failure
	if err != nil {
		//if the rollback also fails return both errors
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf(" tx err:%v,rb err:%v", err, rbErr)
		}
		return err
	}
	//commit the transaction and return its error if it occurs
	return tx.Commit()
}

/* transfer funds,create entries and update account balances as one singular op*/
func (store *Store) TransferTx(ctx context.Context, params TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			Amount:     params.Amount,
			SenderID:   params.SenderID,
			ReceiverID: params.ReceiverID,
		})
		if err != nil {
			return err
		}
		result.SenderEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: result.Transfer.SenderID,
			Amount:    -result.Transfer.Amount,
		})
		if err != nil {
			return err
		}

		result.ReceiverEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: result.Transfer.ReceiverID,
			Amount:    result.Transfer.Amount,
		})
		if err != nil {
			return err
		}

		result.SenderAccount, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
			ID:     params.SenderID,
			Amount: -params.Amount,
		})
		if err != nil {
			return err
		}
		result.ReceiverAccount, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
			ID:     params.ReceiverID,
			Amount: params.Amount,
		})
		if err != nil {
			return err
		}
		return nil
	})

	return result, err
}

type TransferTxParams struct {
	Amount     int64 `json:"amount"`
	SenderID   int64 `json:"sender_id"`
	ReceiverID int64 `json:"receiver_id"`
}

type TransferTxResult struct {
	Transfer        Transfer `json:"transfer"`
	SenderAccount   Account  `json:"sender_account"`
	ReceiverAccount Account  `json:"receiver_account"`
	SenderEntry     Entry    `json:"sender_entry"`
	ReceiverEntry   Entry    `json:"receiver_entry"`
}

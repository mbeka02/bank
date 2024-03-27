// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: transfers.sql

package database

import (
	"context"
)

const createTransfer = `-- name: CreateTransfer :one
INSERT INTO transfers (amount, sender_id, receiver_id)
VALUES ($1,$2,$3)
RETURNING id, amount, sender_id, receiver_id, created_at
`

type CreateTransferParams struct {
	Amount     int64
	SenderID   int64
	ReceiverID int64
}

func (q *Queries) CreateTransfer(ctx context.Context, arg CreateTransferParams) (Transfer, error) {
	row := q.db.QueryRowContext(ctx, createTransfer, arg.Amount, arg.SenderID, arg.ReceiverID)
	var i Transfer
	err := row.Scan(
		&i.ID,
		&i.Amount,
		&i.SenderID,
		&i.ReceiverID,
		&i.CreatedAt,
	)
	return i, err
}

const deleteTransfer = `-- name: DeleteTransfer :exec
DELETE FROM transfers WHERE id=$1
`

func (q *Queries) DeleteTransfer(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteTransfer, id)
	return err
}

const getTransfer = `-- name: GetTransfer :one
SELECT id, amount, sender_id, receiver_id, created_at FROM transfers WHERE id=$1 LIMIT 1
`

func (q *Queries) GetTransfer(ctx context.Context, id int64) (Transfer, error) {
	row := q.db.QueryRowContext(ctx, getTransfer, id)
	var i Transfer
	err := row.Scan(
		&i.ID,
		&i.Amount,
		&i.SenderID,
		&i.ReceiverID,
		&i.CreatedAt,
	)
	return i, err
}

const getTransfers = `-- name: GetTransfers :many
SELECT id, amount, sender_id, receiver_id, created_at FROM transfers ORDER BY id LIMIT $1 OFFSET $2
`

type GetTransfersParams struct {
	Limit  int32
	Offset int32
}

func (q *Queries) GetTransfers(ctx context.Context, arg GetTransfersParams) ([]Transfer, error) {
	rows, err := q.db.QueryContext(ctx, getTransfers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Transfer
	for rows.Next() {
		var i Transfer
		if err := rows.Scan(
			&i.ID,
			&i.Amount,
			&i.SenderID,
			&i.ReceiverID,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
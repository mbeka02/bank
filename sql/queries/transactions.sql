-- name: GetTransactions :many
SELECT * FROM transactions ORDER BY id LIMIT $1 OFFSET $2;

-- name: GetTransaction :one
SELECT * FROM transactions WHERE id=$1 LIMIT 1;


-- name: CreateTransaction :one
INSERT INTO transactions (amount, sender_id, receiver_id)
VALUES ($1,$2,$3)
RETURNING *;

-- name: DeleteTransaction :exec
DELETE FROM transactions WHERE id=$1;

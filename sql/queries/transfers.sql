
-- name: GetTransfers :many
SELECT * FROM transfer
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: GetTransfer :one
SELECT * FROM transfer WHERE id=$1 LIMIT 1;


-- name: CreateTransfer :one
INSERT INTO transfer (amount, sender_id, receiver_id)
VALUES ($1,$2,$3)
RETURNING *;

-- name: DeleteTransfer :exec
DELETE FROM transfer WHERE id=$1;

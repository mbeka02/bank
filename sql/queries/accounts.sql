-- name: GetAccounts :many
SELECT * FROM accounts
WHERE owner=$1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: GetAccount :one
SELECT * FROM accounts WHERE id=$1 LIMIT 1;

-- name: GetAccountForUpdate :one
SELECT * FROM accounts WHERE id=$1 LIMIT 1 FOR NO KEY UPDATE;

-- name: UpdateAccount :one
UPDATE accounts SET balance=$1 WHERE id=$2 RETURNING *;

-- name: UpdateAccountBalance :one
UPDATE accounts SET balance=balance +sqlc.arg(amount) WHERE id=sqlc.arg(id) RETURNING *;


-- name: CreateAccount :one
INSERT INTO accounts (owner, balance, currency)
VALUES ($1,$2,$3)
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts WHERE id=$1;

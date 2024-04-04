-- name: GetUsers :many
SELECT * FROM users
ORDER BY user_name
LIMIT $1
OFFSET $2;

-- name: GetUser :one
SELECT * FROM users WHERE user_name=$1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (user_name, full_name, password , email)
VALUES ($1,$2,$3,$4)
RETURNING *;


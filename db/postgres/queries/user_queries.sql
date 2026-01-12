-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users 
LIMIT $1 OFFSET $2;

-- name: CountUsers :one
SELECT COUNT(*) FROM users;

-- name: CreateUser :one
INSERT INTO users (name, email, password_hash, role, status, token_version) 
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id;

-- name: UpdateUser :one
UPDATE users 
SET name = $1, email = $2, token_version = $3 
WHERE id = $4
RETURNING id;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: UpdateTokenVersion :exec
UPDATE users SET token_version = token_version + 1 WHERE id = $1;

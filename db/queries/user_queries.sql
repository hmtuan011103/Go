-- name: GetUserByID :one
SELECT * FROM users WHERE id = ? LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = ? LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users 
LIMIT ? OFFSET ?;

-- name: CountUsers :one
SELECT COUNT(*) FROM users;

-- name: CreateUser :execresult
INSERT INTO users (name, email, password_hash, role, status, token_version) 
VALUES (?, ?, ?, ?, ?, ?);

-- name: UpdateUser :execresult
UPDATE users 
SET name = ?, email = ?, token_version = ? 
WHERE id = ?;

-- name: DeleteUser :execresult
DELETE FROM users WHERE id = ?;

-- name: UpdateTokenVersion :execresult
UPDATE users SET token_version = token_version + 1 WHERE id = ?;

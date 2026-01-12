-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (user_id, token, expires_at) 
VALUES ($1, $2, $3)
RETURNING id;

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens WHERE token = $1 LIMIT 1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens SET revoked_at = $1 WHERE token = $2;

-- name: RevokeAllUserTokens :exec
UPDATE refresh_tokens SET revoked_at = $1 WHERE user_id = $2 AND revoked_at IS NULL;

-- name: CreateRefreshToken :execresult
INSERT INTO refresh_tokens (user_id, token, expires_at) 
VALUES (?, ?, ?);

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens WHERE token = ? LIMIT 1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens SET revoked_at = ? WHERE token = ?;

-- name: RevokeAllUserTokens :exec
UPDATE refresh_tokens SET revoked_at = ? WHERE user_id = ? AND revoked_at IS NULL;

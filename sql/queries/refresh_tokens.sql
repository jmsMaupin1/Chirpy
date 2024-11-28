-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUserFromRefreshToken :one
SELECT u.id, u.email, r.expires_at, r.revoked_at FROM refresh_tokens r
JOIN users u on r.user_id = u.id
WHERE token = $1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = $2, updated_at = $3
WHERE user_id = $1;

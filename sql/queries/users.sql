-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: UpdateUser :one
UPDATE users u
SET email = $2, hashed_password = $3, updated_at = $4
WHERE id = $1
RETURNING u.id, u.created_at, u.updated_At, u.email;

-- name: DeleteUsers :exec
DELETE FROM users;

-- name: SetUserChirpyRed :one
UPDATE users u
SET is_chirpy_red = true
WHERE id = $1
RETURNING u.id, u.created_at, u.updated_at, u.email;

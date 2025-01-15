-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetAllChirps :many
SELECT * FROM chirps
ORDER BY 
	CASE WHEN $1 = 'asc' THEN created_at END ASC,
	CASE WHEN $1 = 'desc' THEN created_at END DESC;

-- name: GetChirp :one
SELECT * FROM chirps
WHERE id = $1;

-- name: GetChirpsByAuthor :many
SELECT * FROM chirps
WHERE user_id = $1
ORDER BY
	CASE WHEN $2 = 'asc' THEN created_at END ASC,
	CASE WHEN $2 = 'desc' THEN created_at END DESC;

-- name: DeleteChirp :exec
DELETE FROM chirps
WHERE id = $1;

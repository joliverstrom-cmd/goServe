-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, body, user_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: GetAllChirps :many
SELECT * FROM posts
ORDER BY created_at ASC;

-- name: GetOneChirp :one
SELECT * FROM posts
WHERE id = $1;

-- name: DeleteChirp :exec
DELETE from posts
WHERE id = $1 AND user_id = $2;

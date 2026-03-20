-- name: CreateUser :exec
INSERT INTO users (id, username, email, password_hash, display_name, avatar_url, bio, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);

-- name: GetUserByID :one
SELECT id, username, email, password_hash, display_name, avatar_url, bio, created_at, updated_at
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, username, email, password_hash, display_name, avatar_url, bio, created_at, updated_at
FROM users
WHERE email = $1;

-- name: GetUserByUsername :one
SELECT id, username, email, password_hash, display_name, avatar_url, bio, created_at, updated_at
FROM users
WHERE username = $1;

-- name: UpdateUser :exec
UPDATE users
SET display_name = $2, avatar_url = $3, bio = $4, updated_at = $5
WHERE id = $1;

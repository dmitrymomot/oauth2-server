-- name: CreateUser :one
INSERT INTO users (email, password) VALUES ($1, $2) RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: UpdateUserEmail :one
UPDATE users SET email = $1 WHERE id = $2 RETURNING *;

-- name: UpdateUserPassword :one
UPDATE users SET password = $1 WHERE id = $2 RETURNING *;

-- name: UpdateUserVerifiedAt :exec
UPDATE users SET verified_at = now() WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;
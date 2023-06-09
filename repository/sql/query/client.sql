-- name: CreateClient :one
INSERT INTO clients (id, secret, domain, is_public, user_id, allowed_grants, scope) 
VALUES (@id, @secret, @domain, @is_public, @user_id, @allowed_grants, @scope) RETURNING *;

-- name: GetClientByID :one
SELECT * FROM clients WHERE id = $1;

-- name: GetClientByUserID :many
SELECT * FROM clients WHERE user_id = $1 ORDER BY created_at DESC;

-- name: UpdateClientSecret :one
UPDATE clients SET secret = $1 WHERE id = $2 RETURNING *;

-- name: DeleteClient :exec
DELETE FROM clients WHERE id = $1;
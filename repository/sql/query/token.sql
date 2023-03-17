-- name: CreateToken :one
INSERT INTO tokens (
    client_id, 
    user_id, 
    redirect_uri, 
    scope, 
    code,
    code_created_at, 
    code_expires_in,
    code_challenge,
    code_challenge_method,
    access,
    access_created_at,
    access_expires_in,
    refresh,
    refresh_created_at,
    refresh_expires_in
) VALUES (
    @client_id, 
    @user_id, 
    @redirect_uri, 
    @scope, 
    @code,
    @code_created_at, 
    @code_expires_in,
    @code_challenge,
    @code_challenge_method,
    @access,
    @access_created_at,
    @access_expires_in,
    @refresh,
    @refresh_created_at,
    @refresh_expires_in
) RETURNING *;

-- name: GetTokenByCode :one
SELECT * FROM tokens WHERE code = @code;

-- name: GetTokenByAccess :one
SELECT * FROM tokens WHERE access = @access;

-- name: GetTokenByRefresh :one
SELECT * FROM tokens WHERE refresh = @refresh;

-- name: DeleteByCode :exec
DELETE FROM tokens WHERE code = @code;

-- name: DeleteByAccess :exec
DELETE FROM tokens WHERE access = @access;

-- name: DeleteByRefresh :exec
DELETE FROM tokens WHERE refresh = @refresh;

-- name: DeleteExpiredTokens :exec
DELETE FROM tokens 
WHERE (code_expires_in > 0 AND code_created_at + code_expires_in * interval '1 second' < now())
OR (refresh_expires_in > 0 AND refresh_created_at + refresh_expires_in * interval '1 second' < now());
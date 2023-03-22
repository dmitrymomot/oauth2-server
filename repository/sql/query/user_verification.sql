-- name: GetUserVerificationByUserID :one
SELECT *
FROM user_verifications
WHERE request_type = @request_type
    AND user_id = @user_id
    AND expires_at > now()
ORDER BY created_at DESC
LIMIT 1;

-- name: GetUserVerificationByEmail :one
SELECT *
FROM user_verifications
WHERE request_type = @request_type
    AND email = @email
    AND expires_at > now()
ORDER BY created_at DESC
LIMIT 1;

-- name: GetVerificationByUserIDAndEmail :one
SELECT *
FROM user_verifications
WHERE request_type = @request_type
    AND user_id = @user_id
    AND email = @email
    AND expires_at > now()
ORDER BY created_at DESC
LIMIT 1;

-- name: CreateUserVerification :exec
INSERT INTO user_verifications (request_type, user_id, email, verification_code, expires_at)
VALUES (
        @request_type,
        @user_id,
        @email,
        @verification_code,
        @expires_at
    ) ON CONFLICT (request_type, user_id, email) DO
UPDATE
SET verification_code = @verification_code;

-- name: DeleteUserVerificationsByEmail :exec
DELETE FROM user_verifications
WHERE request_type = @request_type
    AND email = @email;

-- name: DeleteUserVerificationsByUserID :exec
DELETE FROM user_verifications
WHERE request_type = @request_type
    AND user_id = @user_id;

-- name: CleanUpExpiredUserVerifications :exec
DELETE FROM user_verifications
WHERE expires_at < now();
-- +migrate Up
-- +migrate StatementBegin
CREATE TYPE user_verification_request_type AS ENUM (
  'email_change', 
  'email_verification',
  'password_reset',
  'delete_account'
);

CREATE TABLE IF NOT EXISTS user_verifications (
    request_type user_verification_request_type NOT NULL,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    email VARCHAR DEFAULT NULL,
    verification_code bytea NOT NULL,
    expires_at TIMESTAMP NOT NULL DEFAULT now() + interval '15 minutes',
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    PRIMARY KEY(request_type, user_id, email)
);
CREATE INDEX email_verifications_user_id ON user_verifications USING BTREE (request_type, user_id);
CREATE INDEX email_verifications_email ON user_verifications USING BTREE (request_type, email) WHERE email IS NOT NULL;
-- +migrate StatementEnd

-- +migrate Down
DROP TABLE IF EXISTS user_verifications;
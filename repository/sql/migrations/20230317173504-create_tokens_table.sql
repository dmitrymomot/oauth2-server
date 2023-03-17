-- +migrate Up
-- +migrate StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS tokens (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    client_id VARCHAR NOT NULL REFERENCES clients (id) ON DELETE CASCADE,
    user_id uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    redirect_uri VARCHAR NOT NULL,
    scope VARCHAR NOT NULL DEFAULT '',
    code VARCHAR NOT NULL DEFAULT '',
    code_created_at TIMESTAMP DEFAULT NULL,
    code_expires_in BIGINT NOT NULL DEFAULT 0,
    code_challenge VARCHAR NOT NULL DEFAULT '',
    code_challenge_method VARCHAR NOT NULL DEFAULT '',
    access VARCHAR NOT NULL DEFAULT '',
    access_created_at TIMESTAMP DEFAULT NULL,
    access_expires_in BIGINT NOT NULL DEFAULT 0,
    refresh VARCHAR NOT NULL DEFAULT '',
    refresh_created_at TIMESTAMP DEFAULT NULL,
    refresh_expires_in BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
CREATE INDEX tokens_client_id ON tokens USING BTREE (client_id);
CREATE INDEX tokens_user_id ON tokens USING BTREE (user_id);
CREATE UNIQUE INDEX tokens_code ON tokens USING BTREE (code) WHERE code <> '';
CREATE UNIQUE INDEX tokens_access ON tokens USING BTREE (access) WHERE access <> '';
CREATE UNIQUE INDEX tokens_refresh ON tokens USING BTREE (refresh) WHERE refresh <> '';
-- +migrate StatementEnd

-- +migrate Down
DROP TABLE IF EXISTS tokens;
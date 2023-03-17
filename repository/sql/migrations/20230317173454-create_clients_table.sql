-- +migrate Up
-- +migrate StatementBegin
CREATE TABLE IF NOT EXISTS clients (
    id VARCHAR PRIMARY KEY,
    secret bytea NOT NULL,
    domain VARCHAR NOT NULL,
    is_public BOOLEAN NOT NULL DEFAULT FALSE,
    user_id uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
CREATE INDEX clients_domain ON clients USING BTREE (domain);
CREATE INDEX clients_user_id ON clients USING BTREE (user_id);
-- +migrate StatementEnd

-- +migrate Down
DROP TABLE IF EXISTS clients;
-- +goose Up
CREATE TABLE sessions
(
    token  TEXT PRIMARY KEY,
    data   BYTEA                    NOT NULL,
    expiry TIMESTAMP WITH TIME ZONE NOT NULL
);
CREATE INDEX sessions_expiry_idx ON sessions (expiry);

-- +goose Down
DROP TABLE IF EXISTS sessions;

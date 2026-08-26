-- +goose Up
CREATE
EXTENSION IF NOT EXISTS citext;
CREATE TABLE users
(
    id            UUID PRIMARY KEY         DEFAULT gen_random_uuid(),
    name          VARCHAR(255)  NOT NULL,
    email         citext UNIQUE NOT NULL,
    password_hash VARCHAR(255)  NOT NULL,
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at    TIMESTAMP WITH TIME ZONE
);
CREATE INDEX idx_users_deleted_at ON users (deleted_at);

-- +goose Down
DROP TABLE users;
DROP
EXTENSION IF EXISTS citext;

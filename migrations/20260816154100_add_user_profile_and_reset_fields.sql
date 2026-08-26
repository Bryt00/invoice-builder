-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_profile_complete BOOLEAN DEFAULT FALSE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_activated BOOLEAN DEFAULT FALSE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS activation_token VARCHAR(255);
ALTER TABLE users ADD COLUMN IF NOT EXISTS activation_expiry TIMESTAMP WITH TIME ZONE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS reset_token VARCHAR(255);
ALTER TABLE users ADD COLUMN IF NOT EXISTS reset_expiry TIMESTAMP WITH TIME ZONE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN IF EXISTS is_profile_complete;
ALTER TABLE users DROP COLUMN IF EXISTS is_activated;
ALTER TABLE users DROP COLUMN IF EXISTS activation_token;
ALTER TABLE users DROP COLUMN IF EXISTS activation_expiry;
ALTER TABLE users DROP COLUMN IF EXISTS reset_token;
ALTER TABLE users DROP COLUMN IF EXISTS reset_expiry;
-- +goose StatementEnd

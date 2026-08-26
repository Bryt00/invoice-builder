-- +goose Up
-- +goose StatementBegin
ALTER TABLE business_profiles ADD COLUMN IF NOT EXISTS registration_number VARCHAR(100);
ALTER TABLE business_profiles ADD COLUMN IF NOT EXISTS registration_date VARCHAR(50);
ALTER TABLE business_profiles ADD COLUMN IF NOT EXISTS business_type VARCHAR(100);
ALTER TABLE business_profiles ADD COLUMN IF NOT EXISTS registered_address TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE business_profiles DROP COLUMN IF EXISTS registration_number;
ALTER TABLE business_profiles DROP COLUMN IF EXISTS registration_date;
ALTER TABLE business_profiles DROP COLUMN IF EXISTS business_type;
ALTER TABLE business_profiles DROP COLUMN IF EXISTS registered_address;
-- +goose StatementEnd

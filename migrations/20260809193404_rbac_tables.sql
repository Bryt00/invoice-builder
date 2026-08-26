-- +goose Up
-- +goose StatementBegin
CREATE TABLE roles
(
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(255) UNIQUE NOT NULL,
    description TEXT
);

CREATE TABLE permissions
(
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(255) UNIQUE NOT NULL,
    description TEXT
);

CREATE TABLE role_permissions
(
    role_id       UUID NOT NULL REFERENCES roles (id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions (id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);

ALTER TABLE users
    ADD COLUMN role_id UUID;

-- Insert default roles
INSERT INTO roles (id, name, description)
VALUES (gen_random_uuid(), 'Admin', 'Administrator with full access'),
       (gen_random_uuid(), 'User', 'Standard user with basic access');

-- Update existing users to have the 'User' role
UPDATE users
SET role_id = (SELECT id FROM roles WHERE name = 'User');

-- Make role_id NOT NULL now that existing rows are populated
ALTER TABLE users
    ALTER COLUMN role_id SET NOT NULL;
ALTER TABLE users
    ADD CONSTRAINT fk_user_role FOREIGN KEY (role_id) REFERENCES roles (id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_user_role;
ALTER TABLE users DROP COLUMN IF EXISTS role_id;
DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS roles;
-- +goose StatementEnd

-- +goose Up
ALTER TABLE users
ADD COLUMN is_active BOOLEAN NOT NULL DEFAULT TRUE,
ADD COLUMN disabled_at TIMESTAMP;

-- +goose Down
ALTER TABLE users
DROP COLUMN disabled_at,
DROP COLUMN is_active;

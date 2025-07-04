-- +goose Up
UPDATE countries SET sort_order = 100 WHERE sort_order IS NULL;
UPDATE countries SET created_at = CURRENT_TIMESTAMP WHERE created_at IS NULL;
UPDATE countries SET updated_at = CURRENT_TIMESTAMP WHERE updated_at IS NULL;

ALTER TABLE countries
    ALTER COLUMN sort_order SET NOT NULL,
    ALTER COLUMN created_at SET NOT NULL,
    ALTER COLUMN updated_at SET NOT NULL;

-- +goose Down
ALTER TABLE countries
    ALTER COLUMN sort_order DROP NOT NULL,
    ALTER COLUMN created_at DROP NOT NULL,
    ALTER COLUMN updated_at DROP NOT NULL;

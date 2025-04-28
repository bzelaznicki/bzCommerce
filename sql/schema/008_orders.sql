-- +goose Up
ALTER TABLE orders
ALTER COLUMN total_price TYPE numeric(10,2);

-- +goose Down
ALTER TABLE orders
ALTER COLUMN total_price NUMERIC;
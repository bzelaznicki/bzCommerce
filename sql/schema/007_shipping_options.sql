-- +goose Up

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE shipping_options (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    description TEXT,
    price NUMERIC(10, 2) NOT NULL,
    estimated_days TEXT,
    sort_order INT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE orders
ADD COLUMN shipping_method_id UUID REFERENCES shipping_options(id) ON DELETE SET NULL;
ALTER TABLE orders
ADD COLUMN shipping_price NUMERIC(10, 2);

-- +goose Down

-- Drop columns from orders table
ALTER TABLE orders
DROP COLUMN IF EXISTS shipping_method_id;

ALTER TABLE orders
DROP COLUMN IF EXISTS shipping_price;

-- Drop index if it exists
DROP INDEX IF EXISTS idx_shipping_options_is_active;

-- Drop the shipping_options table
DROP TABLE IF EXISTS shipping_options;

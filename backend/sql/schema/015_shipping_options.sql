-- +goose Up

-- Create the sequence for sort_order auto-increment if not already present
CREATE SEQUENCE IF NOT EXISTS shipping_options_sort_order_seq;

-- Backfill NULL sort_order with next sequence values
UPDATE shipping_options
SET sort_order = nextval('shipping_options_sort_order_seq')
WHERE sort_order IS NULL;

-- Make sort_order NOT NULL
ALTER TABLE shipping_options
ALTER COLUMN sort_order SET NOT NULL;

-- Set default auto-increment for sort_order
ALTER TABLE shipping_options
ALTER COLUMN sort_order SET DEFAULT nextval('shipping_options_sort_order_seq');

-- Backfill NULL estimated_days with a default string
UPDATE shipping_options
SET estimated_days = 'N/A'
WHERE estimated_days IS NULL;

-- Make estimated_days NOT NULL
ALTER TABLE shipping_options
ALTER COLUMN estimated_days SET NOT NULL;

-- +goose Down

-- Remove default for sort_order
ALTER TABLE shipping_options
ALTER COLUMN sort_order DROP DEFAULT;

-- Allow sort_order to be NULL again
ALTER TABLE shipping_options
ALTER COLUMN sort_order DROP NOT NULL;

-- Allow estimated_days to be NULL again
ALTER TABLE shipping_options
ALTER COLUMN estimated_days DROP NOT NULL;

-- Drop the sequence for sort_order
DROP SEQUENCE IF EXISTS shipping_options_sort_order_seq;

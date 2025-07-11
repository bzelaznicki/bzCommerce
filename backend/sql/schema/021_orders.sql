-- +goose Up

-- 1) Create payment_status ENUM
CREATE TYPE payment_status AS ENUM (
  'pending',
  'paid',
  'failed',
  'refunded'
);

-- 2) Add payment_status column to orders
ALTER TABLE orders
ADD COLUMN payment_status payment_status NOT NULL DEFAULT 'pending';

-- 3) Make created_at and updated_at NOT NULL
ALTER TABLE orders
ALTER COLUMN created_at SET NOT NULL,
ALTER COLUMN updated_at SET NOT NULL;

-- 4) Create a reusable trigger function to auto-update updated_at
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$ BEGIN NEW.updated_at = now();RETURN NEW; END; $$ LANGUAGE plpgsql;


-- 5) Create the trigger on orders table
CREATE TRIGGER set_updated_at
BEFORE UPDATE ON orders
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

-- +goose Down

-- 1) Drop the trigger
DROP TRIGGER IF EXISTS set_updated_at ON orders;

-- 2) Drop the trigger function
DROP FUNCTION IF EXISTS set_updated_at();

-- 3) Remove payment_status column
ALTER TABLE orders
DROP COLUMN payment_status;

-- 4) Drop the ENUM type
DROP TYPE payment_status;

-- 5) Remove NOT NULL constraints on timestamps
ALTER TABLE orders
ALTER COLUMN created_at DROP NOT NULL,
ALTER COLUMN updated_at DROP NOT NULL;
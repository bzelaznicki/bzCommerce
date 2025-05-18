-- +goose Up
CREATE TABLE users_addresses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    name TEXT NOT NULL,
    address TEXT,
    city TEXT,
    postal_code TEXT,
    country_id UUID NOT NULL,
    phone TEXT NOT NULL,
    is_shipping BOOLEAN NOT NULL DEFAULT FALSE,
    is_shipping_default BOOLEAN NOT NULL DEFAULT FALSE,
    is_billing BOOLEAN NOT NULL DEFAULT FALSE,
    is_billing_default BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (country_id) REFERENCES countries(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX uq_users_addresses_shipping_default 
    ON users_addresses (user_id)
    WHERE is_shipping_default = TRUE;
CREATE UNIQUE INDEX uq_users_addresses_billing_default
    ON users_addresses (user_id)
    WHERE is_billing_default = TRUE;

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_user_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  UPDATE users
  SET updated_at = CURRENT_TIMESTAMP
  WHERE id = NEW.user_id;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_update_user_timestamp
AFTER INSERT OR UPDATE OR DELETE ON users_addresses
FOR EACH ROW EXECUTE FUNCTION update_user_timestamp();
-- +goose StatementEnd

-- +goose Down
DROP TRIGGER IF EXISTS trg_update_user_timestamp ON users_addresses;
DROP FUNCTION IF EXISTS update_user_timestamp;
DROP INDEX IF EXISTS uq_users_addresses_shipping_default;
DROP INDEX IF EXISTS uq_users_addresses_billing_default;
DROP TABLE IF EXISTS users_addresses;
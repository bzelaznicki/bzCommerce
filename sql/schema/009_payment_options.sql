-- +goose Up

CREATE TABLE payment_options (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    sort_order INT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE orders
ADD COLUMN payment_option_id UUID REFERENCES payment_options(id) ON DELETE SET NULL;

-- +goose Down

ALTER TABLE orders
DROP COLUMN IF EXISTS payment_option_id;

DROP TABLE IF EXISTS payment_options;

-- +goose Up

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE order_status AS ENUM ('pending', 'paid', 'processing', 'shipped', 'cancelled', 'refunded');

CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID,
    status order_status NOT NULL DEFAULT 'pending',
    total_price NUMERIC NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    customer_email TEXT NOT NULL,
    shipping_name TEXT,
    shipping_address TEXT,
    shipping_city TEXT,
    shipping_postal_code TEXT,
    shipping_country TEXT,
    shipping_phone TEXT NOT NULL,
    billing_name TEXT,
    billing_address TEXT,
    billing_city TEXT,
    billing_postal_code TEXT,
    billing_country TEXT,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE orders_variants (
    order_id UUID NOT NULL,
    product_variant_id UUID NOT NULL,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    price_per_item NUMERIC NOT NULL,
    total_price NUMERIC NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (order_id, product_variant_id),
    CONSTRAINT fk_order FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE CASCADE,
    CONSTRAINT fk_variant FOREIGN KEY (product_variant_id) REFERENCES product_variants (id) ON DELETE CASCADE
);

CREATE INDEX idx_orders_user_id ON orders (user_id);
CREATE INDEX idx_orders_status ON orders (status);


-- +goose Down

DROP INDEX IF EXISTS idx_orders_user_id;
DROP INDEX IF EXISTS idx_orders_status;
DROP TABLE IF EXISTS orders_variants;
DROP TABLE IF EXISTS orders;
DROP TYPE IF EXISTS order_status;
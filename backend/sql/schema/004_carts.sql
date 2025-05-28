-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE cart_status AS ENUM ('new', 'abandoned', 'completed');

CREATE TABLE carts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID,
    status cart_status NOT NULL DEFAULT 'new',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE carts_variants (
    cart_id UUID NOT NULL,
    product_variant_id UUID NOT NULL,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    price_per_item NUMERIC NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (cart_id, product_variant_id),
    CONSTRAINT fk_cart FOREIGN KEY (cart_id) REFERENCES carts (id) ON DELETE CASCADE,
    CONSTRAINT fk_variant FOREIGN KEY (product_variant_id) REFERENCES product_variants (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS carts_variants;
DROP TABLE IF EXISTS carts;
DROP TYPE IF EXISTS cart_status;

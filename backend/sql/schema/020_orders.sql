-- +goose Up

ALTER TABLE orders
    ADD COLUMN shipping_country_id UUID NOT NULL,
    ADD COLUMN billing_country_id UUID NOT NULL;

ALTER TABLE orders
    ADD CONSTRAINT fk_shipping_country FOREIGN KEY (shipping_country_id) REFERENCES countries(id) ON DELETE RESTRICT,
    ADD CONSTRAINT fk_billing_country FOREIGN KEY (billing_country_id) REFERENCES countries(id) ON DELETE RESTRICT;

ALTER TABLE orders
    DROP COLUMN shipping_country,
    DROP COLUMN billing_country;

ALTER TABLE orders
    ALTER COLUMN status SET NOT NULL,
    ALTER COLUMN total_price SET NOT NULL,
    ALTER COLUMN customer_email SET NOT NULL,
    ALTER COLUMN shipping_name SET NOT NULL,
    ALTER COLUMN shipping_address SET NOT NULL,
    ALTER COLUMN shipping_city SET NOT NULL,
    ALTER COLUMN shipping_postal_code SET NOT NULL,
    ALTER COLUMN shipping_phone SET NOT NULL,
    ALTER COLUMN shipping_option_id SET NOT NULL,
    ALTER COLUMN shipping_price SET NOT NULL,
    ALTER COLUMN payment_option_id SET NOT NULL,
    ALTER COLUMN billing_name SET NOT NULL,
    ALTER COLUMN billing_address SET NOT NULL,
    ALTER COLUMN billing_city SET NOT NULL,
    ALTER COLUMN billing_postal_code SET NOT NULL;

-- +goose Down

ALTER TABLE orders
    ADD COLUMN shipping_country TEXT,
    ADD COLUMN billing_country TEXT;

ALTER TABLE orders
    DROP CONSTRAINT fk_shipping_country,
    DROP CONSTRAINT fk_billing_country;

ALTER TABLE orders
    DROP COLUMN shipping_country_id,
    DROP COLUMN billing_country_id;

ALTER TABLE orders
    ALTER COLUMN status DROP NOT NULL,
    ALTER COLUMN total_price DROP NOT NULL,
    ALTER COLUMN customer_email DROP NOT NULL,
    ALTER COLUMN shipping_name DROP NOT NULL,
    ALTER COLUMN shipping_address DROP NOT NULL,
    ALTER COLUMN shipping_city DROP NOT NULL,
    ALTER COLUMN shipping_postal_code DROP NOT NULL,
    ALTER COLUMN shipping_phone DROP NOT NULL,
    ALTER COLUMN shipping_option_id DROP NOT NULL,
    ALTER COLUMN shipping_price DROP NOT NULL,
    ALTER COLUMN payment_option_id DROP NOT NULL,
    ALTER COLUMN billing_name DROP NOT NULL,
    ALTER COLUMN billing_address DROP NOT NULL,
    ALTER COLUMN billing_city DROP NOT NULL,
    ALTER COLUMN billing_postal_code DROP NOT NULL;

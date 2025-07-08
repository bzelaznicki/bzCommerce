-- +goose Up
ALTER TABLE orders
    ADD COLUMN shipping_country_id UUID NOT NULL,
    ADD COLUMN billing_country_id UUID NOT NULL;

ALTER TABLE orders
    ADD CONSTRAINT fk_shipping_country FOREIGN KEY (shipping_country_id) REFERENCES countries (id) ON DELETE RESTRICT,
    ADD CONSTRAINT fk_billing_country FOREIGN KEY (billing_country_id) REFERENCES countries (id) ON DELETE RESTRICT;

ALTER TABLE orders
    DROP COLUMN shipping_country,
    DROP COLUMN billing_country;


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

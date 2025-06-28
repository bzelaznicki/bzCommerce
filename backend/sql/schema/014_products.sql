-- +goose Up
ALTER TABLE products
DROP CONSTRAINT fk_product_category;

ALTER TABLE products
ADD CONSTRAINT fk_product_category
FOREIGN KEY (category_id)
REFERENCES categories(id)
ON DELETE RESTRICT;

-- +goose Down
ALTER TABLE products
DROP CONSTRAINT fk_product_category;

ALTER TABLE products
ADD CONSTRAINT fk_product_category
FOREIGN KEY (category_id)
REFERENCES categories(id)
ON DELETE CASCADE;

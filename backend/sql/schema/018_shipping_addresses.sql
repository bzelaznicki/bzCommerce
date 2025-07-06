-- +goose Up
ALTER TABLE users_addresses
DROP CONSTRAINT users_addresses_country_id_fkey;

ALTER TABLE users_addresses
ADD CONSTRAINT users_addresses_country_id_fkey
FOREIGN KEY (country_id) REFERENCES countries(id) ON DELETE RESTRICT;

-- +goose Down
ALTER TABLE users_addresses
DROP CONSTRAINT users_addresses_country_id_fkey;

ALTER TABLE users_addresses
ADD CONSTRAINT users_addresses_country_id_fkey
FOREIGN KEY (country_id) REFERENCES countries(id) ON DELETE CASCADE;

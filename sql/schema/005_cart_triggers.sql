-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_cart_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  UPDATE carts
  SET updated_at = CURRENT_TIMESTAMP
  WHERE id = NEW.cart_id;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_update_cart_timestamp
AFTER INSERT OR UPDATE OR DELETE ON carts_variants
FOR EACH ROW EXECUTE FUNCTION update_cart_timestamp();
-- +goose StatementEnd

-- +goose Down
DROP TRIGGER IF EXISTS trg_update_cart_timestamp ON carts_variants;
DROP FUNCTION IF EXISTS update_cart_timestamp;

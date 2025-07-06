-- +goose Up

CREATE SEQUENCE IF NOT EXISTS countries_sort_order_seq;


SELECT setval(
  'countries_sort_order_seq',
  COALESCE((SELECT MAX(sort_order) FROM countries), 0)
);


UPDATE countries
SET sort_order = nextval('countries_sort_order_seq')
WHERE sort_order IS NULL;


ALTER TABLE countries
ALTER COLUMN sort_order SET DEFAULT nextval('countries_sort_order_seq');

-- +goose Down


ALTER TABLE countries
ALTER COLUMN sort_order DROP DEFAULT;


DROP SEQUENCE IF EXISTS countries_sort_order_seq;

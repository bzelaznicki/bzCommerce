-- name: GetProductBySlug :one
SELECT * FROM products WHERE slug = $1;

-- name: GetProductVariantsByProductSlug :many
SELECT * FROM product_variants WHERE product_id = $1;

-- name: ListProductsByCategory :many
SELECT * FROM products WHERE category_id = (SELECT id FROM categories WHERE slug = $1) ORDER BY name ASC;

-- name: ListProducts :many
SELECT * FROM products ORDER BY name ASC;
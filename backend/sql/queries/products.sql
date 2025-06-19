-- name: GetProductBySlug :one
SELECT * FROM products WHERE slug = sqlc.arg(slug);

-- name: GetProductById :one
SELECT * FROM products WHERE id = sqlc.arg(id);


-- name: GetProductVariantsByProductSlug :many
SELECT * FROM product_variants WHERE product_id = sqlc.arg(product_id);

-- name: ListProductsWithCategory :many
SELECT
  p.id,
  p.name,
  p.slug,
  p.description,
  p.image_url,
  p.category_id,
  c.name AS category_name,
  c.slug AS category_slug
FROM products p
JOIN categories c ON p.category_id = c.id
ORDER BY p.name;

-- name: ListProductsWithCategoryPaginated :many
SELECT
  p.id,
  p.name,
  p.slug,
  p.description,
  p.image_url,
  p.category_id,
  c.name AS category_name,
  c.slug AS category_slug
FROM products p
JOIN categories c ON p.category_id = c.id
ORDER BY p.name
LIMIT $1 OFFSET $2;

-- name: ListProductsByCategory :many
SELECT * FROM products
WHERE category_id = (SELECT id FROM categories WHERE slug = sqlc.arg(category_slug))
ORDER BY name ASC;

-- name: ListProductsByCategoryRecursive :many
WITH RECURSIVE subcategories AS (
  SELECT id FROM categories WHERE slug = sqlc.arg(slug)
  UNION ALL
  SELECT c.id
  FROM categories c
  INNER JOIN subcategories s ON c.parent_id = s.id
)
SELECT * FROM products
WHERE category_id IN (SELECT id FROM subcategories);



-- name: ListProducts :many
SELECT * FROM products ORDER BY name ASC;


-- name: GetProductVariantsByProductId :many
SELECT * FROM product_variants WHERE product_id = sqlc.arg(product_id);


-- name: DeleteProduct :exec
DELETE FROM products WHERE id = $1;

-- name: CreateProduct :one
INSERT INTO products (name, slug, description, image_url, category_id)
VALUES (
  sqlc.arg(name),
  sqlc.arg(slug),
  sqlc.arg(description),
  sqlc.arg(image_url),
  sqlc.arg(category_id)
)
RETURNING *;

-- name: UpdateProduct :exec
UPDATE products
SET
  name = sqlc.arg(name),
  slug = sqlc.arg(slug),
  description = sqlc.arg(description),
  category_id = sqlc.arg(category_id),
    image_url = sqlc.arg(image_url),
    updated_at = NOW()
WHERE id = sqlc.arg(id);


-- name: CreateProductVariant :one
INSERT INTO product_variants (product_id, variant_name, sku, price, image_url, stock_quantity)
VALUES (
  sqlc.arg(product_id),
    sqlc.arg(name),
  sqlc.arg(sku),
  sqlc.arg(price),
    sqlc.arg(image_url),
  sqlc.arg(stock_quantity)
)
RETURNING *;


-- name: GetVariantsByProductID :many
SELECT * FROM product_variants
WHERE product_id = sqlc.arg('product_id')
ORDER BY created_at;

-- name: GetVariantByID :one
SELECT * FROM product_variants
WHERE id = sqlc.arg('id');

-- name: CreateVariant :one
INSERT INTO product_variants (
  product_id, sku, price, stock_quantity, image_url, variant_name
)
VALUES (
  sqlc.arg('product_id'),
  sqlc.arg('sku'),
  sqlc.arg('price'),
  sqlc.arg('stock_quantity'),
  sqlc.arg('image_url'),
  sqlc.arg('variant_name')
)
RETURNING id;

-- name: UpdateVariant :exec
UPDATE product_variants
SET
  sku = sqlc.arg('sku'),
  price = sqlc.arg('price'),
  stock_quantity = sqlc.arg('stock_quantity'),
  image_url = sqlc.arg('image_url'),
  variant_name = sqlc.arg('variant_name'),
  updated_at = NOW()
WHERE id = sqlc.arg('id');

-- name: DeleteVariant :exec
DELETE FROM product_variants
WHERE id = sqlc.arg('id');


-- name: CountProducts :one
SELECT COUNT(*) FROM products;
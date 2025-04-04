-- name: GetCategories :many
SELECT * FROM categories ORDER BY name ASC;

-- name: GetCategoryById :one

SELECT * FROM categories WHERE id = sqlc.arg(id);

-- name: ListCategoriesWithParent :many
SELECT
  c.id,
  c.name,
  c.slug,
  c.description,
  c.parent_id,
  p.name AS parent_name,
  c.created_at,
  c.updated_at
FROM categories c
LEFT JOIN categories p ON c.parent_id = p.id
ORDER BY c.name;

-- name: CreateCategory :one
INSERT INTO categories (name, slug, description, parent_id)
VALUES (
    sqlc.arg(name),
    sqlc.arg(slug),
    sqlc.arg(description),
    sqlc.arg(parent_id)
)
RETURNING *;

-- name: UpdateCategoryById :exec
UPDATE categories
SET 
    name = sqlc.arg(name),
    slug = sqlc.arg(slug),
    description = sqlc.arg(description),
    parent_id = sqlc.arg(parent_id),
    updated_at = NOW()
    WHERE id = sqlc.arg(id);


-- name: DeleteCategoryById :exec
DELETE FROM categories WHERE id = sqlc.arg(id);
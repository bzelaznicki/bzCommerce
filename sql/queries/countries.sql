-- name: CreateCountry :one
INSERT INTO countries (name, iso_code, is_active, sort_order)
VALUES (
    sqlc.arg(name),
    sqlc.arg(iso_code),
    sqlc.arg(is_active),
    sqlc.arg(sort_order)
    )
    RETURNING *;

-- name: GetCountries :many
SELECT * FROM countries ORDER BY sort_order ASC;

-- name: GetCountryById :one
SELECT * FROM countries WHERE id = sqlc.arg(id);

-- name: UpdateCountryById :exec
UPDATE countries SET 
    name = sqlc.arg(name),
    iso_code = sqlc.arg(iso_code),
    is_active = sqlc.arg(is_active),
    sort_order = sqlc.arg(sort_order)

WHERE id = sqlc.arg(id);

-- name: DeleteCountryById :exec
DELETE FROM countries WHERE id = sqlc.arg(id);
-- name: GetUserAddresses :many
SELECT ua.*, c.name AS country FROM users_addresses ua
JOIN countries c ON ua.country_id = c.id
WHERE ua.user_id = sqlc.arg(user_id)
ORDER BY created_at DESC;


-- name: GetUserAddressById :one
SELECT ua.*, c.name AS country FROM users_addresses ua
JOIN countries c ON ua.country_id = c.id
WHERE ua.id = sqlc.arg(id);

-- name: AddUserAddress :one
INSERT INTO users_addresses (
    user_id, 
    address, 
    city, 
    postal_code, 
    country_id, 
    phone, 
    is_shipping, 
    is_shipping_default,
    is_billing,
    is_billing_default
)
VALUES (
    sqlc.arg(user_id), 
    sqlc.arg(address), 
    sqlc.arg(city), 
    sqlc.arg(postal_code), 
    sqlc.arg(country), 
    sqlc.arg(phone), 
    sqlc.arg(is_shipping), 
    sqlc.arg(is_shipping_default),
    sqlc.arg(is_billing),
    sqlc.arg(is_billing_default)
)
RETURNING *;

-- name: UpdateUserAddress :one
UPDATE users_addresses
SET address = sqlc.arg(address),
    city = sqlc.arg(city),
    postal_code = sqlc.arg(postal_code),
    country_id = sqlc.arg(country_id),
    phone = sqlc.arg(phone),
    is_shipping = sqlc.arg(is_shipping),
    is_shipping_default = sqlc.arg(is_shipping_default),
    is_billing = sqlc.arg(is_billing),
    is_billing_default = sqlc.arg(is_billing_default)
WHERE id = sqlc.arg(id)
RETURNING *;
-- name: DeleteUserAddress :exec
DELETE FROM users_addresses
WHERE id = sqlc.arg(id);
-- name: SetDefaultShippingAddress :exec
UPDATE users_addresses
SET is_shipping_default = TRUE
WHERE id = sqlc.arg(id)
AND user_id = sqlc.arg(user_id)
RETURNING *;
-- name: SetDefaultBillingAddress :exec
UPDATE users_addresses
SET is_billing_default = TRUE
WHERE id = sqlc.arg(id)
AND user_id = sqlc.arg(user_id)
RETURNING *;
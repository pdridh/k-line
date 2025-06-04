-- name: CreateMenuItem :one
INSERT INTO menu_items (
  name,
  description,
  price,
  requires_ticket
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetMenuItems :many
SELECT * FROM menu_items
WHERE (@search::text IS NULL OR name ILIKE '%' || @search::text || '%')
ORDER BY name
LIMIT $1
OFFSET $2;

-- name: GetItemByID :one
SELECT * FROM menu_items
WHERE id = $1 
LIMIT 1;



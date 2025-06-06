-- name: CreateOrder :one
INSERT INTO orders (
  type,
  employee_id,
  table_id
) VALUES (
  $1, $2, $3
) RETURNING id;


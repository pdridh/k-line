-- name: CreateOrder :one
INSERT INTO orders (
  type,
  employee_id,
  table_id
) VALUES (
  $1, $2, $3
) RETURNING id;

-- name: GetOrderByID :one
SELECT * FROM orders 
WHERE id = $1;


-- name: GetOrders :many
SELECT * FROM orders
WHERE status = $1 AND type = $2;


-- name: AddOrderItemsBulk :exec
INSERT INTO order_items (order_id, item_id, quantity, notes)
SELECT $1, unnest(@item_ids::int[]), unnest(@quantity::int[]), unnest(@notes::text[]);

-- name: GetOrderItemByID :one
SELECT * FROM order_items
WHERE order_id = $1 AND id = $2;

-- name: UpdateOrderItemStatus :exec
UPDATE order_items
SET status = $1
WHERE order_id = $2 AND id = $3;

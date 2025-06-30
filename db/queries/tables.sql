-- name: GetTableByID :one
SELECT * FROM tables
WHERE id = $1;


-- name: UpdateTableStatus :exec 
UPDATE tables
SET status = $1
WHERE id = $2;


-- name: GetTables :many
SELECT * FROM tables
WHERE status = $1;


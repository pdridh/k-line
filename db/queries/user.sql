-- name: CreateUser :one
INSERT INTO users (
  email,
  name,
  type,
  password
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users 
WHERE email = $1 LIMIT 1;

-- name: GetUsers :many
SELECT * FROM users
ORDER BY name
LIMIT $1
OFFSET $2;
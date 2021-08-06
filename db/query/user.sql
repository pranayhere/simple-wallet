-- name: CreateUser :one
INSERT INTO users (
    username,
    hashed_password,
    status,
    full_name,
    email
) values (
$1, $2, $3, $4, $5
) RETURNING *;

-- name: GetUser :one
SELECT * from users
where username = $1 LIMIT 1;
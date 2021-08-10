-- name: CreateCurrency :one
INSERT INTO currencies (code,
                        fraction)
VALUES ($1, $2)
ON CONFLICT (code) DO UPDATE SET fraction = $2
returning code, fraction, created_at;

-- name: GetCurrency :one
SELECT *
FROM currencies
where code = $1
LIMIT 1;
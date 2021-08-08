CREATE TABLE "currencies"
(
    "code"       varchar PRIMARY KEY,
    "fraction"   bigint    NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT 'now()'
);

-- name: CreateCurrency :one
INSERT INTO currencies (code,
                        fraction)
VALUES ($1, $2)
returning *;

-- name: GetCurrency :one
SELECT *
FROM currencies
where code = $1
LIMIT 1;
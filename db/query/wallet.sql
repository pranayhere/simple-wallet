-- name: GetWallet :one
SELECT *
FROM wallets
WHERE id = $1
LIMIT 1;

-- name: GetWalletByAddress :one
SELECT *
FROM wallets
WHERE address = $1
LIMIT 1;

-- name: GetWalletForUpdate :one
SELECT *
FROM wallets
WHERE id = $1
LIMIT 1 FOR NO KEY
    UPDATE;

-- name: GetWalletByAddressForUpdate :one
SELECT *
FROM wallets
WHERE address = $1
LIMIT 1 FOR NO KEY
    UPDATE;

-- name: CreateWallet :one
INSERT INTO wallets (name,
                     address,
                     status,
                     user_id,
                     bank_account_id,
                     balance,
                     currency)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: ListWallets :many
SELECT *
FROM wallets
WHERE user_id = $1
ORDER BY id
LIMIT $2 OFFSET $3;

-- name: AddWalletBalance :one
UPDATE wallets
SET balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING *;
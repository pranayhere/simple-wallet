-- name: CreateTransfer :one
INSERT INTO transfers (transfer_type,
                       from_wallet_id,
                       to_wallet_id,
                       amount)
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetTransfer :one
SELECT *
FROM transfers
WHERE id = $1 LIMIT 1;

-- name: ListTransfers :many
SELECT *
FROM transfers
WHERE from_wallet_id = $1
   OR to_wallet_id = $2
ORDER BY id LIMIT $3
    OFFSET $4;
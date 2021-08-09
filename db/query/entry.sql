-- name: CreateEntry :one
INSERT INTO entries (wallet_id,
                     amount,
                     entry_type,
                     balance,
                     transfer_id)
VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetEntry :one
SELECT *
FROM entries
WHERE id = $1 LIMIT 1;

-- name: ListEntries :many
SELECT *
FROM entries
WHERE wallet_id = $1
ORDER BY id LIMIT $2
    OFFSET $3;

-- name: GetEntriesByTransferID :many
SELECT *
FROM entries
WHERE transfer_id = $1;

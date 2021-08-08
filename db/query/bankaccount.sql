-- name: CreateBankAccount :one
INSERT INTO bank_accounts (account_no,
                           ifsc,
                           bank_name,
                           currency,
                           user_id,
                           status)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetBankAccount :one
SELECT *
from bank_accounts
where id = $1
LIMIT 1;

-- name: UpdateBankAccountStatus :one
UPDATE bank_accounts
set Status = $1
where id = $2
RETURNING *;

-- name: ListBankAccounts :many
SELECT *
FROM bank_accounts
WHERE user_id = $1
ORDER BY id
LIMIT $2 OFFSET $3;
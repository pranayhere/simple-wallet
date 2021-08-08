package store

import (
    "context"
    "database/sql"
    "github.com/pranayhere/simple-wallet/domains"
)

type BankAccountRepo interface {
    CreateBankAccount(ctx context.Context, arg CreateBankAccountParams) (domains.BankAccount, error)
    GetBankAccount(ctx context.Context, id int64) (domains.BankAccount, error)
    ListBankAccounts(ctx context.Context, arg ListBankAccountsParams) ([]domains.BankAccount, error)
    UpdateBankAccountStatus(ctx context.Context, arg UpdateBankAccountStatusParams) (domains.BankAccount, error)
}

type bankAccountRepository struct {
    db *sql.DB
}

func NewBankAccountRepo(client *sql.DB) BankAccountRepo {
    return &bankAccountRepository{
        db: client,
    }
}

const createBankAccount = `-- name: CreateBankAccount :one
INSERT INTO bank_accounts (account_no,
                           ifsc,
                           bank_name,
                           currency,
                           user_id,
                           status)
VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, account_no, ifsc, bank_name, status, user_id, currency, created_at, updated_at
`

type CreateBankAccountParams struct {
    AccountNo string                    `json:"account_no"`
    Ifsc      string                    `json:"ifsc"`
    BankName  string                    `json:"bank_name"`
    Currency  string                    `json:"currency"`
    UserID    int64                     `json:"user_id"`
    Status    domains.BankAccountStatus `json:"status"`
}

func (q *bankAccountRepository) CreateBankAccount(ctx context.Context, arg CreateBankAccountParams) (domains.BankAccount, error) {
    row := q.db.QueryRowContext(ctx, createBankAccount,
        arg.AccountNo,
        arg.Ifsc,
        arg.BankName,
        arg.Currency,
        arg.UserID,
        arg.Status,
    )
    var i domains.BankAccount
    err := row.Scan(
        &i.ID,
        &i.AccountNo,
        &i.Ifsc,
        &i.BankName,
        &i.Status,
        &i.UserID,
        &i.Currency,
        &i.CreatedAt,
        &i.UpdatedAt,
    )
    return i, err
}

const getBankAccount = `-- name: GetBankAccount :one
SELECT id, account_no, ifsc, bank_name, status, user_id, currency, created_at, updated_at
from bank_accounts
where id = $1 LIMIT 1
`

func (q *bankAccountRepository) GetBankAccount(ctx context.Context, id int64) (domains.BankAccount, error) {
    row := q.db.QueryRowContext(ctx, getBankAccount, id)
    var i domains.BankAccount
    err := row.Scan(
        &i.ID,
        &i.AccountNo,
        &i.Ifsc,
        &i.BankName,
        &i.Status,
        &i.UserID,
        &i.Currency,
        &i.CreatedAt,
        &i.UpdatedAt,
    )
    return i, err
}

const listBankAccounts = `-- name: ListBankAccounts :many
SELECT id, account_no, ifsc, bank_name, status, user_id, currency, created_at, updated_at
FROM bank_accounts
WHERE user_id = $1
ORDER BY id LIMIT $2
OFFSET $3
`

type ListBankAccountsParams struct {
    UserID int64 `json:"user_id"`
    Limit  int32 `json:"limit"`
    Offset int32 `json:"offset"`
}

func (q *bankAccountRepository) ListBankAccounts(ctx context.Context, arg ListBankAccountsParams) ([]domains.BankAccount, error) {
    rows, err := q.db.QueryContext(ctx, listBankAccounts, arg.UserID, arg.Limit, arg.Offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    items := []domains.BankAccount{}
    for rows.Next() {
        var i domains.BankAccount
        if err := rows.Scan(
            &i.ID,
            &i.AccountNo,
            &i.Ifsc,
            &i.BankName,
            &i.Status,
            &i.UserID,
            &i.Currency,
            &i.CreatedAt,
            &i.UpdatedAt,
        ); err != nil {
            return nil, err
        }
        items = append(items, i)
    }
    if err := rows.Close(); err != nil {
        return nil, err
    }
    if err := rows.Err(); err != nil {
        return nil, err
    }
    return items, nil
}

const updateBankAccountStatus = `-- name: UpdateBankAccountStatus :one
UPDATE bank_accounts
set Status = $1
where id = $2
RETURNING id, account_no, ifsc, bank_name, status, user_id, currency, created_at, updated_at
`

type UpdateBankAccountStatusParams struct {
    Status domains.BankAccountStatus `json:"status"`
    ID     int64                     `json:"id"`
}

func (q *bankAccountRepository) UpdateBankAccountStatus(ctx context.Context, arg UpdateBankAccountStatusParams) (domains.BankAccount, error) {
    row := q.db.QueryRowContext(ctx, updateBankAccountStatus, arg.Status, arg.ID)
    var i domains.BankAccount
    err := row.Scan(
        &i.ID,
        &i.AccountNo,
        &i.Ifsc,
        &i.BankName,
        &i.Status,
        &i.UserID,
        &i.Currency,
        &i.CreatedAt,
        &i.UpdatedAt,
    )
    return i, err
}

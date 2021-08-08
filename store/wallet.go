package store

import (
    "context"
    "database/sql"
    "github.com/pranayhere/simple-wallet/domains"
)

type WalletRepo interface {
    AddWalletBalance(ctx context.Context, arg AddWalletBalanceParams) (domains.Wallet, error)
    CreateWallet(ctx context.Context, arg CreateWalletParams) (domains.Wallet, error)
    GetWallet(ctx context.Context, id int64) (domains.Wallet, error)
    GetWalletForUpdate(ctx context.Context, id int64) (domains.Wallet, error)
    GetWalletByAddress(ctx context.Context, address string) (domains.Wallet, error)
    GetWalletByAddressForUpdate(ctx context.Context, address string) (domains.Wallet, error)
    ListWallets(ctx context.Context, arg ListWalletsParams) ([]domains.Wallet, error)
}

type walletRepository struct {
    db *sql.DB
}

func NewWalletRepo(client *sql.DB) WalletRepo {
    return &walletRepository{
        db: client,
    }
}

const addWalletBalance = `-- name: AddWalletBalance :one
UPDATE wallets
SET balance = balance + $1
WHERE id = $2 RETURNING id, name, address, status, user_id, bank_account_id, balance, currency, created_at, updated_at
`

type AddWalletBalanceParams struct {
    Amount int64 `json:"amount"`
    ID     int64 `json:"id"`
}

func (q *walletRepository) AddWalletBalance(ctx context.Context, arg AddWalletBalanceParams) (domains.Wallet, error) {
    row := q.db.QueryRowContext(ctx, addWalletBalance, arg.Amount, arg.ID)
    var i domains.Wallet
    err := row.Scan(
        &i.ID,
        &i.Name,
        &i.Address,
        &i.Status,
        &i.UserID,
        &i.BankAccountID,
        &i.Balance,
        &i.Currency,
        &i.CreatedAt,
        &i.UpdatedAt,
    )
    return i, err
}

const createWallet = `-- name: CreateWallet :one
INSERT INTO wallets (name,
                     address,
                     status,
                     user_id,
                     bank_account_id,
                     balance,
                     currency)
VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, name, address, status, user_id, bank_account_id, balance, currency, created_at, updated_at
`

type CreateWalletParams struct {
    Name          string               `json:"name"`
    Address       string               `json:"address"`
    Status        domains.WalletStatus `json:"status"`
    UserID        int64                `json:"user_id"`
    BankAccountID int64                `json:"bank_account_id"`
    Balance       int64                `json:"balance"`
    Currency      string               `json:"currency"`
}

func (q *walletRepository) CreateWallet(ctx context.Context, arg CreateWalletParams) (domains.Wallet, error) {
    row := q.db.QueryRowContext(ctx, createWallet,
        arg.Name,
        arg.Address,
        arg.Status,
        arg.UserID,
        arg.BankAccountID,
        arg.Balance,
        arg.Currency,
    )
    var i domains.Wallet
    err := row.Scan(
        &i.ID,
        &i.Name,
        &i.Address,
        &i.Status,
        &i.UserID,
        &i.BankAccountID,
        &i.Balance,
        &i.Currency,
        &i.CreatedAt,
        &i.UpdatedAt,
    )
    return i, err
}

const getWallet = `-- name: GetWallet :one
SELECT id, name, address, status, user_id, bank_account_id, balance, currency, created_at, updated_at
FROM wallets
WHERE id = $1 LIMIT 1
`

func (q *walletRepository) GetWallet(ctx context.Context, id int64) (domains.Wallet, error) {
    row := q.db.QueryRowContext(ctx, getWallet, id)
    var i domains.Wallet
    err := row.Scan(
        &i.ID,
        &i.Name,
        &i.Address,
        &i.Status,
        &i.UserID,
        &i.BankAccountID,
        &i.Balance,
        &i.Currency,
        &i.CreatedAt,
        &i.UpdatedAt,
    )
    return i, err
}

const getWalletByAddress = `-- name: GetWalletByAddress :one
SELECT id, name, address, status, user_id, bank_account_id, balance, currency, created_at, updated_at
FROM wallets
WHERE address = $1 LIMIT 1
`

func (q *walletRepository) GetWalletByAddress(ctx context.Context, address string) (domains.Wallet, error) {
    row := q.db.QueryRowContext(ctx, getWalletByAddress, address)
    var i domains.Wallet
    err := row.Scan(
        &i.ID,
        &i.Name,
        &i.Address,
        &i.Status,
        &i.UserID,
        &i.BankAccountID,
        &i.Balance,
        &i.Currency,
        &i.CreatedAt,
        &i.UpdatedAt,
    )
    return i, err
}

const getWalletByAddressForUpdate = `-- name: GetWalletByAddressForUpdate :one
SELECT id, name, address, status, user_id, bank_account_id, balance, currency, created_at, updated_at
FROM wallets
WHERE address = $1 LIMIT 1
FOR NO KEY
UPDATE
`

func (q *walletRepository) GetWalletByAddressForUpdate(ctx context.Context, address string) (domains.Wallet, error) {
    row := q.db.QueryRowContext(ctx, getWalletByAddressForUpdate, address)
    var i domains.Wallet
    err := row.Scan(
        &i.ID,
        &i.Name,
        &i.Address,
        &i.Status,
        &i.UserID,
        &i.BankAccountID,
        &i.Balance,
        &i.Currency,
        &i.CreatedAt,
        &i.UpdatedAt,
    )
    return i, err
}

const getWalletForUpdate = `-- name: GetWalletForUpdate :one
SELECT id, name, address, status, user_id, bank_account_id, balance, currency, created_at, updated_at
FROM wallets
WHERE id = $1 LIMIT 1
FOR NO KEY
UPDATE
`

func (q *walletRepository) GetWalletForUpdate(ctx context.Context, id int64) (domains.Wallet, error) {
    row := q.db.QueryRowContext(ctx, getWalletForUpdate, id)
    var i domains.Wallet
    err := row.Scan(
        &i.ID,
        &i.Name,
        &i.Address,
        &i.Status,
        &i.UserID,
        &i.BankAccountID,
        &i.Balance,
        &i.Currency,
        &i.CreatedAt,
        &i.UpdatedAt,
    )
    return i, err
}

const listWallets = `-- name: ListWallets :many
SELECT id, name, address, status, user_id, bank_account_id, balance, currency, created_at, updated_at
FROM wallets
WHERE user_id = $1
ORDER BY id LIMIT $2
OFFSET $3
`

type ListWalletsParams struct {
    UserID int64 `json:"user_id"`
    Limit  int32 `json:"limit"`
    Offset int32 `json:"offset"`
}

func (q *walletRepository) ListWallets(ctx context.Context, arg ListWalletsParams) ([]domains.Wallet, error) {
    rows, err := q.db.QueryContext(ctx, listWallets, arg.UserID, arg.Limit, arg.Offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    items := []domains.Wallet{}
    for rows.Next() {
        var i domains.Wallet
        if err := rows.Scan(
            &i.ID,
            &i.Name,
            &i.Address,
            &i.Status,
            &i.UserID,
            &i.BankAccountID,
            &i.Balance,
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

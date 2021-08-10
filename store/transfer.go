package store

import (
    "context"
    "database/sql"
    "github.com/pranayhere/simple-wallet/domains"
)

type TransferRepo interface {
    CreateTransfer(ctx context.Context, arg CreateTransferParams) (domains.Transfer, error)
    GetTransfer(ctx context.Context, id int64) (domains.Transfer, error)
    ListTransfers(ctx context.Context, arg ListTransfersParams) ([]domains.Transfer, error)
}

type transferRepository struct {
    db *sql.DB
}

func NewTransferRepo(client *sql.DB) TransferRepo {
    return &transferRepository{
        db: client,
    }
}

const createTransfer = `-- name: CreateTransfer :one
INSERT INTO transfers (transfer_type,
                       from_wallet_id,
                       to_wallet_id,
                       amount)
VALUES ($1, $2, $3, $4) RETURNING id, transfer_type, from_wallet_id, to_wallet_id, amount, created_at
`

type CreateTransferParams struct {
    TransferType domains.TransferType `json:"transfer_type"`
    FromWalletID int64        `json:"from_wallet_id"`
    ToWalletID   int64        `json:"to_wallet_id"`
    Amount       int64                `json:"amount"`
}

func (q *transferRepository) CreateTransfer(ctx context.Context, arg CreateTransferParams) (domains.Transfer, error) {
    row := q.db.QueryRowContext(ctx, createTransfer,
        arg.TransferType,
        arg.FromWalletID,
        arg.ToWalletID,
        arg.Amount,
    )
    var i domains.Transfer
    err := row.Scan(
        &i.ID,
        &i.TransferType,
        &i.FromWalletID,
        &i.ToWalletID,
        &i.Amount,
        &i.CreatedAt,
    )
    return i, err
}

const getTransfer = `-- name: GetTransfer :one
SELECT id, transfer_type, from_wallet_id, to_wallet_id, amount, created_at
FROM transfers
WHERE id = $1 LIMIT 1
`

func (q *transferRepository) GetTransfer(ctx context.Context, id int64) (domains.Transfer, error) {
    row := q.db.QueryRowContext(ctx, getTransfer, id)
    var i domains.Transfer
    err := row.Scan(
        &i.ID,
        &i.TransferType,
        &i.FromWalletID,
        &i.ToWalletID,
        &i.Amount,
        &i.CreatedAt,
    )
    return i, err
}

const listTransfers = `-- name: ListTransfers :many
SELECT id, transfer_type, from_wallet_id, to_wallet_id, amount, created_at
FROM transfers
WHERE from_wallet_id = $1
OR to_wallet_id = $2
ORDER BY id LIMIT $3
OFFSET $4
`

type ListTransfersParams struct {
    FromWalletID int64 `json:"from_wallet_id"`
    ToWalletID   int64 `json:"to_wallet_id"`
    Limit        int32         `json:"limit"`
    Offset       int32         `json:"offset"`
}

func (q *transferRepository) ListTransfers(ctx context.Context, arg ListTransfersParams) ([]domains.Transfer, error) {
    rows, err := q.db.QueryContext(ctx, listTransfers,
        arg.FromWalletID,
        arg.ToWalletID,
        arg.Limit,
        arg.Offset,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    items := []domains.Transfer{}
    for rows.Next() {
        var i domains.Transfer
        if err := rows.Scan(
            &i.ID,
            &i.TransferType,
            &i.FromWalletID,
            &i.ToWalletID,
            &i.Amount,
            &i.CreatedAt,
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

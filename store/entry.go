package store

import (
    "context"
    "database/sql"
    "github.com/pranayhere/simple-wallet/domains"
)

type EntryRepo interface {
    CreateEntry(ctx context.Context, arg CreateEntryParams) (domains.Entry, error)
    GetEntry(ctx context.Context, id int64) (domains.Entry, error)
    ListEntries(ctx context.Context, arg ListEntriesParams) ([]domains.Entry, error)
}

type entryRepository struct {
    db *sql.DB
}

func NewEntryRepo(client *sql.DB) EntryRepo {
    return &entryRepository{
        db: client,
    }
}

const createEntry = `-- name: CreateEntry :one
INSERT INTO entries (wallet_id,
                     amount,
                     transfer_id)
VALUES ($1, $2, $3) RETURNING id, wallet_id, amount, transfer_id, created_at
`

type CreateEntryParams struct {
    WalletID   int64             `json:"wallet_id"`
    Amount     int64             `json:"amount"`
    TransferID int64             `json:"transfer_id"`
}

func (q *entryRepository) CreateEntry(ctx context.Context, arg CreateEntryParams) (domains.Entry, error) {
    row := q.db.QueryRowContext(ctx, createEntry,
        arg.WalletID,
        arg.Amount,
        arg.TransferID,
    )
    var i domains.Entry
    err := row.Scan(
        &i.ID,
        &i.WalletID,
        &i.Amount,
        &i.TransferID,
        &i.CreatedAt,
    )
    return i, err

}

const getEntriesByTransferID = `-- name: GetEntriesByTransferID :many
SELECT id, wallet_id, amount, transfer_id, created_at
FROM entries
WHERE transfer_id = $1
`

func (q *entryRepository) GetEntriesByTransferID(ctx context.Context, transferID int64) ([]domains.Entry, error) {
    rows, err := q.db.QueryContext(ctx, getEntriesByTransferID, transferID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    items := []domains.Entry{}
    for rows.Next() {
        var i domains.Entry
        if err := rows.Scan(
            &i.ID,
            &i.WalletID,
            &i.Amount,
            &i.TransferID,
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

const getEntry = `-- name: GetEntry :one
SELECT id, wallet_id, amount, transfer_id, created_at
FROM entries
WHERE id = $1 LIMIT 1
`

func (q *entryRepository) GetEntry(ctx context.Context, id int64) (domains.Entry, error) {
    row := q.db.QueryRowContext(ctx, getEntry, id)
    var i domains.Entry
    err := row.Scan(
        &i.ID,
        &i.WalletID,
        &i.Amount,
        &i.TransferID,
        &i.CreatedAt,
    )
    return i, err
}

const listEntries = `-- name: ListEntries :many
SELECT id, wallet_id, amount, transfer_id, created_at
FROM entries
WHERE wallet_id = $1
ORDER BY id LIMIT $2
OFFSET $3
`

type ListEntriesParams struct {
    WalletID int64 `json:"wallet_id"`
    Limit    int32 `json:"limit"`
    Offset   int32 `json:"offset"`
}

func (q *entryRepository) ListEntries(ctx context.Context, arg ListEntriesParams) ([]domains.Entry, error) {
    rows, err := q.db.QueryContext(ctx, listEntries, arg.WalletID, arg.Limit, arg.Offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    items := []domains.Entry{}
    for rows.Next() {
        var i domains.Entry
        if err := rows.Scan(
            &i.ID,
            &i.WalletID,
            &i.Amount,
            &i.TransferID,
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

package store

import (
    "context"
    "database/sql"
    "github.com/pranayhere/simple-wallet/domain"
)

type PaymentRequestRepo interface {
    CreatePaymentRequest(ctx context.Context, arg CreatePaymentRequestParams) (domain.PaymentRequest, error)
    ListPaymentRequests(ctx context.Context, arg ListPaymentRequestsParams) ([]domain.PaymentRequest, error)
    UpdatePaymentRequest(ctx context.Context, arg UpdatePaymentRequestParams) (domain.PaymentRequest, error)
    GetPaymentRequest(ctx context.Context, id int64) (domain.PaymentRequest, error)
}

type paymentRequestRepository struct {
    db *sql.DB
}

func NewPaymentRequestRepo(client *sql.DB) PaymentRequestRepo {
    return &paymentRequestRepository{
        db: client,
    }
}

const createPaymentRequest = `-- name: CreatePaymentRequest :one
INSERT INTO payment_requests (
                       from_wallet_id,
                       to_wallet_id,
                       amount,
                       status)
VALUES ($1, $2, $3, $4) RETURNING id, from_wallet_id, to_wallet_id, amount, status, created_at
`

type CreatePaymentRequestParams struct {
    FromWalletID int64                       `json:"from_wallet_id"`
    ToWalletID   int64                       `json:"to_wallet_id"`
    Amount       int64                       `json:"amount"`
    Status       domain.PaymentRequestStatus `json:"status"`
}

func (q *paymentRequestRepository) CreatePaymentRequest(ctx context.Context, arg CreatePaymentRequestParams) (domain.PaymentRequest, error) {
    row := q.db.QueryRowContext(ctx, createPaymentRequest,
        arg.FromWalletID,
        arg.ToWalletID,
        arg.Amount,
        arg.Status,
    )
    var i domain.PaymentRequest
    err := row.Scan(
        &i.ID,
        &i.FromWalletID,
        &i.ToWalletID,
        &i.Amount,
        &i.Status,
        &i.CreatedAt,
    )
    return i, err
}

const listPaymentRequests = `-- name: ListPaymentRequests :many
SELECT id, from_wallet_id, to_wallet_id, amount, status, created_at
FROM payment_requests
WHERE from_wallet_id = $1
ORDER BY id LIMIT $2
OFFSET $3
`

type ListPaymentRequestsParams struct {
    FromWalletID int64 `json:"from_wallet_id"`
    Limit        int32 `json:"limit"`
    Offset       int32 `json:"offset"`
}

func (q *paymentRequestRepository) ListPaymentRequests(ctx context.Context, arg ListPaymentRequestsParams) ([]domain.PaymentRequest, error) {
    rows, err := q.db.QueryContext(ctx, listPaymentRequests, arg.FromWalletID, arg.Limit, arg.Offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    items := []domain.PaymentRequest{}
    for rows.Next() {
        var i domain.PaymentRequest
        if err := rows.Scan(
            &i.ID,
            &i.FromWalletID,
            &i.ToWalletID,
            &i.Amount,
            &i.Status,
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

const updatePaymentRequest = `-- name: UpdatePaymentRequest :one
UPDATE payment_requests
set Status = $1
where id = $2
RETURNING id, from_wallet_id, to_wallet_id, amount, status, created_at
`

type UpdatePaymentRequestParams struct {
    Status domain.PaymentRequestStatus `json:"status"`
    ID     int64                       `json:"id"`
}

func (q *paymentRequestRepository) UpdatePaymentRequest(ctx context.Context, arg UpdatePaymentRequestParams) (domain.PaymentRequest, error) {
    row := q.db.QueryRowContext(ctx, updatePaymentRequest, arg.Status, arg.ID)
    var i domain.PaymentRequest
    err := row.Scan(
        &i.ID,
        &i.FromWalletID,
        &i.ToWalletID,
        &i.Amount,
        &i.Status,
        &i.CreatedAt,
    )
    return i, err
}

const getPaymentRequest = `-- name: GetPaymentRequest :one
SELECT id, from_wallet_id, to_wallet_id, amount, status, created_at
from payment_requests
where id = $1 LIMIT 1
`

func (q *paymentRequestRepository) GetPaymentRequest(ctx context.Context, id int64) (domain.PaymentRequest, error) {
    row := q.db.QueryRowContext(ctx, getPaymentRequest, id)
    var i domain.PaymentRequest
    err := row.Scan(
        &i.ID,
        &i.FromWalletID,
        &i.ToWalletID,
        &i.Amount,
        &i.Status,
        &i.CreatedAt,
    )
    return i, err
}

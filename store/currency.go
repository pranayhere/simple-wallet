package store

import (
    "context"
    "database/sql"
    "github.com/pranayhere/simple-wallet/domain"
)

type CurrencyRepo interface {
    CreateCurrency(ctx context.Context, arg CreateCurrencyParams) (domain.Currency, error)
    GetCurrency(ctx context.Context, code string) (domain.Currency, error)
}

type currencyRepository struct {
    db *sql.DB
}

func NewCurrencyRepo(client *sql.DB) CurrencyRepo {
    return &currencyRepository{
        db: client,
    }
}

const createCurrency = `-- name: CreateCurrency :one
INSERT INTO currencies (code,
                        fraction)
VALUES ($1, $2)
ON CONFLICT (code) DO UPDATE SET fraction = $2
returning code, fraction, created_at
`

type CreateCurrencyParams struct {
    Code     string `json:"code"`
    Fraction int64  `json:"fraction"`
}

func (q *currencyRepository) CreateCurrency(ctx context.Context, arg CreateCurrencyParams) (domain.Currency, error) {
    row := q.db.QueryRowContext(ctx, createCurrency, arg.Code, arg.Fraction)
    var i domain.Currency
    err := row.Scan(&i.Code, &i.Fraction, &i.CreatedAt)
    return i, err
}

const getCurrency = `-- name: GetCurrency :one
SELECT code, fraction, created_at FROM currencies
where code = $1 LIMIT 1
`

func (q *currencyRepository) GetCurrency(ctx context.Context, code string) (domain.Currency, error) {
    row := q.db.QueryRowContext(ctx, getCurrency, code)
    var i domain.Currency
    err := row.Scan(&i.Code, &i.Fraction, &i.CreatedAt)
    return i, err
}

package store

import (
    "context"
    "database/sql"
    "fmt"
    "github.com/pranayhere/simple-wallet/domains"
    "strings"
)

type BankAccountRepo interface {
    CreateBankAccount(ctx context.Context, arg CreateBankAccountParams) (domains.BankAccount, error)
    GetBankAccount(ctx context.Context, id int64) (domains.BankAccount, error)
    ListBankAccounts(ctx context.Context, arg ListBankAccountsParams) ([]domains.BankAccount, error)
    UpdateBankAccountStatus(ctx context.Context, arg UpdateBankAccountStatusParams) (domains.BankAccount, error)
    CreateBankAccountWithWallet(ctx context.Context, arg CreateBankAccountWithWalletParams) (BankAccountWithWalletResult, error)
    BankAccountVerificationSuccess(ctx context.Context, arg BankAccountVerificationParams) (BankAccountVerificationResult, error)
    BankAccountVerificationFailed(ctx context.Context, arg BankAccountVerificationParams) (BankAccountVerificationResult, error)
}

type bankAccountRepository struct {
    db         *sql.DB
    walletRepo WalletRepo
    userRepo   UserRepo
}

func NewBankAccountRepo(db *sql.DB, walletRepo WalletRepo, userRepo UserRepo) BankAccountRepo {
    return &bankAccountRepository{
        db:         db,
        walletRepo: walletRepo,
        userRepo:   userRepo,
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

type CreateBankAccountWithWalletParams struct {
    AccountNo string `json:"account_no"`
    Ifsc      string `json:"ifsc"`
    BankName  string `json:"bank_name"`
    Currency  string `json:"currency"`
    UserID    int64  `json:"user_id"`
}

type BankAccountWithWalletResult struct {
    BankAccount domains.BankAccount `json:"bank_account"`
    Wallet      domains.Wallet      `json:"wallet"`
}

func (q *bankAccountRepository) CreateBankAccountWithWallet(ctx context.Context, arg CreateBankAccountWithWalletParams) (BankAccountWithWalletResult, error) {
    var result BankAccountWithWalletResult

    err := WithTransaction(q.db, func(tx Transaction) error {
        var err error

        user, err := q.userRepo.GetUser(ctx, arg.UserID)

        if err != nil {
            return err
        }

        result.BankAccount, err = q.CreateBankAccount(ctx, CreateBankAccountParams{
            UserID:    user.ID,
            AccountNo: arg.AccountNo,
            Ifsc:      arg.Ifsc,
            BankName:  arg.BankName,
            Currency:  arg.Currency,
            Status:    domains.BankAccountStatusINVERIFICATION,
        })

        if err != nil {
            return err
        }

        walletAddress := strings.Split(user.Email, "@")[0]
        walletAddress = fmt.Sprintf("%s@my.wallet", walletAddress)

        result.Wallet, err = q.walletRepo.CreateWallet(ctx, CreateWalletParams{
            UserID:        user.ID,
            Currency:      arg.Currency,
            Balance:       0,
            Address:       walletAddress,
            BankAccountID: result.BankAccount.ID,
            Status:        domains.WalletStatusINACTIVE,
        })

        if err != nil {
            return err
        }

        return nil
    })

    return result, err
}

type BankAccountVerificationParams struct {
    BankAccountID int64 `json:"bank_account_id"`
}

type BankAccountVerificationResult struct {
    BankAccount domains.BankAccount `json:"bank_account"`
    Wallet      domains.Wallet      `json:"wallet"`
}

func (q *bankAccountRepository) BankAccountVerificationSuccess(ctx context.Context, arg BankAccountVerificationParams) (BankAccountVerificationResult, error) {
    var result BankAccountVerificationResult

    err := WithTransaction(q.db, func(tx Transaction) error {
        var err error

        result.BankAccount, err = q.UpdateBankAccountStatus(ctx, UpdateBankAccountStatusParams{
            ID:     arg.BankAccountID,
            Status: domains.BankAccountStatusVERIFIED,
        })
        if err != nil {
            return err
        }

        wallet, err := q.walletRepo.GetWalletByBankAccountIDForUpdate(ctx, result.BankAccount.ID)
        if err != nil {
            return err
        }

        result.Wallet, err = q.walletRepo.UpdateWalletStatus(ctx, UpdateWalletStatusParams{
            ID:     wallet.ID,
            Status: domains.WalletStatusACTIVE,
        })
        if err != nil {
            return err
        }

        return nil
    })

    return result, err
}

func (q *bankAccountRepository) BankAccountVerificationFailed(ctx context.Context, arg BankAccountVerificationParams) (BankAccountVerificationResult, error) {
    var result BankAccountVerificationResult

    var err error

    result.BankAccount, err = q.UpdateBankAccountStatus(ctx, UpdateBankAccountStatusParams{
        ID:     arg.BankAccountID,
        Status: domains.BankAccountStatusVERIFICATIONFAILED,
    })
    if err != nil {
        return result, err
    }

    return result, err
}

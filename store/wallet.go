package store

import (
    "context"
    "database/sql"
    "fmt"
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
    UpdateWalletStatus(ctx context.Context, arg UpdateWalletStatusParams) (domains.Wallet, error)
    GetWalletByBankAccountID(ctx context.Context, bankAccountID int64) (domains.Wallet, error)
    GetWalletByBankAccountIDForUpdate(ctx context.Context, bankAccountID int64) (domains.Wallet, error)
    DepositToWallet(ctx context.Context, arg DepositeToWalletParams) (WalletTransferResult, error)
    WithdrawFromWallet(ctx context.Context, arg WithdrawFromWalletParams) (WalletTransferResult, error)
    SendMoney(ctx context.Context, arg SendMoneyParams) (WalletTransferResult, error)
}

type walletRepository struct {
    db           *sql.DB
    transferRepo TransferRepo
    entryRepo    EntryRepo
}

func NewWalletRepo(client *sql.DB, transferRepo TransferRepo, entryRepo EntryRepo) WalletRepo {
    return &walletRepository{
        db:           client,
        transferRepo: transferRepo,
        entryRepo:    entryRepo,
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

const updateWalletStatus = `-- name: UpdateWalletStatus :one
UPDATE wallets
set Status = $1
where id = $2
RETURNING id, name, address, status, user_id, bank_account_id, balance, currency, created_at, updated_at
`

type UpdateWalletStatusParams struct {
    Status domains.WalletStatus `json:"status"`
    ID     int64                `json:"id"`
}

func (q *walletRepository) UpdateWalletStatus(ctx context.Context, arg UpdateWalletStatusParams) (domains.Wallet, error) {
    row := q.db.QueryRowContext(ctx, updateWalletStatus, arg.Status, arg.ID)
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

const getWalletByBankAccountID = `-- name: GetWalletByBankAccountID :one
SELECT id, name, address, status, user_id, bank_account_id, balance, currency, created_at, updated_at
FROM wallets
WHERE bank_account_id = $1 LIMIT 1
`

func (q *walletRepository) GetWalletByBankAccountID(ctx context.Context, bankAccountID int64) (domains.Wallet, error) {
    row := q.db.QueryRowContext(ctx, getWalletByBankAccountID, bankAccountID)
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

const getWalletByBankAccountIDForUpdate = `-- name: GetWalletByBankAccountIDForUpdate :one
SELECT id, name, address, status, user_id, bank_account_id, balance, currency, created_at, updated_at
FROM wallets
WHERE bank_account_id = $1 LIMIT 1
FOR NO KEY
UPDATE
`

func (q *walletRepository) GetWalletByBankAccountIDForUpdate(ctx context.Context, bankAccountID int64) (domains.Wallet, error) {
    row := q.db.QueryRowContext(ctx, getWalletByBankAccountIDForUpdate, bankAccountID)
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

type DepositeToWalletParams struct {
    WalletID int64 `json:"wallet_id"`
    Amount   int64 `json:"amount"`
}

type WalletTransferResult struct {
    Wallet    domains.Wallet   `json:"wallet"`
    FromEntry domains.Entry    `json:"from_entry"`
    ToEntry   domains.Entry    `json:"to_entry"`
    Transfer  domains.Transfer `json:"transfer"`
}

// DepositToWallet transfer money from linked bank account to the wallet
func (q *walletRepository) DepositToWallet(ctx context.Context, arg DepositeToWalletParams) (WalletTransferResult, error) {
    var res WalletTransferResult
    err := ExecTx(q.db, func(tx Tx) error {
        var err error

        wallet, err := q.GetWalletForUpdate(ctx, arg.WalletID)
        if err != nil {
            return err
        }

        if wallet.Status != domains.WalletStatusACTIVE {
            return fmt.Errorf("inactive wallet")
        }

        res.Wallet, err = q.AddWalletBalance(ctx, AddWalletBalanceParams{
            ID:     wallet.ID,
            Amount: arg.Amount,
        })

        res.Transfer, err = q.transferRepo.CreateTransfer(ctx, CreateTransferParams{
            TransferType: domains.TransferTypeDEPOSITTOWALLET,
            Amount:       arg.Amount,
            ToWalletID:   sql.NullInt64{wallet.ID, true},
        })

        if err != nil {
            return err
        }

        res.ToEntry, err = q.entryRepo.CreateEntry(ctx, CreateEntryParams{
            WalletID:   wallet.ID,
            TransferID: res.Transfer.ID,
            Amount:     arg.Amount,
        })

        if err != nil {
            return err
        }

        return nil
    })

    return res, err
}

type WithdrawFromWalletParams struct {
    WalletID int64 `json:"wallet_id"`
    Amount   int64 `json:"amount"`
    UserId   int64 `json:"user_id"`
}

// WithdrawFromWallet transfer money from wallet to the linked bank account
func (q *walletRepository) WithdrawFromWallet(ctx context.Context, arg WithdrawFromWalletParams) (WalletTransferResult, error) {
    var res WalletTransferResult
    err := ExecTx(q.db, func(tx Tx) error {
        var err error

        wallet, err := q.GetWalletForUpdate(ctx, arg.WalletID)
        if err != nil {
            return err
        }

        if wallet.Status != domains.WalletStatusACTIVE {
            return fmt.Errorf("inactive wallet")
        }

        if !wallet.IsBalanceSufficient(arg.Amount) {
            return fmt.Errorf("insufficient wallet balance")
        }

        res.Wallet, err = q.AddWalletBalance(ctx, AddWalletBalanceParams{
            ID:     wallet.ID,
            Amount: arg.Amount * -1,
        })

        res.Transfer, err = q.transferRepo.CreateTransfer(ctx, CreateTransferParams{
            TransferType: domains.TransferTypeWITHDRAWFROMWALLET,
            Amount:       arg.Amount,
            ToWalletID:   sql.NullInt64{wallet.ID, true},
        })

        if err != nil {
            return err
        }

        res.ToEntry, err = q.entryRepo.CreateEntry(ctx, CreateEntryParams{
            WalletID:   wallet.ID,
            TransferID: res.Transfer.ID,
            Amount:     arg.Amount * -1,
        })

        if err != nil {
            return err
        }

        return nil
    })

    return res, err
}

type SendMoneyParams struct {
    FromWalletAddress string `json:"from_account_address"`
    ToWalletAddress   string `json:"to_account_address""`
    Amount             int64  `json:"amount"`
}

func (q *walletRepository) SendMoney(ctx context.Context, arg SendMoneyParams) (WalletTransferResult, error) {
    var res WalletTransferResult

    err := ExecTx(q.db, func(tx Tx) error {
        var err error

        fromWallet, err := q.GetWalletByAddressForUpdate(ctx, arg.FromWalletAddress)
        if err != nil {
            return err
        }

        if fromWallet.Status != domains.WalletStatusACTIVE {
            return fmt.Errorf("inactive wallet")
        }

        if !fromWallet.IsBalanceSufficient(arg.Amount) {
            return fmt.Errorf("insufficient wallet balance")
        }

        toWallet, err := q.GetWalletByAddressForUpdate(ctx, arg.ToWalletAddress)
        if err != nil {
            return err
        }

        if toWallet.Status != domains.WalletStatusACTIVE {
            return fmt.Errorf("inactive wallet")
        }

        res.Transfer, err = q.transferRepo.CreateTransfer(ctx, CreateTransferParams{
            TransferType: domains.TransferTypeSENDMONEY,
            FromWalletID: sql.NullInt64{fromWallet.ID, true},
            ToWalletID:   sql.NullInt64{toWallet.ID, true},
            Amount:       arg.Amount,
        })

        if err != nil {
            return err
        }

        res.FromEntry, err = q.entryRepo.CreateEntry(ctx, CreateEntryParams{
            WalletID: fromWallet.ID,
            Amount:    arg.Amount * -1,
            TransferID: res.Transfer.ID,
        })

        if err != nil {
            return err
        }

        res.ToEntry, err = q.entryRepo.CreateEntry(ctx, CreateEntryParams{
            WalletID: fromWallet.ID,
            Amount:    arg.Amount,
            TransferID: res.Transfer.ID,
        })

        if err != nil {
            return err
        }

        if fromWallet.ID < toWallet.ID {
            fromWallet, toWallet, err = addMoney(ctx, q, fromWallet.ID, -arg.Amount, toWallet.ID, arg.Amount)
        } else {
            toWallet, fromWallet, err = addMoney(ctx, q, toWallet.ID, arg.Amount, fromWallet.ID, -arg.Amount)
        }

        res.Wallet = fromWallet
        return err
    })

    return res, err
}

func addMoney(ctx context.Context, q *walletRepository, walletID1 int64, amount1 int64, walletID2 int64, amount2 int64, )(wallet1 domains.Wallet, wallet2 domains.Wallet, err error) {
    wallet1, err = q.AddWalletBalance(ctx, AddWalletBalanceParams{
        ID:     walletID1,
        Amount: amount1,
    })

    if err != nil {
        return
    }

    wallet2, err = q.AddWalletBalance(ctx, AddWalletBalanceParams{
        ID:     walletID2,
        Amount: amount2,
    })

    if err != nil {
        return
    }

    return
}
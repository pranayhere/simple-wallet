package dto

import (
    "github.com/pranayhere/simple-wallet/domain"
    "github.com/pranayhere/simple-wallet/store"
    "time"
)

type SendMoneyDto struct {
    FromWalletAddress string `json:"from_account_address" validate:"required"`
    ToWalletAddress   string `json:"to_account_address" validate:"required"`
    Amount            int64  `json:"amount" validate:"required,gt=0"`
}

type WalletTransferResultDto struct {
    Wallet    domain.Wallet   `json:"wallet" validate:"required"`
    FromEntry domain.Entry    `json:"from_entry" validate:"required"`
    ToEntry   domain.Entry    `json:"to_entry" validate:"required"`
    Transfer  domain.Transfer `json:"transfer" validate:"required"`
}

type DepositDto struct {
    WalletID int64 `json:"wallet_id" validate:"required"`
    Amount   int64 `json:"amount" validate:"required,gt=0"`
    UserID   int64 `json:"user_id"`
}

type WithdrawDto struct {
    WalletID int64 `json:"wallet_id" validate:"required"`
    Amount   int64 `json:"amount" validate:"required,gt=0"`
    UserId   int64 `json:"user_id"`
}

type WalletDto struct {
    ID                   int64               `json:"id" validate:"required"`
    Address              string              `json:"address" validate:"required"`
    Status               domain.WalletStatus `json:"status" validate:"required"`
    UserID               int64               `json:"user_id" validate:"required"`
    BankAccountID        int64               `json:"bank_account_id" validate:"required"`
    OrganizationWalletID int64               `json:"organization_wallet_id" validate:"required"`
    Balance              int64               `json:"balance" validate:"required"`
    Currency             string              `json:"currency" validate:"required"`
    CreatedAt            time.Time           `json:"created_at" validate:"required"`
    UpdatedAt            time.Time           `json:"updated_at" validate:"required"`
}

func NewWalletTransferDto(wtr store.WalletTransferResult) WalletTransferResultDto {
    return WalletTransferResultDto{
        Wallet:    wtr.Wallet,
        FromEntry: wtr.FromEntry,
        ToEntry:   wtr.ToEntry,
        Transfer:  wtr.Transfer,
    }
}

func NewWalletDto(wallet domain.Wallet) WalletDto {
    return WalletDto{
        ID:                   wallet.ID,
        Address:              wallet.Address,
        Status:               wallet.Status,
        UserID:               wallet.UserID,
        BankAccountID:        wallet.BankAccountID,
        OrganizationWalletID: wallet.OrganizationWalletID,
        Balance:              wallet.Balance,
        Currency:             wallet.Currency,
        CreatedAt:            wallet.CreatedAt,
        UpdatedAt:            wallet.UpdatedAt,
    }
}

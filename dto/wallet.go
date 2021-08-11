package dto

import (
    "github.com/pranayhere/simple-wallet/domain"
    "github.com/pranayhere/simple-wallet/store"
    "time"
)

type SendMoneyDto struct {
    FromWalletAddress string `json:"from_account_address"`
    ToWalletAddress   string `json:"to_account_address""`
    Amount            int64  `json:"amount"`
}

type WalletTransferResultDto struct {
    Wallet    domain.Wallet   `json:"wallet"`
    FromEntry domain.Entry    `json:"from_entry"`
    ToEntry   domain.Entry    `json:"to_entry"`
    Transfer  domain.Transfer `json:"transfer"`
}

type DepositDto struct {
    WalletID int64 `json:"wallet_id"`
    Amount   int64 `json:"amount"`
}

type WithdrawDto struct {
    WalletID int64 `json:"wallet_id"`
    Amount   int64 `json:"amount"`
    UserId   int64 `json:"user_id"`
}

type WalletDto struct {
    ID            int64               `json:"id"`
    Name          string              `json:"name"`
    Address       string              `json:"address"`
    Status        domain.WalletStatus `json:"status"`
    UserID        int64               `json:"user_id"`
    BankAccountID int64               `json:"bank_account_id"`
    Balance       int64               `json:"balance"`
    Currency      string              `json:"currency"`
    CreatedAt     time.Time           `json:"created_at"`
    UpdatedAt     time.Time           `json:"updated_at"`
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
        ID: wallet.ID,
        Name: wallet.Name,
        Address: wallet.Address,
        Status: wallet.Status,
        UserID: wallet.UserID,
        BankAccountID: wallet.BankAccountID,
        Balance: wallet.Balance,
        Currency: wallet.Currency,
        CreatedAt: wallet.CreatedAt,
        UpdatedAt: wallet.UpdatedAt,
    }
}
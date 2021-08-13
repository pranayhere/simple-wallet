package dto

import (
    "github.com/pranayhere/simple-wallet/domain"
    "time"
)

type CreateBankAccountDto struct {
    AccountNo string `json:"account_no" validate:"required"`
    Ifsc      string `json:"ifsc" validate:"required"`
    BankName  string `json:"bank_name" validate:"required"`
    Currency  string `json:"currency" validate:"required"`
    UserID    int64  `json:"user_id" validate:"required"`
}

type BankAccountDto struct {
    ID        int64                    `json:"id" validate:"required"`
    AccountNo string                   `json:"account_no" validate:"required"`
    Ifsc      string                   `json:"ifsc" validate:"required"`
    BankName  string                   `json:"bank_name" validate:"required"`
    Status    domain.BankAccountStatus `json:"status" validate:"required"`
    UserID    int64                    `json:"user_id" validate:"required"`
    Currency  string                   `json:"currency" validate:"required"`
    CreatedAt time.Time                `json:"created_at" validate:"required"`
    UpdatedAt time.Time                `json:"updated_at" validate:"required"`
}

type BankAccountVerificationDto struct {
    BankAccountID int64 `json:"bank_account_id" validate:"required"`
}

func NewBankAccountDto(bankAcct domain.BankAccount) BankAccountDto {
    return BankAccountDto{
        ID:        bankAcct.ID,
        Currency:  bankAcct.Currency,
        AccountNo: bankAcct.AccountNo,
        Ifsc:      bankAcct.Ifsc,
        BankName:  bankAcct.BankName,
        UserID:    bankAcct.UserID,
        Status:    bankAcct.Status,
        CreatedAt: bankAcct.CreatedAt,
        UpdatedAt: bankAcct.UpdatedAt,
    }
}

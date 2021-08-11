package dto

import (
    "github.com/pranayhere/simple-wallet/domain"
    "time"
)

type CreateBankAccountDto struct {
    AccountNo string `json:"account_no"`
    Ifsc      string `json:"ifsc"`
    BankName  string `json:"bank_name"`
    Currency  string `json:"currency"`
    UserID    int64  `json:"user_id"`
}

type BankAccountDto struct {
    ID        int64                    `json:"id"`
    AccountNo string                   `json:"account_no"`
    Ifsc      string                   `json:"ifsc"`
    BankName  string                   `json:"bank_name"`
    Status    domain.BankAccountStatus `json:"status"`
    UserID    int64                    `json:"user_id"`
    Currency  string                   `json:"currency"`
    CreatedAt time.Time                `json:"created_at"`
    UpdatedAt time.Time                `json:"updated_at"`
}

type BankAccountVerificationDto struct {
    BankAccountID int64 `json:"bank_account_id"`
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

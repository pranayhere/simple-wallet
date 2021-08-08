package domains

import (
    "fmt"
    "time"
)

type BankAccountStatus string

const (
    BankAccountStatusINVERIFICATION     BankAccountStatus = "IN_VERIFICATION"
    BankAccountStatusVERIFIED           BankAccountStatus = "VERIFIED"
    BankAccountStatusVERIFICATIONFAILED BankAccountStatus = "VERIFICATION_FAILED"
)

type BankAccount struct {
    ID        int64             `json:"id"`
    AccountNo string            `json:"account_no"`
    Ifsc      string            `json:"ifsc"`
    BankName  string            `json:"bank_name"`
    Status    BankAccountStatus `json:"status"`
    UserID    int64             `json:"user_id"`
    Currency  string            `json:"currency"`
    CreatedAt time.Time         `json:"created_at"`
    UpdatedAt time.Time         `json:"updated_at"`
}

func (e *BankAccountStatus) Scan(src interface{}) error {
    switch s := src.(type) {
    case []byte:
        *e = BankAccountStatus(s)
    case string:
        *e = BankAccountStatus(s)
    default:
        return fmt.Errorf("unsupported scan type for BankAccountStatus: %T", src)
    }
    return nil
}

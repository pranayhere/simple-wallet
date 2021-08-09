package domains

import (
    "fmt"
    "time"
)

type WalletStatus string

const (
    WalletStatusACTIVE   WalletStatus = "ACTIVE"
    WalletStatusINACTIVE WalletStatus = "INACTIVE"
    WalletStatusBLOCKED  WalletStatus = "BLOCKED"
)

type Wallet struct {
    ID            int64        `json:"id"`
    Name          string       `json:"name"`
    Address       string       `json:"address"`
    Status        WalletStatus `json:"status"`
    UserID        int64        `json:"user_id"`
    BankAccountID int64        `json:"bank_account_id"`
    Balance       int64        `json:"balance"`
    Currency      string       `json:"currency"`
    CreatedAt     time.Time    `json:"created_at"`
    UpdatedAt     time.Time    `json:"updated_at"`
}

func (e *Wallet) IsBalanceSufficient(expectedAmount int64) bool {
    return e.Balance >= expectedAmount
}

func (e *WalletStatus) Scan(src interface{}) error {
    switch s := src.(type) {
    case []byte:
        *e = WalletStatus(s)
    case string:
        *e = WalletStatus(s)
    default:
        return fmt.Errorf("unsupported scan type for WalletStatus: %T", src)
    }
    return nil
}

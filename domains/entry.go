package domains

import (
    "fmt"
    "time"
)

type EntryType string

const (
    EntryTypeCREDIT EntryType = "CREDIT"
    EntryTypeDEBIT  EntryType = "DEBIT"
)

type Entry struct {
    ID         int64     `json:"id"`
    EntryType  EntryType `json:"entry_type"`
    WalletID   int64     `json:"wallet_id"`
    Amount     int64     `json:"amount"`
    Balance    int64     `json:"balance"`
    TransferID int64     `json:"transfer_id"`
    CreatedAt  time.Time `json:"created_at"`
}

func (e *EntryType) Scan(src interface{}) error {
    switch s := src.(type) {
    case []byte:
        *e = EntryType(s)
    case string:
        *e = EntryType(s)
    default:
        return fmt.Errorf("unsupported scan type for EntryType: %T", src)
    }
    return nil
}

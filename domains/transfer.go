package domains

import (
    "database/sql"
    "fmt"
    "time"
)

type TransferType string

const (
    TransferTypeDEPOSITTOBANK    TransferType = "DEPOSIT_TO_BANK"
    TransferTypeWITHDRAWFROMBANK TransferType = "WITHDRAW_FROM_BANK"
    TransferTypeSENDMONEY        TransferType = "SEND_MONEY"
    TransferTypePURCHASE         TransferType = "PURCHASE"
    TransferTypeREFUND           TransferType = "REFUND"
)

type Transfer struct {
    ID           int64         `json:"id"`
    TransferType         TransferType  `json:"type"`
    FromWalletID sql.NullInt64 `json:"from_wallet_id"`
    ToWalletID   sql.NullInt64 `json:"to_wallet_id"`
    Amount       int64         `json:"amount"`
    CreatedAt    time.Time     `json:"created_at"`
}

func (e *TransferType) Scan(src interface{}) error {
    switch s := src.(type) {
    case []byte:
        *e = TransferType(s)
    case string:
        *e = TransferType(s)
    default:
        return fmt.Errorf("unsupported scan type for TransferType: %T", src)
    }
    return nil
}

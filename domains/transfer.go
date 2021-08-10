package domains

import (
    "fmt"
    "time"
)

type TransferType string

const (
    TransferTypeDEPOSITTOWALLET    TransferType = "DEPOSIT_TO_WALLET"
    TransferTypeWITHDRAWFROMWALLET TransferType = "WITHDRAW_FROM_WALLET"
    TransferTypeSENDMONEY          TransferType = "SEND_MONEY"
)

type Transfer struct {
    ID           int64         `json:"id"`
    TransferType TransferType  `json:"type"`
    FromWalletID int64 `json:"from_wallet_id"`
    ToWalletID   int64 `json:"to_wallet_id"`
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

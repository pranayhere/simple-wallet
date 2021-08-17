package domain

import (
    "fmt"
    "time"
)

type PaymentRequestStatus string

const (
    PaymentRequestStatusWAITINGAPPROVAL PaymentRequestStatus = "WAITING_APPROVAL"
    PaymentRequestStatusAPPROVED        PaymentRequestStatus = "APPROVED"
    PaymentRequestStatusREFUSED         PaymentRequestStatus = "REFUSED"
    PaymentRequestStatusPAYMENTSUCCESS  PaymentRequestStatus = "PAYMENT_SUCCESS"
    PaymentRequestStatusPAYMENTFAILED   PaymentRequestStatus = "PAYMENT_FAILED"
)

func (e *PaymentRequestStatus) Scan(src interface{}) error {
    switch s := src.(type) {
    case []byte:
        *e = PaymentRequestStatus(s)
    case string:
        *e = PaymentRequestStatus(s)
    default:
        return fmt.Errorf("unsupported scan type for PaymentRequestStatus: %T", src)
    }
    return nil
}

type PaymentRequest struct {
    ID           int64                `json:"id"`
    FromWalletID int64                `json:"from_wallet_id"`
    ToWalletID   int64                `json:"to_wallet_id"`
    Amount       int64                `json:"amount"`
    Status       PaymentRequestStatus `json:"status"`
    CreatedAt    time.Time            `json:"created_at"`
}

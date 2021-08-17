package dto

import (
    "github.com/pranayhere/simple-wallet/domain"
    "time"
)

type PaymentRequestDto struct {
    ID                int64                       `json:"id"`
    FromWalletID      int64                       `json:"from_wallet_id"`
    ToWalletID        int64                       `json:"to_wallet_id"`
    FromWalletAddress string                      `json:"from_wallet_address,omitempty" validate:"required"`
    ToWalletAddress   string                      `json:"to_wallet_address,omitempty"  validate:"required"`
    Amount            int64                       `json:"amount"  validate:"required,gt=0"`
    Status            domain.PaymentRequestStatus `json:"status"`
    CreatedAt         time.Time                   `json:"created_at"`
}

type ListPaymentRequestsDto struct {
    FromWalletID int64 `json:"from_wallet_id"`
    Limit        int32 `json:"limit"`
    Offset       int32 `json:"offset"`
}

func NewPaymentRequestDto(domain domain.PaymentRequest) PaymentRequestDto {
    return PaymentRequestDto{
        ID:           domain.ID,
        FromWalletID: domain.FromWalletID,
        ToWalletID:   domain.ToWalletID,
        Amount:       domain.Amount,
        Status:       domain.Status,
        CreatedAt:    domain.CreatedAt,
    }
}

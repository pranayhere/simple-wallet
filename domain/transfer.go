package domain

import (
    "time"
)

type Transfer struct {
    ID           int64     `json:"id"`
    FromWalletID int64     `json:"from_wallet_id"`
    ToWalletID   int64     `json:"to_wallet_id"`
    Amount       int64     `json:"amount"`
    CreatedAt    time.Time `json:"created_at"`
}

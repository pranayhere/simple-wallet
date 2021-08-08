package domains

import "time"

type Entry struct {
    ID        int64     `json:"id"`
    WalletID  int64     `json:"wallet_id"`
    Amount    int64     `json:"amount"`
    CreatedAt time.Time `json:"created_at"`
}

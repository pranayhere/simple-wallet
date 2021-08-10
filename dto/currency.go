package dto

import "time"

type CurrencyDto struct {
    Code      string    `json:"code"`
    Fraction  int64     `json:"fraction"`
    CreatedAt time.Time `json:"created_at"`
}
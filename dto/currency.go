package dto

import (
    "time"
)

type CurrencyDto struct {
    Code      string    `json:"code" validate:"required"`
    Fraction  int64     `json:"fraction" validate:"required,gt=0,lte=3"`
    CreatedAt time.Time `json:"created_at"`
}

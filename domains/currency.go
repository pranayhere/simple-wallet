package domains

import "time"

type Currency struct {
    Code      string    `json:"code"`
    Fraction  int64     `json:"fraction"`
    CreatedAt time.Time `json:"created_at"`
}

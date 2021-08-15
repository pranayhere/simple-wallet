package token

import (
    "errors"
    "github.com/google/uuid"
    "time"
)

var (
    ErrExpiredToken = errors.New("token has expired")
    ErrInvalidToken = errors.New("invalid token")
)

// Payload contains the payload data of the token
type Payload struct {
    ID        uuid.UUID `json:"id"`
    UserID    int64     `json:"userID"`
    IssuedAt  time.Time `json:"issued_at"`
    ExpiredAt time.Time `json:"expire_at"`
}

// NewPayload creates a new token payload with specified username and duration
func NewPayload(userID int64, duration time.Duration) (*Payload, error) {
    tokenID, err := uuid.NewRandom()
    if err != nil {
        return nil, err
    }

    payload := &Payload{
        ID:        tokenID,
        UserID:    userID,
        IssuedAt:  time.Now(),
        ExpiredAt: time.Now().Add(duration),
    }

    return payload, nil
}

func (p *Payload) Valid() error {
    if time.Now().After(p.ExpiredAt) {
        return ErrExpiredToken
    }

    return nil
}

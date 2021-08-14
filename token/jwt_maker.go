package token

import (
    "errors"
    "fmt"
    "github.com/dgrijalva/jwt-go"
    "time"
)

const minSecretKey = 32

type JWTMaker struct {
    secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {
    if len(secretKey) < minSecretKey {
        return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKey)
    }
    return &JWTMaker{secretKey}, nil
}

// CreateToken creates a new token for specified username and duration
func (maker *JWTMaker) CreateToken(userID int64, duration time.Duration) (string, error) {
    payload, err := NewPayload(userID, duration)
    if err != nil {
        return "", err
    }

    jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
    return jwtToken.SignedString([]byte(maker.secretKey))
}

// VerifyToken checks the jwt token and validate it
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
    keyFunc := func(token *jwt.Token) (interface{}, error) {
        _, ok := token.Method.(*jwt.SigningMethodHMAC)
        if !ok {
            return nil, ErrInvalidToken
        }

        return []byte(maker.secretKey), nil
    }
    jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
    if err != nil {
        verr, ok := err.(*jwt.ValidationError)
        if ok && errors.Is(verr.Inner, ErrExpiredToken) {
            return nil, ErrExpiredToken
        }
        return nil, ErrInvalidToken
    }

    payload, ok := jwtToken.Claims.(*Payload)
    if !ok {
        return nil, ErrInvalidToken
    }

    return payload, nil
}

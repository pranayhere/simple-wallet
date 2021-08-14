package token

import (
    "github.com/dgrijalva/jwt-go"
    "github.com/pranayhere/simple-wallet/util"

    "github.com/stretchr/testify/require"
    "testing"
    "time"
)

func TestJWTMaker(t *testing.T) {
    maker, err := NewJWTMaker(util.RandomString(32))
    require.NoError(t, err)

    userID := util.RandomInt(1, 1000)
    duration := time.Minute

    issuedAt := time.Now()
    expiredAt := time.Now().Add(duration)

    token, err := maker.CreateToken(userID, duration)
    require.NoError(t, err)
    require.NotEmpty(t, token)

    payload, err := maker.VerifyToken(token)
    require.NoError(t, err)
    require.NotEmpty(t, payload)

    require.NotZero(t, payload.ID)
    require.Equal(t, userID, payload.UserID)
    require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
    require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredJWTTest(t *testing.T) {
    maker, err := NewJWTMaker(util.RandomString(32))
    require.NoError(t, err)

    token, err := maker.CreateToken(util.RandomInt(1, 1000), -time.Minute)
    require.NoError(t, err)
    require.NotEmpty(t, token)

    payload, err := maker.VerifyToken(token)
    require.Error(t, err)
    require.EqualError(t, err, ErrExpiredToken.Error())
    require.Nil(t, payload)
}

func TestInvalidJWTTokenALgoNone(t *testing.T) {
    payload, err := NewPayload(util.RandomInt(1, 1000), time.Minute)
    require.NoError(t, err)

    jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
    token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
    require.NoError(t, err)

    maker, err := NewJWTMaker(util.RandomString(32))
    require.NoError(t, err)

    payload, err = maker.VerifyToken(token)
    require.Error(t, err)
    require.EqualError(t, err, ErrInvalidToken.Error())
    require.Nil(t, payload)
}

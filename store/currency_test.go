package store_test

import (
    "context"
    "github.com/pranayhere/simple-wallet/domain"
    "github.com/pranayhere/simple-wallet/store"
    "github.com/pranayhere/simple-wallet/util"
    "github.com/stretchr/testify/require"
    "testing"
)

func createRandomCurrency(t *testing.T, code string) domain.Currency {
    currencyRepo := store.NewCurrencyRepo(testDb)
    args := store.CreateCurrencyParams{
        Code:     code,
        Fraction: util.RandomInt(1, 3),
    }

    currency, err := currencyRepo.CreateCurrency(context.Background(), args)
    require.NoError(t, err)
    require.NotEmpty(t, currency)

    require.Equal(t, args.Code, currency.Code)
    require.Equal(t, args.Fraction, currency.Fraction)

    require.NotZero(t, currency.CreatedAt)
    return currency
}

func TestCreateCurrency(t *testing.T) {
    createRandomCurrency(t, util.RandomString(3))
}

func TestGetCurrency(t *testing.T) {
    currencyRepo := store.NewCurrencyRepo(testDb)

    currency1 := createRandomCurrency(t, util.RandomString(3))
    currency2, err := currencyRepo.GetCurrency(context.Background(), currency1.Code)

    require.NoError(t, err)
    require.NotEmpty(t, currency2)

    require.Equal(t, currency1.Code, currency2.Code)
    require.Equal(t, currency1.Fraction, currency2.Fraction)
    require.Equal(t, currency1.CreatedAt, currency2.CreatedAt)
}

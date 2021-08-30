package service

import (
    "context"
    "database/sql"
    "github.com/pranayhere/simple-wallet/dto"
    "github.com/pranayhere/simple-wallet/pkg/errors"
    "github.com/pranayhere/simple-wallet/store"
    "strings"
)

type CurrencySvc interface {
    CreateCurrency(ctx context.Context, currencyDto dto.CurrencyDto) (dto.CurrencyDto, error)
    GetCurrency(ctx context.Context, currencyCode string) (dto.CurrencyDto, error)
}

type currencyService struct {
    currencyRepo store.CurrencyRepo
}

func NewCurrencyService(currencyRepo store.CurrencyRepo) CurrencySvc {
    return &currencyService{
        currencyRepo: currencyRepo,
    }
}

func (c *currencyService) CreateCurrency(ctx context.Context, currencyDto dto.CurrencyDto) (dto.CurrencyDto, error) {
    var res dto.CurrencyDto

    arg := store.CreateCurrencyParams{
        Code:     strings.ToUpper(currencyDto.Code),
        Fraction: currencyDto.Fraction,
    }

    currency, err := c.currencyRepo.CreateCurrency(ctx, arg)
    if err != nil {
        return res, err
    }

    res = dto.CurrencyDto{
        Code:      currency.Code,
        Fraction:  currency.Fraction,
        CreatedAt: currency.CreatedAt,
    }

    return res, nil
}

func (c *currencyService) GetCurrency(ctx context.Context, currencyCode string) (dto.CurrencyDto, error) {
    var res dto.CurrencyDto

    currency, err := c.currencyRepo.GetCurrency(ctx, strings.ToUpper(currencyCode))
    if err != nil {
        if err == sql.ErrNoRows {
            return res, errors.ErrCurrencyNotFound
        }
        return res, err
    }

    res = dto.CurrencyDto{
        Code:      currency.Code,
        Fraction:  currency.Fraction,
        CreatedAt: currency.CreatedAt,
    }

    return res, nil
}
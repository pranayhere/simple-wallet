package service_test

import (
    "context"
    "database/sql"
    "github.com/golang/mock/gomock"
    "github.com/pranayhere/simple-wallet/domain"
    "github.com/pranayhere/simple-wallet/dto"
    "github.com/pranayhere/simple-wallet/pkg/errors"
    "github.com/pranayhere/simple-wallet/service"
    "github.com/pranayhere/simple-wallet/store"
    mockdb "github.com/pranayhere/simple-wallet/store/mock"
    "github.com/pranayhere/simple-wallet/util"
    "github.com/stretchr/testify/require"
    "strings"
    "testing"
)

func TestCreateCurrency(t *testing.T) {
    testcases := []struct {
        name      string
        buildStub func(mockCurrencyRepo *mockdb.MockCurrencyRepo, currencyDto dto.CurrencyDto, currency domain.Currency)
        checkResp func(t *testing.T, currencyDto dto.CurrencyDto, res dto.CurrencyDto, err error)
    }{
        {
            name: "Ok",
            buildStub: func(mockCurrencyRepo *mockdb.MockCurrencyRepo, currencyDto dto.CurrencyDto, currency domain.Currency) {
                arg := store.CreateCurrencyParams{
                    Code:     strings.ToUpper(currencyDto.Code),
                    Fraction: currencyDto.Fraction,
                }

                mockCurrencyRepo.EXPECT().CreateCurrency(gomock.Any(), arg).Times(1).Return(currency, nil)
            },
            checkResp: func(t *testing.T, currencyDto dto.CurrencyDto, res dto.CurrencyDto, err error) {
                require.NoError(t, err)
                require.NotEmpty(t, res)

                require.Equal(t, currencyDto.Code, res.Code)
                require.Equal(t, currencyDto.Fraction, res.Fraction)
            },
        },
        {
            name: "DbConnectionClosed",
            buildStub: func(mockCurrencyRepo *mockdb.MockCurrencyRepo, currencyDto dto.CurrencyDto, currency domain.Currency) {
                mockCurrencyRepo.EXPECT().CreateCurrency(gomock.Any(), gomock.Any()).Times(1).Return(domain.Currency{}, sql.ErrConnDone)
            },
            checkResp: func(t *testing.T, currencyDto dto.CurrencyDto, res dto.CurrencyDto, err error) {
                require.Error(t, err)
                require.EqualError(t, err, sql.ErrConnDone.Error())
            },
        },
    }

    for _, tc := range testcases {
        t.Run(tc.name, func(t *testing.T) {
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()

            currencyDto := util.RandomCurrencyDto()
            currency := util.RandomCurrency(currencyDto)

            mockCurrencyRepo := mockdb.NewMockCurrencyRepo(ctrl)
            tc.buildStub(mockCurrencyRepo, currencyDto, currency)

            ctx := context.TODO()
            currencySvc := service.NewCurrencyService(mockCurrencyRepo)

            res, err := currencySvc.CreateCurrency(ctx, currencyDto)
            tc.checkResp(t, currencyDto, res, err)
        })
    }
}

func TestGetCurrency(t *testing.T) {
    testcases := []struct {
        name      string
        reqDto    func(t *testing.T) dto.CurrencyDto
        buildStub func(mockCurrencyRepo *mockdb.MockCurrencyRepo, currencyDto dto.CurrencyDto, currency domain.Currency)
        checkResp func(t *testing.T, currencyDto dto.CurrencyDto, res dto.CurrencyDto, err error)
    }{
        {
            name: "Ok",
            reqDto: func(t *testing.T) dto.CurrencyDto {
                return util.RandomCurrencyDto()
            },
            buildStub: func(mockCurrencyRepo *mockdb.MockCurrencyRepo, currencyDto dto.CurrencyDto, currency domain.Currency) {
                mockCurrencyRepo.EXPECT().GetCurrency(gomock.Any(), strings.ToUpper(currencyDto.Code)).Times(1).Return(currency, nil)
            },
            checkResp: func(t *testing.T, currencyDto dto.CurrencyDto, res dto.CurrencyDto, err error) {
                require.NoError(t, err)
                require.NotEmpty(t, res)

                require.Equal(t, currencyDto.Code, res.Code)
                require.Equal(t, currencyDto.Fraction, res.Fraction)
            },
        },
        {
            name: "CurrencyNotFound",
            reqDto: func(t *testing.T) dto.CurrencyDto {
                return dto.CurrencyDto{
                    Code:     "Random Currency",
                    Fraction: 2,
                }
            },
            buildStub: func(mockCurrencyRepo *mockdb.MockCurrencyRepo, currencyDto dto.CurrencyDto, currency domain.Currency) {
                mockCurrencyRepo.EXPECT().GetCurrency(gomock.Any(), strings.ToUpper(currencyDto.Code)).Times(1).Return(domain.Currency{}, sql.ErrNoRows)
            },
            checkResp: func(t *testing.T, currencyDto dto.CurrencyDto, res dto.CurrencyDto, err error) {
                require.Error(t, err)
                require.EqualError(t, err, errors.ErrCurrencyNotFound.Error())
            },
        },
        {
            name: "CurrencyNotFound",
            buildStub: func(mockCurrencyRepo *mockdb.MockCurrencyRepo, currencyDto dto.CurrencyDto, currency domain.Currency) {
                mockCurrencyRepo.EXPECT().GetCurrency(gomock.Any(), gomock.Any()).Times(1).Return(domain.Currency{}, sql.ErrConnDone)
            },
            checkResp: func(t *testing.T, currencyDto dto.CurrencyDto, res dto.CurrencyDto, err error) {
                require.Error(t, err)
                require.EqualError(t, err, sql.ErrConnDone.Error())
            },
        },
    }

    for _, tc := range testcases {
        t.Run(tc.name, func(t *testing.T) {
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()

            currencyDto := util.RandomCurrencyDto()
            currency := util.RandomCurrency(currencyDto)

            mockCurrencyRepo := mockdb.NewMockCurrencyRepo(ctrl)
            tc.buildStub(mockCurrencyRepo, currencyDto, currency)

            ctx := context.TODO()
            currencySvc := service.NewCurrencyService(mockCurrencyRepo)

            res, err := currencySvc.GetCurrency(ctx, currencyDto.Code)
            tc.checkResp(t, currencyDto, res, err)
        })
    }
}

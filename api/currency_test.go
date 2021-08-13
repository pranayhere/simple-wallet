package api_test

import (
    "bytes"
    "database/sql"
    "encoding/json"
    "fmt"
    "github.com/go-chi/chi"
    "github.com/golang/mock/gomock"
    "github.com/pranayhere/simple-wallet/api"
    "github.com/pranayhere/simple-wallet/dto"
    "github.com/pranayhere/simple-wallet/pkg/errors"
    mocksvc "github.com/pranayhere/simple-wallet/service/mock"
    "github.com/pranayhere/simple-wallet/util"
    "github.com/stretchr/testify/require"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestCreateCurrency(t *testing.T) {
    currencyDto := util.RandomCurrencyDto()

    testcases := []struct{
        name string
        body map[string]interface{}
        buildStub func(mockCurrencySvc *mocksvc.MockCurrencySvc)
        checkResp func(recorder *httptest.ResponseRecorder)
    }{
        {
            name: "Ok",
            body: map[string]interface{}{
                "code": currencyDto.Code,
                "fraction" : currencyDto.Fraction,
            },
            buildStub: func(mockCurrencySvc *mocksvc.MockCurrencySvc) {
                mockCurrencySvc.EXPECT().CreateCurrency(gomock.Any(), currencyDto).Times(1).Return(currencyDto, nil)
            },
            checkResp: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusOK, recorder.Code)
            },
        },
        {
            name: "JsonInvalid",
            body: map[string]interface{}{
                "code": currencyDto.Code,
                "fraction" : "abc",
            },
            buildStub: func(mockCurrencySvc *mocksvc.MockCurrencySvc) {
                mockCurrencySvc.EXPECT().CreateCurrency(gomock.Any(), currencyDto).Times(0)
            },
            checkResp: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusBadRequest, recorder.Code)
            },
        },
        {
            name: "EmptyCurrencyCode",
            body: map[string]interface{}{
                "fraction" : currencyDto.Fraction,
            },
            buildStub: func(mockCurrencySvc *mocksvc.MockCurrencySvc) {
                mockCurrencySvc.EXPECT().CreateCurrency(gomock.Any(), currencyDto).Times(0)
            },
            checkResp: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusBadRequest, recorder.Code)
            },
        },
        {
            name: "InvalidFraction",
            body: map[string]interface{}{
                "code": currencyDto.Code,
                "fraction" : 8,
            },
            buildStub: func(mockCurrencySvc *mocksvc.MockCurrencySvc) {
                mockCurrencySvc.EXPECT().CreateCurrency(gomock.Any(), currencyDto).Times(0)
            },
            checkResp: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusBadRequest, recorder.Code)
            },
        },
        {
            name: "DbConnectionClosed",
            body: map[string]interface{}{
                "code": currencyDto.Code,
                "fraction" : currencyDto.Fraction,
            },
            buildStub: func(mockCurrencySvc *mocksvc.MockCurrencySvc) {
                mockCurrencySvc.EXPECT().CreateCurrency(gomock.Any(), currencyDto).Times(1).Return(dto.CurrencyDto{}, sql.ErrConnDone)
            },
            checkResp: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusInternalServerError, recorder.Code)
            },
        },
    }

    for _, tc := range testcases {
        t.Run(tc.name, func(t *testing.T) {
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()

            mockCurrencySvc := mocksvc.NewMockCurrencySvc(ctrl)
            tc.buildStub(mockCurrencySvc)

            recorder := httptest.NewRecorder()
            router := chi.NewRouter()

            currencyApi := api.NewCurrencyResource(mockCurrencySvc)
            currencyApi.RegisterRoutes(router)

            data, err := json.Marshal(tc.body)
            require.NoError(t, err)

            url := "/currencies"
            request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
            require.NoError(t, err)

            router.ServeHTTP(recorder, request)

            tc.checkResp(recorder)
        })
    }
}

func TestGetCurrency(t *testing.T) {
    currencyDto := util.RandomCurrencyDto()

    testcases := []struct{
        name string
        buildStub func(mockCurrencySvc *mocksvc.MockCurrencySvc)
        checkResp func(recorder *httptest.ResponseRecorder)
    }{
        {
            name: "Ok",
            buildStub: func(mockCurrencySvc *mocksvc.MockCurrencySvc) {
                mockCurrencySvc.EXPECT().GetCurrency(gomock.Any(), currencyDto.Code).Times(1).Return(currencyDto, nil)
            },
            checkResp: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusOK, recorder.Code)
            },
        },
        {
            name: "NotFound",
            buildStub: func(mockCurrencySvc *mocksvc.MockCurrencySvc) {
                mockCurrencySvc.EXPECT().GetCurrency(gomock.Any(), currencyDto.Code).Times(1).Return(currencyDto, errors.ErrCurrencyNotFound)
            },
            checkResp: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusNotFound, recorder.Code)
            },
        },
    }

    for _, tc := range testcases {
        t.Run(tc.name, func(t *testing.T) {
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()

            mockCurrencySvc := mocksvc.NewMockCurrencySvc(ctrl)
            tc.buildStub(mockCurrencySvc)

            recorder := httptest.NewRecorder()
            router := chi.NewRouter()

            currencyApi := api.NewCurrencyResource(mockCurrencySvc)
            currencyApi.RegisterRoutes(router)

            url := fmt.Sprintf("/currencies/%s", currencyDto.Code)
            request, err := http.NewRequest(http.MethodGet, url, nil)
            require.NoError(t, err)

            router.ServeHTTP(recorder, request)
            tc.checkResp(recorder)
        })
    }
}
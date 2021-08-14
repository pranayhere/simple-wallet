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

func TestGetWallet(t *testing.T) {
    testcases := []struct {
        name      string
        url       string
        buildStub func(mockWalletSvc *mocksvc.MockWalletSvc)
        checkResp func(recorder *httptest.ResponseRecorder)
    }{
        {
            name: "Ok",
            url:  fmt.Sprintf("/wallets/%d", 1),
            buildStub: func(mockWalletSvc *mocksvc.MockWalletSvc) {
                mockWalletSvc.EXPECT().GetWalletById(gomock.Any(), gomock.Any()).Times(1)
            },
            checkResp: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusOK, recorder.Code)
            },
        },
        {
            name: "WalletNotFound",
            url:  fmt.Sprintf("/wallets/%d", 1),
            buildStub: func(mockWalletSvc *mocksvc.MockWalletSvc) {
                mockWalletSvc.EXPECT().GetWalletById(gomock.Any(), gomock.Any()).Times(1).Return(dto.WalletDto{}, errors.ErrWalletNotFound)
            },
            checkResp: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusNotFound, recorder.Code)
            },
        },
        {
            name: "InternalServerError",
            url:  fmt.Sprintf("/wallets/%d", 1),
            buildStub: func(mockWalletSvc *mocksvc.MockWalletSvc) {
                mockWalletSvc.EXPECT().GetWalletById(gomock.Any(), gomock.Any()).Times(1).Return(dto.WalletDto{}, sql.ErrConnDone)
            },
            checkResp: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusInternalServerError, recorder.Code)
            },
        },
        {
            name: "InvalidWalletID",
            url:  fmt.Sprintf("/wallets/%s", "invalid-id"),
            buildStub: func(mockWalletSvc *mocksvc.MockWalletSvc) {
                mockWalletSvc.EXPECT().GetWalletById(gomock.Any(), gomock.Any()).Times(0)
            },
            checkResp: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusBadRequest, recorder.Code)
            },
        },
    }

    for _, tc := range testcases {
        t.Run(tc.name, func(t *testing.T) {
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()

            mockWalletSvc := mocksvc.NewMockWalletSvc(ctrl)
            tc.buildStub(mockWalletSvc)

            recorder := httptest.NewRecorder()
            router := chi.NewRouter()

            walletApi := api.NewWalletResource(mockWalletSvc)
            walletApi.RegisterRoutes(router)

            url := tc.url
            request, err := http.NewRequest(http.MethodGet, url, nil)
            require.NoError(t, err)

            router.ServeHTTP(recorder, request)

            tc.checkResp(recorder)
            fmt.Println(recorder.Body)
        })
    }
}

func TestDepositToWallet(t *testing.T) {
    testcases := []struct {
        name      string
        body      map[string]interface{}
        buildStub func(mockWalletSvc *mocksvc.MockWalletSvc)
        checkResp func(recorder *httptest.ResponseRecorder)
    }{
        {
            name: "Ok",
            body: map[string]interface{}{
                "wallet_id": util.RandomInt(1, 1000),
                "amount":    util.RandomInt(1, 1000),
            },
            buildStub: func(mockWalletSvc *mocksvc.MockWalletSvc) {
                mockWalletSvc.EXPECT().Deposit(gomock.Any(), gomock.Any()).Times(1)
            },
            checkResp: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusOK, recorder.Code)
            },
        },
        {
            name: "InvalidWalletID",
            body: map[string]interface{}{
                "amount": util.RandomInt(1, 1000),
            },
            buildStub: func(mockWalletSvc *mocksvc.MockWalletSvc) {
                mockWalletSvc.EXPECT().Deposit(gomock.Any(), gomock.Any()).Times(0)
            },
            checkResp: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusBadRequest, recorder.Code)
            },
        },
        {
            name: "NegativeWalletAmount",
            body: map[string]interface{}{
                "wallet_id": util.RandomInt(1, 1000),
                "amount":    -1,
            },
            buildStub: func(mockWalletSvc *mocksvc.MockWalletSvc) {
                mockWalletSvc.EXPECT().Deposit(gomock.Any(), gomock.Any()).Times(0)
            },
            checkResp: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusBadRequest, recorder.Code)
            },
        },
        {
            name: "JsonInvalid",
            body: map[string]interface{}{
                "wallet_id": "invalid-json",
                "amount":    -1,
            },
            buildStub: func(mockWalletSvc *mocksvc.MockWalletSvc) {
                mockWalletSvc.EXPECT().Deposit(gomock.Any(), gomock.Any()).Times(0)
            },
            checkResp: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusBadRequest, recorder.Code)
            },
        },
        {
            name: "WalletNotFound",
            body: map[string]interface{}{
                "wallet_id": util.RandomInt(1, 1000),
                "amount":    util.RandomInt(1, 1000),
            },
            buildStub: func(mockWalletSvc *mocksvc.MockWalletSvc) {
                mockWalletSvc.EXPECT().Deposit(gomock.Any(), gomock.Any()).Times(1).Return(dto.WalletTransferResultDto{}, errors.ErrWalletNotFound)
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

            body := tc.body
            mockWalletSvc := mocksvc.NewMockWalletSvc(ctrl)
            tc.buildStub(mockWalletSvc)

            recorder := httptest.NewRecorder()
            router := chi.NewRouter()

            walletApi := api.NewWalletResource(mockWalletSvc)
            walletApi.RegisterRoutes(router)

            data, err := json.Marshal(body)
            require.NoError(t, err)

            url := "/wallets/deposit"
            request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
            require.NoError(t, err)

            router.ServeHTTP(recorder, request)
            tc.checkResp(recorder)
            fmt.Println(recorder.Body)
        })
    }
}

func TestWithdrawFromWallet(t *testing.T) {
    testcases := []struct {
        name      string
        body      map[string]interface{}
        buildStub func(mockWalletSvc *mocksvc.MockWalletSvc)
        checkResp func(recorder *httptest.ResponseRecorder)
    }{
        {
            name: "Ok",
            body: map[string]interface{}{
                "wallet_id": util.RandomInt(1, 1000),
                "amount":    util.RandomInt(1, 1000),
                "user_id":   util.RandomInt(1, 1000),
            },
            buildStub: func(mockWalletSvc *mocksvc.MockWalletSvc) {
                mockWalletSvc.EXPECT().Withdraw(gomock.Any(), gomock.Any()).Times(1)
            },
            checkResp: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusOK, recorder.Code)
            },
        },
        {
            name: "InvalidWalletID",
            body: map[string]interface{}{
                "amount":  util.RandomInt(1, 1000),
                "user_id": util.RandomInt(1, 1000),
            },
            buildStub: func(mockWalletSvc *mocksvc.MockWalletSvc) {
                mockWalletSvc.EXPECT().Withdraw(gomock.Any(), gomock.Any()).Times(0)
            },
            checkResp: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusBadRequest, recorder.Code)
            },
        },
        {
            name: "NegativeWalletAmount",
            body: map[string]interface{}{
                "wallet_id": util.RandomInt(1, 1000),
                "amount":    -1,
                "user_id":   util.RandomInt(1, 1000),
            },
            buildStub: func(mockWalletSvc *mocksvc.MockWalletSvc) {
                mockWalletSvc.EXPECT().Withdraw(gomock.Any(), gomock.Any()).Times(0)
            },
            checkResp: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusBadRequest, recorder.Code)
            },
        },
        {
            name: "JsonInvalid",
            body: map[string]interface{}{
                "wallet_id": "invalid-json",
                "amount":    -1,
                "user_id":   util.RandomInt(1, 1000),
            },
            buildStub: func(mockWalletSvc *mocksvc.MockWalletSvc) {
                mockWalletSvc.EXPECT().Withdraw(gomock.Any(), gomock.Any()).Times(0)
            },
            checkResp: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusBadRequest, recorder.Code)
            },
        },
        {
            name: "WalletNotFound",
            body: map[string]interface{}{
                "wallet_id": util.RandomInt(1, 1000),
                "amount":    util.RandomInt(1, 1000),
                "user_id":   util.RandomInt(1, 1000),
            },
            buildStub: func(mockWalletSvc *mocksvc.MockWalletSvc) {
                mockWalletSvc.EXPECT().Withdraw(gomock.Any(), gomock.Any()).Times(1).Return(dto.WalletTransferResultDto{}, errors.ErrWalletNotFound)
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

            body := tc.body
            mockWalletSvc := mocksvc.NewMockWalletSvc(ctrl)
            tc.buildStub(mockWalletSvc)

            recorder := httptest.NewRecorder()
            router := chi.NewRouter()

            walletApi := api.NewWalletResource(mockWalletSvc)
            walletApi.RegisterRoutes(router)

            data, err := json.Marshal(body)
            require.NoError(t, err)

            url := "/wallets/withdraw"
            request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
            require.NoError(t, err)

            router.ServeHTTP(recorder, request)
            tc.checkResp(recorder)
            fmt.Println(recorder.Body)
        })
    }
}

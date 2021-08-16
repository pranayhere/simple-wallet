package api_test

import (
    "database/sql"
    "fmt"
    "github.com/go-chi/chi"
    "github.com/golang/mock/gomock"
    "github.com/pranayhere/simple-wallet/api"
    "github.com/pranayhere/simple-wallet/dto"
    "github.com/pranayhere/simple-wallet/pkg/errors"
    mocksvc "github.com/pranayhere/simple-wallet/service/mock"
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
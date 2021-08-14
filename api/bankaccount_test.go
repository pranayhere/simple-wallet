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
    "github.com/pranayhere/simple-wallet/middleware"
    "github.com/pranayhere/simple-wallet/pkg/constant"
    "github.com/pranayhere/simple-wallet/pkg/errors"
    mocksvc "github.com/pranayhere/simple-wallet/service/mock"
    "github.com/pranayhere/simple-wallet/token"
    "github.com/pranayhere/simple-wallet/util"
    "github.com/stretchr/testify/require"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"
)

func AddAuthorization(t *testing.T, request *http.Request, tokenMaker token.Maker, authorizationType string, userID int64, duration time.Duration) {
    accessToken, err := tokenMaker.CreateToken(userID, duration)
    require.NoError(t, err)

    authorizationHeader := fmt.Sprintf("%s %s", authorizationType, accessToken)
    request.Header.Set(constant.AuthorizationHeaderKey, authorizationHeader)
}

func TestCreateBankAccount(t *testing.T) {
    createBankAccountDto := util.RandomCreateBankAccountDto("INR")
    bankAccount := util.RandomBankAccount(createBankAccountDto)
    bankAccountDto := dto.NewBankAccountDto(bankAccount)

    testcases := []struct {
        name      string
        body      map[string]interface{}
        setupAuth func(t *testing.T, request *http.Request, tokenMaker token.Maker)
        buildStub func(mockBankAcctSvc *mocksvc.MockBankAccountSvc)
        checkResp func(recorder *httptest.ResponseRecorder)
    }{
        {
            name: "Ok",
            body: map[string]interface{}{
                "account_no": createBankAccountDto.AccountNo,
                "ifsc":       createBankAccountDto.Ifsc,
                "bank_name":  createBankAccountDto.BankName,
                "currency":   createBankAccountDto.Currency,
            },
            setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
                AddAuthorization(t, request, tokenMaker, constant.AuthorizationTypeBearer, bankAccount.UserID, time.Minute)
            },
            buildStub: func(mockBankAcctSvc *mocksvc.MockBankAccountSvc) {
                mockBankAcctSvc.EXPECT().CreateBankAccount(gomock.Any(), createBankAccountDto).Times(1).Return(bankAccountDto, nil)
            },
            checkResp: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusOK, recorder.Code)
            },
        },
        {
            name: "InternalServerErr",
            body: map[string]interface{}{
                "account_no": createBankAccountDto.AccountNo,
                "ifsc":       createBankAccountDto.Ifsc,
                "bank_name":  createBankAccountDto.BankName,
                "currency":   createBankAccountDto.Currency,
            },
            setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
                AddAuthorization(t, request, tokenMaker, constant.AuthorizationTypeBearer, bankAccount.UserID, time.Minute)
            },
            buildStub: func(mockBankAcctSvc *mocksvc.MockBankAccountSvc) {
                mockBankAcctSvc.EXPECT().CreateBankAccount(gomock.Any(), createBankAccountDto).Times(1).Return(dto.BankAccountDto{}, sql.ErrConnDone)
            },
            checkResp: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusInternalServerError, recorder.Code)
            },
        },
        {
            name: "JsonInvalid",
            body: map[string]interface{}{
                "account_no": 1,
                "ifsc":       createBankAccountDto.Ifsc,
                "bank_name":  createBankAccountDto.BankName,
                "currency":   createBankAccountDto.Currency,
            },
            setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
                AddAuthorization(t, request, tokenMaker, constant.AuthorizationTypeBearer, bankAccount.UserID, time.Minute)
            },
            buildStub: func(mockBankAcctSvc *mocksvc.MockBankAccountSvc) {
                mockBankAcctSvc.EXPECT().CreateBankAccount(gomock.Any(), gomock.Any()).Times(0)
            },
            checkResp: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusBadRequest, recorder.Code)
            },
        },
        {
            name: "ValidationError",
            body: map[string]interface{}{
                "ifsc":      createBankAccountDto.Ifsc,
                "bank_name": createBankAccountDto.BankName,
                "currency":  createBankAccountDto.Currency,
            },
            setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
                AddAuthorization(t, request, tokenMaker, constant.AuthorizationTypeBearer, bankAccount.UserID, time.Minute)
            },
            buildStub: func(mockBankAcctSvc *mocksvc.MockBankAccountSvc) {
                mockBankAcctSvc.EXPECT().CreateBankAccount(gomock.Any(), gomock.Any()).Times(0)
            },
            checkResp: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusBadRequest, recorder.Code)
            },
        },
        {
            name: "DuplicateBankAccount",
            body: map[string]interface{}{
                "account_no": createBankAccountDto.AccountNo,
                "ifsc":       createBankAccountDto.Ifsc,
                "bank_name":  createBankAccountDto.BankName,
                "currency":   createBankAccountDto.Currency,
            },
            setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
                AddAuthorization(t, request, tokenMaker, constant.AuthorizationTypeBearer, bankAccount.UserID, time.Minute)
            },
            buildStub: func(mockBankAcctSvc *mocksvc.MockBankAccountSvc) {
                mockBankAcctSvc.EXPECT().CreateBankAccount(gomock.Any(), gomock.Any()).Times(1).Return(dto.BankAccountDto{}, errors.ErrBankAccountAlreadyExist)
            },
            checkResp: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusForbidden, recorder.Code)
            },
        },
    }

    for _, tc := range testcases {
        t.Run(tc.name, func(t *testing.T) {
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()

            tokenMaker, _ := token.NewJWTMaker(constant.SymmetricKey)
            mockBankAcctSvc := mocksvc.NewMockBankAccountSvc(ctrl)
            tc.buildStub(mockBankAcctSvc)

            recorder := httptest.NewRecorder()
            router := chi.NewRouter().With(middleware.Auth(tokenMaker))

            bankAcctApi := api.NewBankAccountResource(mockBankAcctSvc)
            bankAcctApi.RegisterRoutes(router)

            data, err := json.Marshal(tc.body)
            require.NoError(t, err)

            url := "/bank-accounts"
            request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
            require.NoError(t, err)
            tc.setupAuth(t, request, tokenMaker)

            router.ServeHTTP(recorder, request)

            tc.checkResp(recorder)
        })
    }
}

func TestGetBankAccount(t *testing.T) {
    createBankAccountDto := util.RandomCreateBankAccountDto("INR")
    bankAccount := util.RandomBankAccount(createBankAccountDto)
    bankAccountDto := dto.NewBankAccountDto(bankAccount)

    testcases := []struct {
        name      string
        url       string
        buildStub func(mockBankAcctSvc *mocksvc.MockBankAccountSvc)
        checkResp func(recorder *httptest.ResponseRecorder)
    }{
        {
            name: "Ok",
            url:  fmt.Sprintf("/bank-accounts/%d", 1),
            buildStub: func(mockBankAcctSvc *mocksvc.MockBankAccountSvc) {
                mockBankAcctSvc.EXPECT().GetBankAccount(gomock.Any(), int64(1)).Times(1).Return(bankAccountDto, nil)
            },
            checkResp: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusOK, recorder.Code)
            },
        },
        {
            name: "InternalServerError",
            url:  fmt.Sprintf("/bank-accounts/%d", 1),
            buildStub: func(mockBankAcctSvc *mocksvc.MockBankAccountSvc) {
                mockBankAcctSvc.EXPECT().GetBankAccount(gomock.Any(), int64(1)).Times(1).Return(bankAccountDto, sql.ErrConnDone)
            },
            checkResp: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusInternalServerError, recorder.Code)
            },
        },
        {
            name: "BankAccountNotFound",
            url:  fmt.Sprintf("/bank-accounts/%d", 1),
            buildStub: func(mockBankAcctSvc *mocksvc.MockBankAccountSvc) {
                mockBankAcctSvc.EXPECT().GetBankAccount(gomock.Any(), gomock.Any()).Times(1).Return(dto.BankAccountDto{}, errors.ErrBankAccountNotFound)
            },
            checkResp: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusNotFound, recorder.Code)
            },
        },
        {
            name: "InvalidBankAccountId",
            url:  fmt.Sprintf("/bank-accounts/%s", "invalid-id"),
            buildStub: func(mockBankAcctSvc *mocksvc.MockBankAccountSvc) {
                mockBankAcctSvc.EXPECT().GetBankAccount(gomock.Any(), gomock.Any()).Times(0)
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

            mockBankAcctSvc := mocksvc.NewMockBankAccountSvc(ctrl)
            tc.buildStub(mockBankAcctSvc)

            recorder := httptest.NewRecorder()
            router := chi.NewRouter()

            bankAcctApi := api.NewBankAccountResource(mockBankAcctSvc)
            bankAcctApi.RegisterRoutes(router)

            url := tc.url
            request, err := http.NewRequest(http.MethodGet, url, nil)
            require.NoError(t, err)

            router.ServeHTTP(recorder, request)
            tc.checkResp(recorder)
        })
    }
}

func TestBankAccountVerificationSuccess(t *testing.T) {
    createBankAccountDto := util.RandomCreateBankAccountDto("INR")
    bankAccount := util.RandomBankAccount(createBankAccountDto)
    bankAccountDto := dto.NewBankAccountDto(bankAccount)

    testcases := []struct {
        name      string
        url       string
        buildStub func(mockBankAcctSvc *mocksvc.MockBankAccountSvc)
        checkResp func(recorder *httptest.ResponseRecorder)
    }{
        {
            name: "Ok",
            url:  fmt.Sprintf("/bank-accounts/%v/verification-success", 1),
            buildStub: func(mockBankAcctSvc *mocksvc.MockBankAccountSvc) {
                mockBankAcctSvc.EXPECT().VerificationSuccess(gomock.Any(), gomock.Any()).Times(1).Return(bankAccountDto, nil)
            },
            checkResp: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusOK, recorder.Code)
            },
        },
        {
            name: "InvalidBankAccountId",
            url:  fmt.Sprintf("/bank-accounts/%s/verification-success", "invalid-id"),
            buildStub: func(mockBankAcctSvc *mocksvc.MockBankAccountSvc) {
                mockBankAcctSvc.EXPECT().VerificationSuccess(gomock.Any(), gomock.Any()).Times(0)
            },
            checkResp: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusBadRequest, recorder.Code)
            },
        },
        {
            name: "InternalServerError",
            url:  fmt.Sprintf("/bank-accounts/%d/verification-success", 1),
            buildStub: func(mockBankAcctSvc *mocksvc.MockBankAccountSvc) {
                mockBankAcctSvc.EXPECT().VerificationSuccess(gomock.Any(), gomock.Any()).Times(1).Return(dto.BankAccountDto{}, sql.ErrConnDone)
            },
            checkResp: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusInternalServerError, recorder.Code)
            },
        },
        {
            name: "BankAccountNotFound",
            url:  fmt.Sprintf("/bank-accounts/%d/verification-success", 1),
            buildStub: func(mockBankAcctSvc *mocksvc.MockBankAccountSvc) {
                mockBankAcctSvc.EXPECT().VerificationSuccess(gomock.Any(), gomock.Any()).Times(1).Return(dto.BankAccountDto{}, errors.ErrBankAccountNotFound)
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

            mockBankAcctSvc := mocksvc.NewMockBankAccountSvc(ctrl)
            tc.buildStub(mockBankAcctSvc)

            recorder := httptest.NewRecorder()
            router := chi.NewRouter()

            bankAcctApi := api.NewBankAccountResource(mockBankAcctSvc)
            bankAcctApi.RegisterRoutes(router)

            request, err := http.NewRequest(http.MethodPatch, tc.url, nil)
            require.NoError(t, err)

            router.ServeHTTP(recorder, request)
            tc.checkResp(recorder)
        })
    }
}

func TestBankAccountVerificationFailed(t *testing.T) {
    createBankAccountDto := util.RandomCreateBankAccountDto("INR")
    bankAccount := util.RandomBankAccount(createBankAccountDto)
    bankAccountDto := dto.NewBankAccountDto(bankAccount)

    testcases := []struct {
        name      string
        url       string
        buildStub func(mockBankAcctSvc *mocksvc.MockBankAccountSvc)
        checkResp func(recorder *httptest.ResponseRecorder)
    }{
        {
            name: "Ok",
            url:  fmt.Sprintf("/bank-accounts/%v/verification-failed", 1),
            buildStub: func(mockBankAcctSvc *mocksvc.MockBankAccountSvc) {
                mockBankAcctSvc.EXPECT().VerificationFailed(gomock.Any(), gomock.Any()).Times(1).Return(bankAccountDto, nil)
            },
            checkResp: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusOK, recorder.Code)
            },
        },
        {
            name: "InvalidBankAccountId",
            url:  fmt.Sprintf("/bank-accounts/%s/verification-failed", "invalid-id"),
            buildStub: func(mockBankAcctSvc *mocksvc.MockBankAccountSvc) {
                mockBankAcctSvc.EXPECT().VerificationFailed(gomock.Any(), gomock.Any()).Times(0)
            },
            checkResp: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusBadRequest, recorder.Code)
            },
        },
        {
            name: "InternalServerError",
            url:  fmt.Sprintf("/bank-accounts/%d/verification-failed", 1),
            buildStub: func(mockBankAcctSvc *mocksvc.MockBankAccountSvc) {
                mockBankAcctSvc.EXPECT().VerificationFailed(gomock.Any(), gomock.Any()).Times(1).Return(dto.BankAccountDto{}, sql.ErrConnDone)
            },
            checkResp: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusInternalServerError, recorder.Code)
            },
        },
        {
            name: "BankAccountNotFound",
            url:  fmt.Sprintf("/bank-accounts/%d/verification-failed", 1),
            buildStub: func(mockBankAcctSvc *mocksvc.MockBankAccountSvc) {
                mockBankAcctSvc.EXPECT().VerificationFailed(gomock.Any(), gomock.Any()).Times(1).Return(dto.BankAccountDto{}, errors.ErrBankAccountNotFound)
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

            mockBankAcctSvc := mocksvc.NewMockBankAccountSvc(ctrl)
            tc.buildStub(mockBankAcctSvc)

            recorder := httptest.NewRecorder()
            router := chi.NewRouter()

            bankAcctApi := api.NewBankAccountResource(mockBankAcctSvc)
            bankAcctApi.RegisterRoutes(router)

            request, err := http.NewRequest(http.MethodPatch, tc.url, nil)
            require.NoError(t, err)

            router.ServeHTTP(recorder, request)
            tc.checkResp(recorder)
        })
    }
}

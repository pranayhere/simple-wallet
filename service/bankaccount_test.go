package service_test

import (
    "context"
    "database/sql"
    "github.com/golang/mock/gomock"
    "github.com/pranayhere/simple-wallet/common"
    "github.com/pranayhere/simple-wallet/domain"
    "github.com/pranayhere/simple-wallet/dto"
    "github.com/pranayhere/simple-wallet/service"
    "github.com/pranayhere/simple-wallet/store"
    mockdb "github.com/pranayhere/simple-wallet/store/mock"
    "github.com/pranayhere/simple-wallet/util"
    "github.com/stretchr/testify/require"
    "strings"
    "testing"
)

func TestCreateBankAccount(t *testing.T) {
    testcases := []struct {
        name        string
        currency    func(t *testing.T) (dto.CurrencyDto, domain.Currency)
        bankAccount func(t *testing.T, currency string) (dto.CreateBankAccountDto, domain.BankAccount)
        buildStub   func(mockCurrencyRepo *mockdb.MockCurrencyRepo, mockBankAcctRepo *mockdb.MockBankAccountRepo, currency domain.Currency, bankAccountDto dto.CreateBankAccountDto, bankAccount domain.BankAccount)
        checkResp   func(t *testing.T, dto dto.CreateBankAccountDto, res dto.BankAccountDto, err error)
    }{
        {
            name: "Ok",
            currency: func(t *testing.T) (dto.CurrencyDto, domain.Currency) {
                currencyDto := util.RandomCurrencyDto()
                currency := util.RandomCurrency(currencyDto)
                return currencyDto, currency
            },
            bankAccount: func(t *testing.T, currency string) (dto.CreateBankAccountDto, domain.BankAccount) {
                bankAccountDto := randomCreateBankAccountDto(t, currency)
                bankAccount := randomBankAccount(t, bankAccountDto)
                return bankAccountDto, bankAccount
            },
            buildStub: func(mockCurrencyRepo *mockdb.MockCurrencyRepo, mockBankAcctRepo *mockdb.MockBankAccountRepo, currency domain.Currency, bankAccountDto dto.CreateBankAccountDto, bankAccount domain.BankAccount) {
                mockCurrencyRepo.EXPECT().GetCurrency(gomock.Any(), strings.ToUpper(currency.Code)).Times(1).Return(currency, nil)

                arg := store.CreateBankAccountWithWalletParams{
                    AccountNo: bankAccountDto.AccountNo,
                    Currency:  bankAccountDto.Currency,
                    UserID:    bankAccountDto.UserID,
                    BankName:  bankAccountDto.BankName,
                    Ifsc:      bankAccountDto.Ifsc,
                }

                bankAcctWithWalletRes := store.BankAccountWithWalletResult{
                    BankAccount: bankAccount,
                }

                mockBankAcctRepo.EXPECT().CreateBankAccountWithWallet(gomock.Any(), arg).Times(1).Return(bankAcctWithWalletRes, nil)
            },
            checkResp: func(t *testing.T, dto dto.CreateBankAccountDto, res dto.BankAccountDto, err error) {
                require.NoError(t, err)
                require.NotEmpty(t, res)

                require.Equal(t, dto.AccountNo, res.AccountNo)
                require.Equal(t, dto.BankName, res.BankName)
                require.Equal(t, dto.Ifsc, res.Ifsc)
                require.Equal(t, dto.Currency, res.Currency)
                require.Equal(t, dto.UserID, res.UserID)
            },
        },
        {
            name: "BankAccountAlreadyExist",
            currency: func(t *testing.T) (dto.CurrencyDto, domain.Currency) {
                currencyDto := util.RandomCurrencyDto()
                currency := util.RandomCurrency(currencyDto)
                return currencyDto, currency
            },
            bankAccount: func(t *testing.T, currency string) (dto.CreateBankAccountDto, domain.BankAccount) {
                bankAccountDto := randomCreateBankAccountDto(t, currency)
                bankAccount := randomBankAccount(t, bankAccountDto)
                return bankAccountDto, bankAccount
            },
            buildStub: func(mockCurrencyRepo *mockdb.MockCurrencyRepo, mockBankAcctRepo *mockdb.MockBankAccountRepo, currency domain.Currency, bankAccountDto dto.CreateBankAccountDto, bankAccount domain.BankAccount) {
                mockCurrencyRepo.EXPECT().GetCurrency(gomock.Any(), strings.ToUpper(currency.Code)).Times(1).Return(currency, nil)

                arg := store.CreateBankAccountWithWalletParams{
                    AccountNo: bankAccountDto.AccountNo,
                    Currency:  bankAccountDto.Currency,
                    UserID:    bankAccountDto.UserID,
                    BankName:  bankAccountDto.BankName,
                    Ifsc:      bankAccountDto.Ifsc,
                }

                bankAcctWithWalletRes := store.BankAccountWithWalletResult{
                    BankAccount: bankAccount,
                }

                mockBankAcctRepo.EXPECT().CreateBankAccountWithWallet(gomock.Any(), arg).Times(1).Return(bankAcctWithWalletRes, common.ErrBankAccountAlreadyExist)
            },
            checkResp: func(t *testing.T, dto dto.CreateBankAccountDto, res dto.BankAccountDto, err error) {
                require.Error(t, err)
                require.EqualError(t, err, common.ErrBankAccountAlreadyExist.Error())
            },
        },
        {
            name: "CurrencyNotFound",
            currency: func(t *testing.T) (dto.CurrencyDto, domain.Currency) {
                currencyDto := util.RandomCurrencyDto()
                currency := util.RandomCurrency(currencyDto)
                return currencyDto, currency
            },
            bankAccount: func(t *testing.T, currency string) (dto.CreateBankAccountDto, domain.BankAccount) {
                bankAccountDto := randomCreateBankAccountDto(t, currency)
                bankAccount := randomBankAccount(t, bankAccountDto)
                return bankAccountDto, bankAccount
            },
            buildStub: func(mockCurrencyRepo *mockdb.MockCurrencyRepo, mockBankAcctRepo *mockdb.MockBankAccountRepo, currency domain.Currency, bankAccountDto dto.CreateBankAccountDto, bankAccount domain.BankAccount) {
                mockCurrencyRepo.EXPECT().GetCurrency(gomock.Any(), strings.ToUpper(currency.Code)).Times(1).Return(currency, sql.ErrNoRows)
            },
            checkResp: func(t *testing.T, dto dto.CreateBankAccountDto, res dto.BankAccountDto, err error) {
                require.Error(t, err)
                require.EqualError(t, err, common.ErrCurrencyNotFound.Error())
            },
        },
        {
            name: "DbConnectionDown",
            currency: func(t *testing.T) (dto.CurrencyDto, domain.Currency) {
                currencyDto := util.RandomCurrencyDto()
                currency := util.RandomCurrency(currencyDto)
                return currencyDto, currency
            },
            bankAccount: func(t *testing.T, currency string) (dto.CreateBankAccountDto, domain.BankAccount) {
                bankAccountDto := randomCreateBankAccountDto(t, currency)
                bankAccount := randomBankAccount(t, bankAccountDto)
                return bankAccountDto, bankAccount
            },
            buildStub: func(mockCurrencyRepo *mockdb.MockCurrencyRepo, mockBankAcctRepo *mockdb.MockBankAccountRepo, currency domain.Currency, bankAccountDto dto.CreateBankAccountDto, bankAccount domain.BankAccount) {
                mockCurrencyRepo.EXPECT().GetCurrency(gomock.Any(), strings.ToUpper(currency.Code)).Times(1).Return(currency, sql.ErrConnDone)
            },
            checkResp: func(t *testing.T, dto dto.CreateBankAccountDto, res dto.BankAccountDto, err error) {
                require.Error(t, err)
                require.EqualError(t, err, sql.ErrConnDone.Error())
            },
        },
    }

    for _, tc := range testcases {
        t.Run(tc.name, func(t *testing.T) {
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()

            _, currency := tc.currency(t)
            bankAccountDto, bankAccount := tc.bankAccount(t, currency.Code)

            mockCurrencyRepo := mockdb.NewMockCurrencyRepo(ctrl)
            mockBankAcctRepo := mockdb.NewMockBankAccountRepo(ctrl)

            tc.buildStub(mockCurrencyRepo, mockBankAcctRepo, currency, bankAccountDto, bankAccount)

            ctx := context.TODO()
            currencySvc := service.NewCurrencyService(mockCurrencyRepo)
            bankAcctSvc := service.NewBankAccountService(mockBankAcctRepo, currencySvc)

            res, err := bankAcctSvc.CreateBankAccount(ctx, bankAccountDto)

            tc.checkResp(t, bankAccountDto, res, err)
        })
    }
}

func TestGetBankAccount(t *testing.T) {
    testcases := []struct {
        name      string
        buildStub func(mockBankAcctRepo *mockdb.MockBankAccountRepo)
        checkResp func(t *testing.T, err error)
    }{
        {
            name: "Ok",
            buildStub: func(mockBankAcctRepo *mockdb.MockBankAccountRepo) {
                mockBankAcctRepo.EXPECT().GetBankAccount(gomock.Any(), gomock.Any()).Times(1).Return(domain.BankAccount{}, nil)
            },
            checkResp: func(t *testing.T, err error) {
                require.NoError(t, err)
            },
        },
        {
            name: "BankAccountNotFound",
            buildStub: func(mockBankAcctRepo *mockdb.MockBankAccountRepo) {
                mockBankAcctRepo.EXPECT().GetBankAccount(gomock.Any(), gomock.Any()).Times(1).Return(domain.BankAccount{}, sql.ErrNoRows)
            },
            checkResp: func(t *testing.T, err error) {
                require.Error(t, err)
                require.EqualError(t, err, common.ErrBankAccountNotFound.Error())
            },
        },
        {
            name: "ConnectionError",
            buildStub: func(mockBankAcctRepo *mockdb.MockBankAccountRepo) {
                mockBankAcctRepo.EXPECT().GetBankAccount(gomock.Any(), gomock.Any()).Times(1).Return(domain.BankAccount{}, sql.ErrConnDone)
            },
            checkResp: func(t *testing.T, err error) {
                require.Error(t, err)
                require.EqualError(t, err, sql.ErrConnDone.Error())
            },
        },
    }
    for _, tc := range testcases {
        t.Run(tc.name, func(t *testing.T) {
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()

            mockCurrencyRepo := mockdb.NewMockCurrencyRepo(ctrl)
            mockBankAcctRepo := mockdb.NewMockBankAccountRepo(ctrl)
            tc.buildStub(mockBankAcctRepo)

            ctx := context.TODO()
            currencySvc := service.NewCurrencyService(mockCurrencyRepo)
            bankAcctSvc := service.NewBankAccountService(mockBankAcctRepo, currencySvc)

            _, err := bankAcctSvc.GetBankAccount(ctx, 1)
            tc.checkResp(t, err)
        })
    }
}

func TestBankAccountVerificationSuccess(t *testing.T) {
    testcases := []struct {
        name      string
        buildStub func(mockBankAcctRepo *mockdb.MockBankAccountRepo, bankAcct domain.BankAccount)
        checkResp func(t *testing.T, res dto.BankAccountDto, err error)
    }{
        {
            name: "Ok",
            buildStub: func(mockBankAcctRepo *mockdb.MockBankAccountRepo, bankAcct domain.BankAccount) {
                verifiedBankAcct := bankAcct
                verifiedBankAcct.Status = domain.BankAccountStatusVERIFIED

                bankAcctRes := store.BankAccountVerificationResult{
                    BankAccount: verifiedBankAcct,
                }

                arg := store.BankAccountVerificationParams{
                    BankAccountID: bankAcct.ID,
                }
                mockBankAcctRepo.EXPECT().GetBankAccount(gomock.Any(), gomock.Any()).Times(1).Return(bankAcct, nil)
                mockBankAcctRepo.EXPECT().BankAccountVerificationSuccess(gomock.Any(), arg).Times(1).Return(bankAcctRes, nil)
            },
            checkResp: func(t *testing.T, res dto.BankAccountDto, err error) {
                require.NoError(t, err)
                require.NotEmpty(t, res)
                require.Equal(t, domain.BankAccountStatusVERIFIED, res.Status)
            },
        },
        {
            name: "BankAcctVerificationError",
            buildStub: func(mockBankAcctRepo *mockdb.MockBankAccountRepo, bankAcct domain.BankAccount) {
                mockBankAcctRepo.EXPECT().GetBankAccount(gomock.Any(), gomock.Any()).Times(1).Return(bankAcct, nil)
                mockBankAcctRepo.EXPECT().BankAccountVerificationSuccess(gomock.Any(), gomock.Any()).Times(1).Return(store.BankAccountVerificationResult{}, sql.ErrConnDone)
            },
            checkResp: func(t *testing.T, res dto.BankAccountDto, err error) {
                require.Error(t, err)
                require.EqualError(t, err, sql.ErrConnDone.Error())
            },
        },
    }

    for _, tc := range testcases {
        t.Run(tc.name, func(t *testing.T) {
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()

            mockCurrencyRepo := mockdb.NewMockCurrencyRepo(ctrl)
            mockBankAcctRepo := mockdb.NewMockBankAccountRepo(ctrl)

            bankAcct := randomBankAccount(t, randomCreateBankAccountDto(t, "INR"))
            tc.buildStub(mockBankAcctRepo, bankAcct)

            ctx := context.TODO()
            currencySvc := service.NewCurrencyService(mockCurrencyRepo)
            bankAcctSvc := service.NewBankAccountService(mockBankAcctRepo, currencySvc)

            verificationDto := dto.BankAccountVerificationDto{
                BankAccountID: bankAcct.ID,
            }

            res, err := bankAcctSvc.VerificationSuccess(ctx, verificationDto)
            tc.checkResp(t, res, err)
        })
    }
}

func TestBankAccountVerificationFailed(t *testing.T) {
    testcases := []struct {
        name      string
        buildStub func(mockBankAcctRepo *mockdb.MockBankAccountRepo, bankAcct domain.BankAccount)
        checkResp func(t *testing.T, res dto.BankAccountDto, err error)
    }{
        {
            name: "Ok",
            buildStub: func(mockBankAcctRepo *mockdb.MockBankAccountRepo, bankAcct domain.BankAccount) {
                verifiedBankAcct := bankAcct
                verifiedBankAcct.Status = domain.BankAccountStatusVERIFICATIONFAILED

                bankAcctRes := store.BankAccountVerificationResult{
                    BankAccount: verifiedBankAcct,
                }

                arg := store.BankAccountVerificationParams{
                    BankAccountID: bankAcct.ID,
                }
                mockBankAcctRepo.EXPECT().GetBankAccount(gomock.Any(), gomock.Any()).Times(1).Return(bankAcct, nil)
                mockBankAcctRepo.EXPECT().BankAccountVerificationFailed(gomock.Any(), arg).Times(1).Return(bankAcctRes, nil)
            },
            checkResp: func(t *testing.T, res dto.BankAccountDto, err error) {
                require.NoError(t, err)
                require.NotEmpty(t, res)
                require.Equal(t, domain.BankAccountStatusVERIFICATIONFAILED, res.Status)
            },
        },
        {
            name: "BankAcctVerificationError",
            buildStub: func(mockBankAcctRepo *mockdb.MockBankAccountRepo, bankAcct domain.BankAccount) {
                mockBankAcctRepo.EXPECT().GetBankAccount(gomock.Any(), gomock.Any()).Times(1).Return(bankAcct, nil)
                mockBankAcctRepo.EXPECT().BankAccountVerificationFailed(gomock.Any(), gomock.Any()).Times(1).Return(store.BankAccountVerificationResult{}, sql.ErrConnDone)
            },
            checkResp: func(t *testing.T, res dto.BankAccountDto, err error) {
                require.Error(t, err)
                require.EqualError(t, err, sql.ErrConnDone.Error())
            },
        },
    }

    for _, tc := range testcases {
        t.Run(tc.name, func(t *testing.T) {
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()

            mockCurrencyRepo := mockdb.NewMockCurrencyRepo(ctrl)
            mockBankAcctRepo := mockdb.NewMockBankAccountRepo(ctrl)

            bankAcct := randomBankAccount(t, randomCreateBankAccountDto(t, "INR"))
            tc.buildStub(mockBankAcctRepo, bankAcct)

            ctx := context.TODO()
            currencySvc := service.NewCurrencyService(mockCurrencyRepo)
            bankAcctSvc := service.NewBankAccountService(mockBankAcctRepo, currencySvc)

            verificationDto := dto.BankAccountVerificationDto{
                BankAccountID: bankAcct.ID,
            }

            res, err := bankAcctSvc.VerificationFailed(ctx, verificationDto)
            tc.checkResp(t, res, err)
        })
    }
}

func randomCreateBankAccountDto(t *testing.T, currencyCode string) dto.CreateBankAccountDto {
    return dto.CreateBankAccountDto{
        AccountNo: util.RandomString(10),
        Ifsc:      util.RandomString(7),
        BankName:  util.RandomString(5),
        UserID:    util.RandomInt(1, 1000),
        Currency:  currencyCode,
    }
}

func randomBankAccount(t *testing.T, createBankAcctDto dto.CreateBankAccountDto) domain.BankAccount {
    return domain.BankAccount{
        UserID:    createBankAcctDto.UserID,
        BankName:  createBankAcctDto.BankName,
        Ifsc:      createBankAcctDto.Ifsc,
        AccountNo: createBankAcctDto.AccountNo,
        Currency:  createBankAcctDto.Currency,
    }
}

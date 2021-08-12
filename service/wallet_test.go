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
    "testing"
)

func TestSendMoney(t *testing.T)  {
    amount := int64(10)

    testcases := []struct{
        name string
        buildStub func(mockWalletRepo *mockdb.MockWalletRepo)
        checkResp func(t *testing.T, err error)
    }{
        {
            name: "Ok",
            buildStub: func(mockWalletRepo *mockdb.MockWalletRepo) {
                mockWalletRepo.EXPECT().SendMoney(gomock.Any(), gomock.Any()).Times(1)
            },
            checkResp: func(t *testing.T, err error) {
                require.NoError(t, err)
            },
        },
        {
            name: "SendMoneyTxErr",
            buildStub: func(mockWalletRepo *mockdb.MockWalletRepo) {
                mockWalletRepo.EXPECT().SendMoney(gomock.Any(), gomock.Any()).Times(1).Return(store.WalletTransferResult{}, sql.ErrTxDone)
            },
            checkResp: func(t *testing.T, err error) {
                require.Error(t, err)
                require.EqualError(t, err, sql.ErrTxDone.Error())
            },
        },
    }

    for _, tc := range testcases {
        t.Run(tc.name, func(t *testing.T) {
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()

            mockWalletRepo := mockdb.NewMockWalletRepo(ctrl)
            tc.buildStub(mockWalletRepo)

            ctx := context.TODO()
            walletSvc := service.NewWalletService(mockWalletRepo)

            sendMoneyDto := dto.SendMoneyDto{
                FromWalletAddress: util.RandomWalletAddress(util.RandomEmail()),
                ToWalletAddress: util.RandomWalletAddress(util.RandomEmail()),
                Amount: amount,
            }
            _, err := walletSvc.SendMoney(ctx, sendMoneyDto)
            tc.checkResp(t, err)
        })
    }
}

func TestDeposit(t *testing.T) {
    amount := int64(10)

    testcases := []struct{
        name string
        buildStub func(mockWalletRepo *mockdb.MockWalletRepo)
        checkResp func(t *testing.T, err error)
    }{
        {
            name: "Ok",
            buildStub: func(mockWalletRepo *mockdb.MockWalletRepo) {
                mockWalletRepo.EXPECT().DepositToWallet(gomock.Any(), gomock.Any()).Times(1)
            },
            checkResp: func(t *testing.T, err error) {
                require.NoError(t, err)
            },
        },
        {
            name: "DepositTxErr",
            buildStub: func(mockWalletRepo *mockdb.MockWalletRepo) {
                mockWalletRepo.EXPECT().DepositToWallet(gomock.Any(), gomock.Any()).Times(1).Return(store.WalletTransferResult{}, sql.ErrTxDone)
            },
            checkResp: func(t *testing.T, err error) {
                require.Error(t, err)
                require.EqualError(t, err, sql.ErrTxDone.Error())
            },
        },
    }

    for _, tc := range testcases {
        t.Run(tc.name, func(t *testing.T) {
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()

            mockWalletRepo := mockdb.NewMockWalletRepo(ctrl)
            tc.buildStub(mockWalletRepo)

            ctx := context.TODO()
            walletSvc := service.NewWalletService(mockWalletRepo)
            depositDto := dto.DepositDto{
                WalletID: util.RandomInt(1, 1000),
                Amount: amount,
            }

            _, err := walletSvc.Deposit(ctx, depositDto)
            tc.checkResp(t, err)
        })
    }
}

func TestWithdraw(t *testing.T) {
    amount := int64(10)

    testcases := []struct{
        name string
        buildStub func(mockWalletRepo *mockdb.MockWalletRepo)
        checkResp func(t *testing.T, err error)
    }{
        {
            name: "Ok",
            buildStub: func(mockWalletRepo *mockdb.MockWalletRepo) {
                mockWalletRepo.EXPECT().WithdrawFromWallet(gomock.Any(), gomock.Any()).Times(1)
            },
            checkResp: func(t *testing.T, err error) {
                require.NoError(t, err)
            },
        },
        {
            name: "DepositTxErr",
            buildStub: func(mockWalletRepo *mockdb.MockWalletRepo) {
                mockWalletRepo.EXPECT().WithdrawFromWallet(gomock.Any(), gomock.Any()).Times(1).Return(store.WalletTransferResult{}, sql.ErrTxDone)
            },
            checkResp: func(t *testing.T, err error) {
                require.Error(t, err)
                require.EqualError(t, err, sql.ErrTxDone.Error())
            },
        },
    }

    for _, tc := range testcases {
        t.Run(tc.name, func(t *testing.T) {
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()

            mockWalletRepo := mockdb.NewMockWalletRepo(ctrl)
            tc.buildStub(mockWalletRepo)

            ctx := context.TODO()
            walletSvc := service.NewWalletService(mockWalletRepo)

            withdrawDto := dto.WithdrawDto{
                WalletID: util.RandomInt(1, 1000),
                Amount: amount,
                UserId: util.RandomInt(1, 1000),
            }

            _, err := walletSvc.Withdraw(ctx, withdrawDto)
            tc.checkResp(t, err)
        })
    }
}

func TestGetWalletById(t *testing.T) {
    walletDto := randomWalletDto(util.RandomInt(1, 1000), util.RandomEmail())
    wallet := randomWallet(t, walletDto)

    testcases := []struct{
        name string
        buildStub func(mockWalletRepo *mockdb.MockWalletRepo)
        checkResp func(t *testing.T, res dto.WalletDto, err error)
    }{
        {
            name: "Ok",
            buildStub: func(mockWalletRepo *mockdb.MockWalletRepo) {
                mockWalletRepo.EXPECT().GetWallet(gomock.Any(), walletDto.ID).Times(1).Return(wallet, nil)
            },
            checkResp: func(t *testing.T, res dto.WalletDto, err error) {
                require.NoError(t, err)
                require.NotEmpty(t, res)
                require.Equal(t, walletDto.Name, res.Name)
                require.Equal(t, walletDto.Address, res.Address)
                require.Equal(t, walletDto.Status, res.Status)
                require.Equal(t, walletDto.UserID, res.UserID)
                require.Equal(t, walletDto.BankAccountID, res.BankAccountID)
                require.Equal(t, walletDto.Balance, res.Balance)
                require.Equal(t, walletDto.Currency, res.Currency)
            },
        },
        {
            name: "WalletNotFound",
            buildStub: func(mockWalletRepo *mockdb.MockWalletRepo) {
                mockWalletRepo.EXPECT().GetWallet(gomock.Any(), gomock.Any()).Times(1).Return(domain.Wallet{}, sql.ErrNoRows)
            },
            checkResp: func(t *testing.T, res dto.WalletDto, err error) {
                require.Error(t, err)
                require.EqualError(t, err, common.ErrWalletNotFound.Error())
            },
        },
        {
            name: "DbConnectionClosed",
            buildStub: func(mockWalletRepo *mockdb.MockWalletRepo) {
                mockWalletRepo.EXPECT().GetWallet(gomock.Any(), gomock.Any()).Times(1).Return(domain.Wallet{}, sql.ErrConnDone)
            },
            checkResp: func(t *testing.T, res dto.WalletDto, err error) {
                require.Error(t, err)
                require.EqualError(t, err, sql.ErrConnDone.Error())
            },
        },
    }

    for _, tc := range testcases {
        t.Run(tc.name, func(t *testing.T) {
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()

            mockWalletRepo := mockdb.NewMockWalletRepo(ctrl)
            tc.buildStub(mockWalletRepo)

            ctx := context.TODO()
            walletSvc := service.NewWalletService(mockWalletRepo)
            res, err := walletSvc.GetWalletById(ctx, walletDto.ID)
            tc.checkResp(t, res, err)
        })
    }
}

func TestGetWalletByAddress(t *testing.T) {
    walletDto := randomWalletDto(util.RandomInt(1, 1000), util.RandomEmail())
    wallet := randomWallet(t, walletDto)

    testcases := []struct{
        name string
        buildStub func(mockWalletRepo *mockdb.MockWalletRepo)
        checkResp func(t *testing.T, res dto.WalletDto, err error)
    }{
        {
            name: "Ok",
            buildStub: func(mockWalletRepo *mockdb.MockWalletRepo) {
                mockWalletRepo.EXPECT().GetWalletByAddress(gomock.Any(), walletDto.Address).Times(1).Return(wallet, nil)
            },
            checkResp: func(t *testing.T, res dto.WalletDto, err error) {
                require.NoError(t, err)
                require.NotEmpty(t, res)
                require.Equal(t, walletDto.Name, res.Name)
                require.Equal(t, walletDto.Address, res.Address)
                require.Equal(t, walletDto.Status, res.Status)
                require.Equal(t, walletDto.UserID, res.UserID)
                require.Equal(t, walletDto.BankAccountID, res.BankAccountID)
                require.Equal(t, walletDto.Balance, res.Balance)
                require.Equal(t, walletDto.Currency, res.Currency)
            },
        },
        {
            name: "WalletNotFound",
            buildStub: func(mockWalletRepo *mockdb.MockWalletRepo) {
                mockWalletRepo.EXPECT().GetWalletByAddress(gomock.Any(), gomock.Any()).Times(1).Return(domain.Wallet{}, sql.ErrNoRows)
            },
            checkResp: func(t *testing.T, res dto.WalletDto, err error) {
                require.Error(t, err)
                require.EqualError(t, err, common.ErrWalletNotFound.Error())
            },
        },
        {
            name: "DbConnectionClosed",
            buildStub: func(mockWalletRepo *mockdb.MockWalletRepo) {
                mockWalletRepo.EXPECT().GetWalletByAddress(gomock.Any(), gomock.Any()).Times(1).Return(domain.Wallet{}, sql.ErrConnDone)
            },
            checkResp: func(t *testing.T, res dto.WalletDto, err error) {
                require.Error(t, err)
                require.EqualError(t, err, sql.ErrConnDone.Error())
            },
        },
    }

    for _, tc := range testcases {
        t.Run(tc.name, func(t *testing.T) {
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()

            mockWalletRepo := mockdb.NewMockWalletRepo(ctrl)
            tc.buildStub(mockWalletRepo)

            ctx := context.TODO()
            walletSvc := service.NewWalletService(mockWalletRepo)
            res, err := walletSvc.GetWalletByAddress(ctx, walletDto.Address)
            tc.checkResp(t, res, err)
        })
    }
}

func randomWalletDto(userId int64, email string) dto.WalletDto {
    return dto.WalletDto{
        ID: util.RandomInt(1, 1000),
        Name: util.RandomString(6),
        Address: util.RandomWalletAddress(email),
        Status: domain.WalletStatusACTIVE,
        UserID: userId,
        BankAccountID: util.RandomInt(1, 1000),
        Balance: util.RandomMoney(),
        Currency: util.RandomString(3),
    }
}

func randomWallet(t *testing.T, walletDto dto.WalletDto) domain.Wallet {
    return domain.Wallet{
        ID: walletDto.ID,
        Name: walletDto.Name,
        Address: walletDto.Address,
        Status: walletDto.Status,
        UserID: walletDto.UserID,
        BankAccountID: walletDto.BankAccountID,
        Balance: walletDto.Balance,
        Currency: walletDto.Currency,
    }
}
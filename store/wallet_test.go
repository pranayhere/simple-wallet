package store_test

import (
    "context"
    "fmt"
    "github.com/pranayhere/simple-wallet/domains"
    "github.com/pranayhere/simple-wallet/store"
    "github.com/pranayhere/simple-wallet/util"
    "github.com/stretchr/testify/require"
    "strings"
    "testing"
)

func InitWalletRepo(t *testing.T) store.WalletRepo {
    transferRepo := store.NewTransferRepo(testDb)
    entryRepo := store.NewEntryRepo(testDb)
    walletRepo := store.NewWalletRepo(testDb, transferRepo, entryRepo)

    require.NotEmpty(t, transferRepo)
    require.NotEmpty(t, entryRepo)
    require.NotEmpty(t, walletRepo)

    return walletRepo
}

func createRandomWallet(t *testing.T) domains.Wallet {
    walletRepo := InitWalletRepo(t)

    user := createRandomUser(t)
    bankAccount := createRandomBankAccount(t)
    currency := createRandomCurrency(t, util.RandomString(3))
    walletAddress := strings.Split(user.Email, "@")[0]
    walletAddress = fmt.Sprintf("%s@my.wallet", walletAddress)

    args := store.CreateWalletParams{
        Name:          util.RandomString(6),
        Status:        domains.WalletStatusINACTIVE,
        UserID:        user.ID,
        BankAccountID: bankAccount.ID,
        Balance:       0,
        Currency:      currency.Code,
        Address:       walletAddress,
    }

    wallet, err := walletRepo.CreateWallet(context.Background(), args)
    require.NoError(t, err)
    require.NotEmpty(t, wallet)

    require.Equal(t, args.Name, wallet.Name)
    require.Equal(t, args.Address, wallet.Address)
    require.Equal(t, args.Currency, wallet.Currency)
    require.Equal(t, args.UserID, wallet.UserID)
    require.Equal(t, args.BankAccountID, wallet.BankAccountID)
    require.Equal(t, args.Balance, wallet.Balance)
    require.Equal(t, args.Status, wallet.Status)

    require.NotZero(t, wallet.ID)
    require.NotZero(t, wallet.CreatedAt)
    require.NotZero(t, wallet.UpdatedAt)

    return wallet
}

func TestCreateWallet(t *testing.T) {
    createRandomWallet(t)
}

func TestGetWallet(t *testing.T) {
    walletRepo := InitWalletRepo(t)
    wallet1 := createRandomWallet(t)

    wallet2, err := walletRepo.GetWallet(context.Background(), wallet1.ID)
    require.NoError(t, err)
    require.NotEmpty(t, wallet2)

    require.Equal(t, wallet1.ID, wallet2.ID)
    require.Equal(t, wallet1.Name, wallet2.Name)
    require.Equal(t, wallet1.Address, wallet2.Address)
    require.Equal(t, wallet1.Currency, wallet2.Currency)
    require.Equal(t, wallet1.UserID, wallet2.UserID)
    require.Equal(t, wallet1.BankAccountID, wallet2.BankAccountID)
    require.Equal(t, wallet1.Balance, wallet2.Balance)
    require.Equal(t, wallet1.Status, wallet2.Status)
    require.Equal(t, wallet1.CreatedAt, wallet2.CreatedAt)
    require.Equal(t, wallet1.UpdatedAt, wallet2.UpdatedAt)
}

func TestGetWalletByAddress(t *testing.T) {
    walletRepo := InitWalletRepo(t)
    wallet1 := createRandomWallet(t)

    wallet2, err := walletRepo.GetWalletByAddress(context.Background(), wallet1.Address)
    require.NoError(t, err)
    require.NotEmpty(t, wallet2)

    require.Equal(t, wallet1.ID, wallet2.ID)
    require.Equal(t, wallet1.Name, wallet2.Name)
    require.Equal(t, wallet1.Address, wallet2.Address)
    require.Equal(t, wallet1.Currency, wallet2.Currency)
    require.Equal(t, wallet1.UserID, wallet2.UserID)
    require.Equal(t, wallet1.BankAccountID, wallet2.BankAccountID)
    require.Equal(t, wallet1.Balance, wallet2.Balance)
    require.Equal(t, wallet1.Status, wallet2.Status)
    require.Equal(t, wallet1.CreatedAt, wallet2.CreatedAt)
    require.Equal(t, wallet1.UpdatedAt, wallet2.UpdatedAt)
}

func TestListWallet(t *testing.T) {
    walletRepo := InitWalletRepo(t)
    var lastWallet domains.Wallet
    for i := 0; i < 5; i++ {
        lastWallet = createRandomWallet(t)
    }

    args := store.ListWalletsParams{
        UserID: lastWallet.UserID,
        Limit:  5,
        Offset: 0,
    }

    wallets, err := walletRepo.ListWallets(context.Background(), args)
    require.NoError(t, err)
    require.NotEmpty(t, wallets)

    for _, wallet := range wallets {
        require.NotEmpty(t, wallet)
        require.Equal(t, lastWallet.Address, wallet.Address)
    }
}

func TestAddWalletBalance(t *testing.T) {
    walletRepo := InitWalletRepo(t)
    wallet1 := createRandomWallet(t)

    args := store.AddWalletBalanceParams{
        ID:     wallet1.ID,
        Amount: 100,
    }

    wallet2, err := walletRepo.AddWalletBalance(context.Background(), args)
    require.NoError(t, err)
    require.NotEmpty(t, wallet2)

    require.Equal(t, wallet1.Balance+args.Amount, wallet2.Balance)
}

func TestGetWalletByBankAccountID(t *testing.T) {
    walletRepo := InitWalletRepo(t)
    wallet1 := createRandomWallet(t)

    wallet2, err := walletRepo.GetWalletByBankAccountID(context.Background(), wallet1.BankAccountID)
    require.NoError(t, err)
    require.NotEmpty(t, wallet2)

    require.Equal(t, wallet1.ID, wallet2.ID)
    require.Equal(t, wallet1.Name, wallet2.Name)
    require.Equal(t, wallet1.Address, wallet2.Address)
    require.Equal(t, wallet1.Currency, wallet2.Currency)
    require.Equal(t, wallet1.UserID, wallet2.UserID)
    require.Equal(t, wallet1.BankAccountID, wallet2.BankAccountID)
    require.Equal(t, wallet1.Balance, wallet2.Balance)
    require.Equal(t, wallet1.Status, wallet2.Status)
    require.Equal(t, wallet1.CreatedAt, wallet2.CreatedAt)
    require.Equal(t, wallet1.UpdatedAt, wallet2.UpdatedAt)
}

func TestDepositToWallet(t *testing.T) {
    walletRepo := InitWalletRepo(t)
    wallet := createRandomWallet(t)

    verifyBankAccount(t, wallet.BankAccountID)

    n := 5
    amount := int64(100)

    errs := make(chan error)
    results := make(chan store.WalletTransferResult)

    for i := 0; i < n; i++ {
        txName := fmt.Sprintf("tx %d", i+1)
        go func() {
            ctx := context.WithValue(context.Background(), store.TxKey, txName)

            result, err := walletRepo.DepositToWallet(ctx, store.DepositeToWalletParams{
                WalletID: wallet.ID,
                Amount:   amount,
            })
            errs <- err
            results <- result
        }()
    }

    for i := 0; i < n; i++ {
        err := <-errs
        require.NoError(t, err)

        res := <-results
        require.NotEmpty(t, res)
    }

    updatedWallet, err := walletRepo.GetWallet(context.Background(), wallet.ID)
    require.NoError(t, err)
    require.NotEmpty(t, updatedWallet)

    require.Equal(t, wallet.Balance+int64(n)*amount, updatedWallet.Balance)
}

func TestWithdrawFromWallet(t *testing.T) {
    walletRepo := InitWalletRepo(t)
    wallet := createRandomWallet(t)

    verifyBankAccount(t, wallet.BankAccountID)

    initAmount := int64(2000)
    result, err := walletRepo.DepositToWallet(context.Background(), store.DepositeToWalletParams{
        WalletID: wallet.ID,
        Amount:   initAmount,
    })

    wallet = result.Wallet

    n := 5
    amount := int64(100)

    errs := make(chan error)
    results := make(chan store.WalletTransferResult)

    for i := 0; i < n; i++ {
        txName := fmt.Sprintf("tx %d", i+1)
        go func() {
            ctx := context.WithValue(context.Background(), store.TxKey, txName)

            result, err := walletRepo.WithdrawFromWallet(ctx, store.WithdrawFromWalletParams{
                WalletID: wallet.ID,
                Amount:   amount,
            })
            errs <- err
            results <- result
        }()
    }

    for i := 0; i < n; i++ {
        err := <-errs
        require.NoError(t, err)

        res := <-results
        require.NotEmpty(t, res)
    }

    updatedWallet, err := walletRepo.GetWallet(context.Background(), wallet.ID)
    require.NoError(t, err)
    require.NotEmpty(t, updatedWallet)

    require.Equal(t, wallet.Balance-int64(n)*amount, updatedWallet.Balance)
}

func TestSendMoney(t *testing.T) {
    walletRepo := InitWalletRepo(t)

    fromWallet := createRandomWallet(t)
    verifyBankAccount(t, fromWallet.BankAccountID)

    toWallet := createRandomWallet(t)
    verifyBankAccount(t, toWallet.BankAccountID)

    initAmount := int64(50)
    result, err := walletRepo.DepositToWallet(context.Background(), store.DepositeToWalletParams{
        WalletID: fromWallet.ID,
        Amount:   initAmount,
    })

    fromWallet = result.Wallet

    n := 5
    amount := int64(10)

    errs := make(chan error)
    results := make(chan store.WalletTransferResult)

    for i := 0; i < n; i++ {
        txName := fmt.Sprintf("tx %d", i+1)
        go func() {
            ctx := context.WithValue(context.Background(), store.TxKey, txName)

            result, err := walletRepo.SendMoney(ctx, store.SendMoneyParams{
                FromWalletAddress: fromWallet.Address,
                ToWalletAddress: toWallet.Address,
                Amount:        amount,
            })

            errs <- err
            results <- result
        }()
    }

    for i := 0; i < n; i++ {
        err := <-errs
        require.NoError(t, err)

        res := <-results
        require.NotEmpty(t, res)
    }

    updatedFromWallet, err := walletRepo.GetWallet(context.Background(), fromWallet.ID)
    require.NoError(t, err)
    require.NotEmpty(t, updatedFromWallet)
    require.Equal(t, fromWallet.Balance-int64(n)*amount, updatedFromWallet.Balance)

    updatedToWallet, err := walletRepo.GetWallet(context.Background(), toWallet.ID)
    require.NoError(t, err)
    require.NotEmpty(t, updatedToWallet)
    require.Equal(t, toWallet.Balance+int64(n)*amount, updatedToWallet.Balance)
}
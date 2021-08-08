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

func createRandomWallet(t *testing.T) domains.Wallet {
    walletRepo := store.NewWalletRepo(testDb)

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
    walletRepo := store.NewWalletRepo(testDb)
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
    walletRepo := store.NewWalletRepo(testDb)
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
    walletRepo := store.NewWalletRepo(testDb)
    var lastWallet domains.Wallet
    for i := 0; i < 5; i++ {
        lastWallet = createRandomWallet(t)
    }

    args := store.ListWalletsParams{
        UserID: lastWallet.UserID,
        Limit: 5,
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
    walletRepo := store.NewWalletRepo(testDb)
    wallet1 := createRandomWallet(t)

    args := store.AddWalletBalanceParams{
        ID: wallet1.ID,
        Amount: 100,
    }

    wallet2, err := walletRepo.AddWalletBalance(context.Background(), args)
    require.NoError(t, err)
    require.NotEmpty(t, wallet2)

    require.Equal(t, wallet1.Balance + args.Amount, wallet2.Balance)
}
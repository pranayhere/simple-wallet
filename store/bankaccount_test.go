package store_test

import (
    "context"
    "github.com/pranayhere/simple-wallet/domain"
    "github.com/pranayhere/simple-wallet/store"
    "github.com/pranayhere/simple-wallet/util"
    "github.com/stretchr/testify/require"
    "testing"
)

func InitBankAccountRepo(t *testing.T) store.BankAccountRepo {
    transferRepo := store.NewTransferRepo(testDb)
    entryRepo := store.NewEntryRepo(testDb)
    walletRepo := store.NewWalletRepo(testDb, transferRepo, entryRepo)
    userRepo := store.NewUserRepo(testDb)
    bankAcctRepo := store.NewBankAccountRepo(testDb, walletRepo, userRepo)

    require.NotEmpty(t, transferRepo)
    require.NotEmpty(t, entryRepo)
    require.NotEmpty(t, walletRepo)
    require.NotEmpty(t, userRepo)
    require.NotEmpty(t, bankAcctRepo)

    return bankAcctRepo
}

func createRandomBankAccount(t *testing.T) domain.BankAccount {
    bankAcctRepo := InitBankAccountRepo(t)

    user := createRandomUser(t)
    currency := createRandomCurrency(t, util.RandomString(3))

    arg := store.CreateBankAccountParams{
        AccountNo: util.RandomString(10),
        Ifsc:      util.RandomString(7),
        BankName:  util.RandomString(5),
        UserID:    user.ID,
        Status:    domain.BankAccountStatusINVERIFICATION,
        Currency:  currency.Code,
    }

    bankAcct, err := bankAcctRepo.CreateBankAccount(context.Background(), arg)
    require.NoError(t, err)
    require.NotEmpty(t, bankAcct)

    require.Equal(t, arg.AccountNo, bankAcct.AccountNo)
    require.Equal(t, arg.Ifsc, bankAcct.Ifsc)
    require.Equal(t, arg.BankName, bankAcct.BankName)
    require.Equal(t, arg.Status, bankAcct.Status)
    require.Equal(t, arg.Currency, bankAcct.Currency)

    require.NotZero(t, bankAcct.ID)
    require.NotZero(t, bankAcct.CreatedAt)
    require.NotZero(t, bankAcct.UpdatedAt)

    return bankAcct
}

func verifyBankAccount(t *testing.T, bankAccountID int64) store.BankAccountVerificationResult {
    bankAcctRepo := InitBankAccountRepo(t)
    res, err := bankAcctRepo.BankAccountVerificationSuccess(context.Background(), store.BankAccountVerificationParams{
        BankAccountID: bankAccountID,
    })
    require.NoError(t, err)
    require.NotEmpty(t, res)

    return res
}

func TestCreateBankAccount(t *testing.T) {
    createRandomBankAccount(t)
}

func TestGetBankAccount(t *testing.T) {
    bankAcctRepo := InitBankAccountRepo(t)

    bankAcct1 := createRandomBankAccount(t)
    bankAcct2, err := bankAcctRepo.GetBankAccount(context.Background(), bankAcct1.ID)
    require.NoError(t, err)
    require.NotEmpty(t, bankAcct2)

    require.Equal(t, bankAcct1.AccountNo, bankAcct2.AccountNo)
    require.Equal(t, bankAcct1.Ifsc, bankAcct2.Ifsc)
    require.Equal(t, bankAcct1.UserID, bankAcct2.UserID)
    require.Equal(t, bankAcct1.BankName, bankAcct2.BankName)
    require.Equal(t, bankAcct1.Status, bankAcct2.Status)
    require.Equal(t, bankAcct1.Currency, bankAcct2.Currency)
    require.Equal(t, bankAcct1.ID, bankAcct2.ID)
}

func TestUpdateBankAccountStatus(t *testing.T) {
    bankAcctRepo := InitBankAccountRepo(t)

    bankAcct1 := createRandomBankAccount(t)

    args := store.UpdateBankAccountStatusParams{
        ID:     bankAcct1.ID,
        Status: domain.BankAccountStatusVERIFIED,
    }

    bankAcct2, err := bankAcctRepo.UpdateBankAccountStatus(context.Background(), args)
    require.NoError(t, err)
    require.NotEmpty(t, bankAcct2)

    require.Equal(t, bankAcct1.AccountNo, bankAcct2.AccountNo)
    require.Equal(t, bankAcct1.Ifsc, bankAcct2.Ifsc)
    require.Equal(t, bankAcct1.UserID, bankAcct2.UserID)
    require.Equal(t, bankAcct1.BankName, bankAcct2.BankName)
    require.Equal(t, domain.BankAccountStatusVERIFIED, bankAcct2.Status)
    require.Equal(t, bankAcct1.Currency, bankAcct2.Currency)
    require.Equal(t, bankAcct1.ID, bankAcct2.ID)
}

func TestListBankAccounts(t *testing.T) {
    bankAcctRepo := InitBankAccountRepo(t)

    var lastBankAccount domain.BankAccount
    for i := 0; i < 5; i++ {
        lastBankAccount = createRandomBankAccount(t)
    }

    args := store.ListBankAccountsParams{
        UserID: lastBankAccount.UserID,
        Limit:  5,
        Offset: 0,
    }

    accounts, err := bankAcctRepo.ListBankAccounts(context.Background(), args)
    require.NoError(t, err)
    require.NotEmpty(t, accounts)

    for _, account := range accounts {
        require.NotEmpty(t, account)
        require.Equal(t, lastBankAccount.UserID, account.UserID)
    }
}

func TestCreateBankAccountWithWallet(t *testing.T) {
    bankAcctRepo := InitBankAccountRepo(t)

    user := createRandomUser(t)
    currency := createRandomCurrency(t, util.RandomString(3))

    res, err := bankAcctRepo.CreateBankAccountWithWallet(context.Background(), store.CreateBankAccountWithWalletParams{
        AccountNo: util.RandomString(10),
        Ifsc:      util.RandomString(7),
        BankName:  util.RandomString(5),
        UserID:    user.ID,
        Currency:  currency.Code,
    })

    require.NoError(t, err)
    require.NotEmpty(t, res)

    bankAcct := res.BankAccount
    wallet := res.Wallet

    require.NotEmpty(t, bankAcct)
    require.NotZero(t, bankAcct.ID)
    require.Equal(t, domain.BankAccountStatusINVERIFICATION, bankAcct.Status)

    require.NotEmpty(t, wallet)
    require.NotZero(t, wallet.ID)
    require.Equal(t, bankAcct.ID, wallet.BankAccountID)
    require.Equal(t, wallet.Balance, int64(0))
    require.Equal(t, domain.WalletStatusINACTIVE, wallet.Status)
}

func TestBankAccountVerificationSuccess(t *testing.T) {
    bankAcctRepo := InitBankAccountRepo(t)

    user := createRandomUser(t)
    currency := createRandomCurrency(t, util.RandomString(3))

    res, err := bankAcctRepo.CreateBankAccountWithWallet(context.Background(), store.CreateBankAccountWithWalletParams{
        AccountNo: util.RandomString(10),
        Ifsc:      util.RandomString(7),
        BankName:  util.RandomString(5),
        UserID:    user.ID,
        Currency:  currency.Code,
    })

    require.NoError(t, err)
    require.NotEmpty(t, res)
    bankAcct := res.BankAccount
    wallet := res.Wallet

    require.NotEmpty(t, bankAcct)
    require.NotEmpty(t, wallet)
    require.Equal(t, domain.BankAccountStatusINVERIFICATION, bankAcct.Status)
    require.Equal(t, domain.WalletStatusINACTIVE, wallet.Status)

    verificationRes, err := bankAcctRepo.BankAccountVerificationSuccess(context.Background(), store.BankAccountVerificationParams{
        BankAccountID: res.BankAccount.ID,
    })

    require.NoError(t, err)
    require.NotEmpty(t, res)

    verifiedBankAccount := verificationRes.BankAccount
    verifiedWallet := verificationRes.Wallet
    require.NotEmpty(t, verifiedBankAccount)
    require.NotEmpty(t, verifiedWallet)
    require.Equal(t, domain.BankAccountStatusVERIFIED, verifiedBankAccount.Status)
    require.Equal(t, domain.WalletStatusACTIVE, verifiedWallet.Status)
}

func TestAccountVerificationFailed(t *testing.T) {
    bankAcctRepo := InitBankAccountRepo(t)

    user := createRandomUser(t)
    currency := createRandomCurrency(t, util.RandomString(3))

    res, err := bankAcctRepo.CreateBankAccountWithWallet(context.Background(), store.CreateBankAccountWithWalletParams{
        AccountNo: util.RandomString(10),
        Ifsc:      util.RandomString(7),
        BankName:  util.RandomString(5),
        UserID:    user.ID,
        Currency:  currency.Code,
    })

    require.NoError(t, err)
    require.NotEmpty(t, res)
    bankAcct := res.BankAccount
    wallet := res.Wallet

    require.NotEmpty(t, bankAcct)
    require.NotEmpty(t, wallet)
    require.Equal(t, domain.BankAccountStatusINVERIFICATION, bankAcct.Status)
    require.Equal(t, domain.WalletStatusINACTIVE, wallet.Status)

    verificationRes, err := bankAcctRepo.BankAccountVerificationFailed(context.Background(), store.BankAccountVerificationParams{
        BankAccountID: res.BankAccount.ID,
    })

    require.NoError(t, err)
    require.NotEmpty(t, res)

    verifiedBankAccount := verificationRes.BankAccount

    require.NotEmpty(t, verifiedBankAccount)
    require.Equal(t, domain.BankAccountStatusVERIFICATIONFAILED, verifiedBankAccount.Status)
}

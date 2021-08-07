package db_test

import (
    "context"
    db "github.com/pranayhere/simple-wallet/db/sqlc"
    "github.com/pranayhere/simple-wallet/util"
    "github.com/stretchr/testify/require"
    "testing"
)

func createRandomBankAccount(t *testing.T) db.BankAccount {
    user := createRandomUser(t)
    currency := createRandomCurrency(t, util.RandomString(3))

    arg := db.CreateBankAccountParams{
        AccountNo: util.RandomString(10),
        Ifsc: util.RandomString(7),
        BankName: util.RandomString(5),
        UserID: user.ID,
        Status: db.BankAccountStatusINVERIFICATION,
        Currency: currency.Code,
    }

    bankAcct, err := testQueries.CreateBankAccount(context.Background(), arg)
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

func TestCreateBankAccount(t *testing.T) {
    createRandomBankAccount(t)
}

func TestGetBankAccount(t *testing.T) {
    bankAcct1 := createRandomBankAccount(t)
    bankAcct2, err := testQueries.GetBankAccount(context.Background(), bankAcct1.ID)
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
    bankAcct1 := createRandomBankAccount(t)

    args := db.UpdateBankAccountStatusParams {
        ID: bankAcct1.ID,
        Status: db.BankAccountStatusVERIFIED,
    }

    bankAcct2, err := testQueries.UpdateBankAccountStatus(context.Background(), args)
    require.NoError(t, err)
    require.NotEmpty(t, bankAcct2)

    require.Equal(t, bankAcct1.AccountNo, bankAcct2.AccountNo)
    require.Equal(t, bankAcct1.Ifsc, bankAcct2.Ifsc)
    require.Equal(t, bankAcct1.UserID, bankAcct2.UserID)
    require.Equal(t, bankAcct1.BankName, bankAcct2.BankName)
    require.Equal(t, db.BankAccountStatusVERIFIED, bankAcct2.Status)
    require.Equal(t, bankAcct1.Currency, bankAcct2.Currency)
    require.Equal(t, bankAcct1.ID, bankAcct2.ID)
}

func TestListBankAccounts(t *testing.T) {
    var lastBankAccount db.BankAccount
    for i := 0; i < 5; i++ {
        lastBankAccount = createRandomBankAccount(t)
    }

    args := db.ListBankAccountsParams{
        UserID: lastBankAccount.UserID,
        Limit:  5,
        Offset: 0,
    }

    accounts, err := testQueries.ListBankAccounts(context.Background(), args)
    require.NoError(t, err)
    require.NotEmpty(t, accounts)

    for _, account := range accounts {
        require.NotEmpty(t, account)
        require.Equal(t, lastBankAccount.UserID, account.UserID)
    }
}
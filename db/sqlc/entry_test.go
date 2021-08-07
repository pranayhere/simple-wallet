package db_test

import (
    "context"
    db "github.com/pranayhere/simple-wallet/db/sqlc"
    "github.com/pranayhere/simple-wallet/util"
    "github.com/stretchr/testify/require"
    "testing"
    "time"
)

func createRandomEntry(t *testing.T, wallet db.Wallet) db.Entry {
    args := db.CreateEntryParams{
        WalletID: wallet.ID,
        Amount:    util.RandomMoney(),
    }

    entry, err := testQueries.CreateEntry(context.Background(), args)
    require.NoError(t, err)
    require.NotEmpty(t, entry)

    require.Equal(t, args.WalletID, entry.WalletID)
    require.Equal(t, args.Amount, entry.Amount)

    require.NotZero(t, entry.ID)
    require.NotZero(t, entry.CreatedAt)

    return entry
}

func TestCreateEntry(t *testing.T) {
    account := createRandomWallet(t)
    createRandomEntry(t, account)
}

func TestGetEntry(t *testing.T) {
    account := createRandomWallet(t)
    entry1 := createRandomEntry(t, account)
    entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)

    require.NoError(t, err)
    require.NotEmpty(t, entry2)

    require.Equal(t, entry1.ID, entry2.ID)
    require.Equal(t, entry1.WalletID, entry2.WalletID)
    require.Equal(t, entry1.Amount, entry2.Amount)

    require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
}

func TestListEntries(t *testing.T) {
    account := createRandomWallet(t)
    for i := 0; i < 10; i++ {
        createRandomEntry(t, account)
    }

    arg := db.ListEntriesParams{
        WalletID: account.ID,
        Limit:     5,
        Offset:    5,
    }

    entries, err := testQueries.ListEntries(context.Background(), arg)
    require.NoError(t, err)
    require.Len(t, entries, 5)

    for _, entry := range entries {
        require.NotEmpty(t, entry)
        require.Equal(t, arg.WalletID, entry.WalletID)
    }
}
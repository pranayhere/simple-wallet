package store_test

import (
    "context"
    "github.com/pranayhere/simple-wallet/domains"
    "github.com/pranayhere/simple-wallet/store"
    "github.com/pranayhere/simple-wallet/util"
    "github.com/stretchr/testify/require"
    "testing"
    "time"
)

func createRandomEntry(t *testing.T, wallet domains.Wallet) domains.Entry {
    entryRepo := store.NewEntryRepo(testDb)
    args := store.CreateEntryParams{
        WalletID: wallet.ID,
        Amount:    util.RandomMoney(),
    }

    entry, err := entryRepo.CreateEntry(context.Background(), args)
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
    entryRepo := store.NewEntryRepo(testDb)
    account := createRandomWallet(t)
    entry1 := createRandomEntry(t, account)
    entry2, err := entryRepo.GetEntry(context.Background(), entry1.ID)

    require.NoError(t, err)
    require.NotEmpty(t, entry2)

    require.Equal(t, entry1.ID, entry2.ID)
    require.Equal(t, entry1.WalletID, entry2.WalletID)
    require.Equal(t, entry1.Amount, entry2.Amount)

    require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
}

func TestListEntries(t *testing.T) {
    entryRepo := store.NewEntryRepo(testDb)
    account := createRandomWallet(t)
    for i := 0; i < 10; i++ {
        createRandomEntry(t, account)
    }

    arg := store.ListEntriesParams{
        WalletID: account.ID,
        Limit:     5,
        Offset:    5,
    }

    entries, err := entryRepo.ListEntries(context.Background(), arg)
    require.NoError(t, err)
    require.Len(t, entries, 5)

    for _, entry := range entries {
        require.NotEmpty(t, entry)
        require.Equal(t, arg.WalletID, entry.WalletID)
    }
}
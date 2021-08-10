package store_test

import (
    "context"
    "github.com/pranayhere/simple-wallet/domain"
    "github.com/pranayhere/simple-wallet/store"
    "github.com/pranayhere/simple-wallet/util"
    "github.com/stretchr/testify/require"
    "testing"
)

func createRandomTransfer(t *testing.T, wallet1, wallet2 domain.Wallet) domain.Transfer {
    transferRepo := store.NewTransferRepo(testDb)
    arg := store.CreateTransferParams{
        FromWalletID: wallet1.ID,
        ToWalletID:   wallet2.ID,
        Amount:       util.RandomMoney(),
        TransferType: domain.TransferTypeSENDMONEY,
    }

    transfer, err := transferRepo.CreateTransfer(context.Background(), arg)
    require.NoError(t, err)
    require.NotEmpty(t, transfer)

    require.Equal(t, arg.FromWalletID, transfer.FromWalletID)
    require.Equal(t, arg.ToWalletID, transfer.ToWalletID)
    require.Equal(t, arg.Amount, transfer.Amount)

    require.NotZero(t, transfer.ID)
    require.NotZero(t, transfer.CreatedAt)

    return transfer
}

func TestCreateTransfer(t *testing.T) {
    wallet1 := createRandomWallet(t)
    wallet2 := createRandomWallet(t)
    createRandomTransfer(t, wallet1, wallet2)
}

func TestGetTransfer(t *testing.T) {
    transferRepo := store.NewTransferRepo(testDb)
    wallet1 := createRandomWallet(t)
    wallet2 := createRandomWallet(t)
    transfer1 := createRandomTransfer(t, wallet1, wallet2)

    transfer2, err := transferRepo.GetTransfer(context.Background(), transfer1.ID)
    require.NoError(t, err)
    require.NotEmpty(t, transfer2)

    require.Equal(t, transfer1.ID, transfer2.ID)
    require.Equal(t, transfer1.FromWalletID, transfer2.FromWalletID)
    require.Equal(t, transfer1.ToWalletID, transfer2.ToWalletID)
    require.Equal(t, transfer1.Amount, transfer2.Amount)
}

func TestListTransfer(t *testing.T) {
    transferRepo := store.NewTransferRepo(testDb)
    wallet1 := createRandomWallet(t)
    wallet2 := createRandomWallet(t)

    for i := 0; i < 5; i++ {
        createRandomTransfer(t, wallet1, wallet2)
        createRandomTransfer(t, wallet2, wallet1)
    }

    arg := store.ListTransfersParams{
        FromWalletID: wallet1.ID,
        ToWalletID:   wallet1.ID,
        Limit:        5,
        Offset:       5,
    }

    transfers, err := transferRepo.ListTransfers(context.Background(), arg)
    require.NoError(t, err)
    require.Len(t, transfers, 5)

    for _, transfer := range transfers {
        require.NotEmpty(t, transfer)
        require.True(t, transfer.FromWalletID == wallet1.ID || transfer.ToWalletID == wallet1.ID)
    }
}

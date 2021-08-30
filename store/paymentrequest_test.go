package store_test

import (
    "context"
    "github.com/pranayhere/simple-wallet/domain"
    "github.com/pranayhere/simple-wallet/store"
    "github.com/pranayhere/simple-wallet/util"
    "github.com/stretchr/testify/require"
    "testing"
)

func createRandomPaymentRequest(t *testing.T, wallet1, wallet2 domain.Wallet) domain.PaymentRequest {
    payReqRepo := store.NewPaymentRequestRepo(testDb)
    arg := store.CreatePaymentRequestParams{
        FromWalletID: wallet1.ID,
        ToWalletID:   wallet2.ID,
        Amount:       util.RandomMoney(),
        Status:       domain.PaymentRequestStatusWAITINGAPPROVAL,
    }

    payReq, err := payReqRepo.CreatePaymentRequest(context.Background(), arg)
    require.NoError(t, err)
    require.NotEmpty(t, payReq)

    require.Equal(t, arg.FromWalletID, payReq.FromWalletID)
    require.Equal(t, arg.ToWalletID, payReq.ToWalletID)
    require.Equal(t, arg.Amount, payReq.Amount)
    require.Equal(t, arg.Status, payReq.Status)

    require.NotZero(t, payReq.ID)
    require.NotZero(t, payReq.CreatedAt)

    return payReq
}

func TestCreatePaymentRequest(t *testing.T) {
    wallet1 := createRandomWallet(t)
    wallet2 := createRandomWallet(t)
    createRandomPaymentRequest(t, wallet1, wallet2)
}

func TestListPaymentRequests(t *testing.T) {
    payReqRepo := store.NewPaymentRequestRepo(testDb)

    var lastPayReq domain.PaymentRequest
    for i := 0; i < 5; i++ {
        wallet1 := createRandomWallet(t)
        wallet2 := createRandomWallet(t)
        lastPayReq = createRandomPaymentRequest(t, wallet1, wallet2)
    }

    args := store.ListPaymentRequestsParams{
        FromWalletID: lastPayReq.FromWalletID,
        Limit:        5,
        Offset:       0,
    }

    payReqs, err := payReqRepo.ListPaymentRequests(context.Background(), args)
    require.NoError(t, err)
    require.NotEmpty(t, payReqs)

    for _, payReq := range payReqs {
        require.NotEmpty(t, payReq)
        require.Equal(t, lastPayReq.FromWalletID, payReq.FromWalletID)
    }
}

func TestPaymentRequestStatus(t *testing.T) {
    payReqRepo := store.NewPaymentRequestRepo(testDb)
    wallet1 := createRandomWallet(t)
    wallet2 := createRandomWallet(t)

    payReq1 := createRandomPaymentRequest(t, wallet1, wallet2)
    require.Equal(t, domain.PaymentRequestStatusWAITINGAPPROVAL, payReq1.Status)

    args := store.UpdatePaymentRequestParams{
        ID:     payReq1.ID,
        Status: domain.PaymentRequestStatusAPPROVED,
    }

    payReq2, err := payReqRepo.UpdatePaymentRequest(context.Background(), args)
    require.NoError(t, err)
    require.NotEmpty(t, payReq2)

    require.Equal(t, domain.PaymentRequestStatusAPPROVED, payReq2.Status)
}

func TestGetPaymentRequest(t *testing.T) {
    payReqRepo := store.NewPaymentRequestRepo(testDb)
    wallet1 := createRandomWallet(t)
    wallet2 := createRandomWallet(t)

    payReq := createRandomPaymentRequest(t, wallet1, wallet2)

    savedPayReq, err := payReqRepo.GetPaymentRequest(context.Background(), payReq.ID)
    require.NoError(t, err)
    require.NotEmpty(t, savedPayReq)
    require.Equal(t, payReq.ID, savedPayReq.ID)
    require.Equal(t, payReq.FromWalletID, savedPayReq.FromWalletID)
    require.Equal(t, payReq.ToWalletID, savedPayReq.ToWalletID)
    require.Equal(t, payReq.Amount, savedPayReq.Amount)
    require.Equal(t, payReq.Status, savedPayReq.Status)
}
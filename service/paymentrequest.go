package service

import (
    "context"
    "database/sql"
    "github.com/pranayhere/simple-wallet/domain"
    "github.com/pranayhere/simple-wallet/dto"
    "github.com/pranayhere/simple-wallet/pkg/errors"
    "github.com/pranayhere/simple-wallet/store"
)

type PaymentRequestSvc interface {
    Create(ctx context.Context, patReqDto dto.PaymentRequestDto) (dto.PaymentRequestDto, error)
    List(ctx context.Context, patReqDto dto.ListPaymentRequestsDto) ([]dto.PaymentRequestDto, error)
    Approve(ctx context.Context, id int64) (dto.PaymentRequestDto, error)
    Refuse(ctx context.Context, id int64) (dto.PaymentRequestDto, error)
    UpdateStatus(ctx context.Context, id int64, status domain.PaymentRequestStatus) (dto.PaymentRequestDto, error)
    Get(ctx context.Context, id int64) (dto.PaymentRequestDto, error)
}

type paymentRequestService struct {
    paymentRequestRepo store.PaymentRequestRepo
    walletSvc          WalletSvc
}

func NewPaymentRequestService(paymentRequestRepo store.PaymentRequestRepo, walletSvc WalletSvc) PaymentRequestSvc {
    return &paymentRequestService{
        paymentRequestRepo: paymentRequestRepo,
        walletSvc:          walletSvc,
    }
}

func (p *paymentRequestService) Create(ctx context.Context, payReqDto dto.PaymentRequestDto) (dto.PaymentRequestDto, error) {
    var res dto.PaymentRequestDto

    fromWallet, err := p.walletSvc.GetWalletByAddress(ctx, payReqDto.FromWalletAddress)
    if err != nil {
        return res, err
    }

    toWallet, err := p.walletSvc.GetWalletByAddress(ctx, payReqDto.ToWalletAddress)
    if err != nil {
        return res, err
    }

    arg := store.CreatePaymentRequestParams{
        FromWalletID: fromWallet.ID,
        ToWalletID:   toWallet.ID,
        Amount:       payReqDto.Amount,
        Status:       domain.PaymentRequestStatusWAITINGAPPROVAL,
    }

    payReq, err := p.paymentRequestRepo.CreatePaymentRequest(ctx, arg)
    if err != nil {
        return res, err
    }

    res = dto.NewPaymentRequestDto(payReq)
    return res, nil
}

func (p *paymentRequestService) Approve(ctx context.Context, id int64) (dto.PaymentRequestDto, error) {
    var res dto.PaymentRequestDto

    _, err := p.Get(ctx, id)
    if err != nil {
        return res, err
    }

    arg := store.UpdatePaymentRequestParams{
        ID:     id,
        Status: domain.PaymentRequestStatusAPPROVED,
    }

    payReq, err := p.paymentRequestRepo.UpdatePaymentRequest(ctx, arg)
    if err != nil {
        return res, err
    }

    transferArg := dto.TransferMoneyByWalletIDDto{
        FromWalletID: payReq.FromWalletID,
        ToWalletID:   payReq.ToWalletID,
        Amount:       payReq.Amount,
    }

    _, err = p.walletSvc.PayByWalletID(ctx, transferArg)
    if err != nil {
        res, err = p.UpdateStatus(ctx, id, domain.PaymentRequestStatusPAYMENTFAILED)
        if err != nil {
            return res, err
        }

        return res, err
    }

    res, err = p.UpdateStatus(ctx, id, domain.PaymentRequestStatusPAYMENTSUCCESS)
    if err != nil {
        return res, err
    }

    return res, err

}

func (p *paymentRequestService) Refuse(ctx context.Context, id int64) (dto.PaymentRequestDto, error) {
    var res dto.PaymentRequestDto

    arg := store.UpdatePaymentRequestParams{
        ID:     id,
        Status: domain.PaymentRequestStatusAPPROVED,
    }

    payReq, err := p.paymentRequestRepo.UpdatePaymentRequest(ctx, arg)
    if err != nil {
        return res, err
    }

    res = dto.NewPaymentRequestDto(payReq)
    return res, nil
}

func (p *paymentRequestService) UpdateStatus(ctx context.Context, id int64, status domain.PaymentRequestStatus) (dto.PaymentRequestDto, error) {
    var res dto.PaymentRequestDto

    arg := store.UpdatePaymentRequestParams{
        ID:     id,
        Status: status,
    }

    payReq, err := p.paymentRequestRepo.UpdatePaymentRequest(ctx, arg)
    if err != nil {
        return res, err
    }

    res = dto.NewPaymentRequestDto(payReq)
    return res, nil
}

func (p *paymentRequestService) Get(ctx context.Context, id int64) (dto.PaymentRequestDto, error) {
    var res dto.PaymentRequestDto

    payReq, err := p.paymentRequestRepo.GetPaymentRequest(ctx, id)
    if err != nil {
        if err == sql.ErrNoRows {
            return res, errors.ErrPaymentRequestNotFound
        }

        return res, err
    }

    res = dto.NewPaymentRequestDto(payReq)
    return res, nil
}

func (p *paymentRequestService) List(ctx context.Context, listPayReqDto dto.ListPaymentRequestsDto) ([]dto.PaymentRequestDto, error) {
    res := []dto.PaymentRequestDto{}

    arg := store.ListPaymentRequestsParams{
        FromWalletID: listPayReqDto.FromWalletID,
        Limit:        listPayReqDto.Limit,
        Offset:       listPayReqDto.Offset,
    }

    payReqs, err := p.paymentRequestRepo.ListPaymentRequests(ctx, arg)
    if err != nil {
        return res, err
    }

    for _, pr := range payReqs {
        res = append(res, dto.NewPaymentRequestDto(pr))
    }

    return res, nil
}

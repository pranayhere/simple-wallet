package service

import (
    "context"
    "database/sql"
    "github.com/pranayhere/simple-wallet/dto"
    "github.com/pranayhere/simple-wallet/pkg/errors"
    "github.com/pranayhere/simple-wallet/store"
)

type WalletSvc interface {
    Pay(ctx context.Context, transferMoneyDto dto.TransferMoneyDto) (dto.WalletTransferResultDto, error)
    PayByWalletID(ctx context.Context, transferMoneyDto dto.TransferMoneyByWalletIDDto) (dto.WalletTransferResultDto, error)
    GetWalletById(ctx context.Context, id int64) (dto.WalletDto, error)
    GetWalletByAddress(ctx context.Context, address string) (dto.WalletDto, error)
}

type walletService struct {
    walletRepo store.WalletRepo
}

func NewWalletService(walletRepo store.WalletRepo) WalletSvc {
    return &walletService{
        walletRepo: walletRepo,
    }
}

func (w *walletService) Pay(ctx context.Context, transferMoneyDto dto.TransferMoneyDto) (dto.WalletTransferResultDto, error) {
    var txnResDto dto.WalletTransferResultDto

    arg := store.SendMoneyParams{
        FromWalletAddress: transferMoneyDto.FromWalletAddress,
        ToWalletAddress:   transferMoneyDto.ToWalletAddress,
        Amount:            transferMoneyDto.Amount,
    }

    res, err := w.walletRepo.SendMoney(ctx, arg)
    if err != nil {
        return txnResDto, err
    }

    txnResDto = dto.NewWalletTransferDto(res)
    return txnResDto, nil
}

func (w *walletService) PayByWalletID(ctx context.Context, transferMoneyDto dto.TransferMoneyByWalletIDDto) (dto.WalletTransferResultDto, error) {
    var res dto.WalletTransferResultDto

    fromWallet, err := w.GetWalletById(ctx, transferMoneyDto.FromWalletID)
    if err != nil {
        return res, err
    }

    toWallet, err := w.GetWalletById(ctx, transferMoneyDto.ToWalletID)
    if err != nil {
        return res, err
    }

    arg := dto.TransferMoneyDto{
        FromWalletAddress: fromWallet.Address,
        ToWalletAddress:   toWallet.Address,
        Amount:            transferMoneyDto.Amount,
    }

    res, err = w.Pay(ctx, arg)
    if err != nil {
        return dto.WalletTransferResultDto{}, err
    }

    return res, nil
}

func (w *walletService) GetWalletById(ctx context.Context, id int64) (dto.WalletDto, error) {
    var walletDto dto.WalletDto

    wallet, err := w.walletRepo.GetWallet(ctx, id)
    if err != nil {
        if err == sql.ErrNoRows {
            return walletDto, errors.ErrWalletNotFound
        }

        return walletDto, err
    }

    walletDto = dto.NewWalletDto(wallet)
    return walletDto, nil
}

func (w *walletService) GetWalletByAddress(ctx context.Context, address string) (dto.WalletDto, error) {
    var walletDto dto.WalletDto

    wallet, err := w.walletRepo.GetWalletByAddress(ctx, address)
    if err != nil {
        if err == sql.ErrNoRows {
            return walletDto, errors.ErrWalletNotFound
        }

        return walletDto, err
    }

    walletDto = dto.NewWalletDto(wallet)
    return walletDto, nil
}

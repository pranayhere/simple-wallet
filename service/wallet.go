package service

import (
    "context"
    "database/sql"
    "github.com/pranayhere/simple-wallet/dto"
    "github.com/pranayhere/simple-wallet/pkg/errors"
    "github.com/pranayhere/simple-wallet/store"
)

type WalletSvc interface {
    SendMoney(ctx context.Context, sendMoneyDto dto.SendMoneyDto) (dto.WalletTransferResultDto, error)
    Deposit(ctx context.Context, depositDto dto.DepositDto) (dto.WalletTransferResultDto, error)
    Withdraw(ctx context.Context, withdrawDto dto.WithdrawDto) (dto.WalletTransferResultDto, error)
    GetWalletById(ctx context.Context, id int64) (dto.WalletDto, error)
}

type walletService struct {
    walletRepo store.WalletRepo
}

func NewWalletService(walletRepo store.WalletRepo) WalletSvc {
    return &walletService{
        walletRepo: walletRepo,
    }
}

func (w *walletService) SendMoney(ctx context.Context, sendMoneyDto dto.SendMoneyDto) (dto.WalletTransferResultDto, error) {
    var txnResDto dto.WalletTransferResultDto

    arg := store.SendMoneyParams{
        FromWalletAddress: sendMoneyDto.FromWalletAddress,
        ToWalletAddress:   sendMoneyDto.ToWalletAddress,
        Amount:            sendMoneyDto.Amount,
    }

    res, err := w.walletRepo.SendMoney(ctx, arg)
    if err != nil {
        return txnResDto, err
    }

    txnResDto = dto.NewWalletTransferDto(res)
    return txnResDto, nil
}

func (w *walletService) Deposit(ctx context.Context, depositDto dto.DepositDto) (dto.WalletTransferResultDto, error) {
    var txnResDto dto.WalletTransferResultDto

    arg := store.DepositeToWalletParams{
        WalletID: depositDto.WalletID,
        Amount:   depositDto.Amount,
    }

    res, err := w.walletRepo.DepositToWallet(ctx, arg)
    if err != nil {
        return txnResDto, err
    }

    txnResDto = dto.NewWalletTransferDto(res)
    return txnResDto, nil
}

func (w *walletService) Withdraw(ctx context.Context, withdrawDto dto.WithdrawDto) (dto.WalletTransferResultDto, error) {
    var txnResDto dto.WalletTransferResultDto

    arg := store.WithdrawFromWalletParams{
        WalletID: withdrawDto.WalletID,
        Amount:   withdrawDto.Amount,
    }

    res, err := w.walletRepo.WithdrawFromWallet(ctx, arg)
    if err != nil {
        return txnResDto, err
    }

    txnResDto = dto.NewWalletTransferDto(res)
    return txnResDto, nil
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

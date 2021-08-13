package service

import (
    "context"
    "database/sql"
    "github.com/pranayhere/simple-wallet/common"
    "github.com/pranayhere/simple-wallet/dto"
    "github.com/pranayhere/simple-wallet/store"
)

type BankAccountSvc interface {
    CreateBankAccount(ctx context.Context, newBankAcctDto dto.CreateBankAccountDto) (dto.BankAccountDto, error)
    VerificationSuccess(ctx context.Context, verificationDto dto.BankAccountVerificationDto) (dto.BankAccountDto, error)
    VerificationFailed(ctx context.Context, verificationDto dto.BankAccountVerificationDto) (dto.BankAccountDto, error)
    GetBankAccount(ctx context.Context, bankAccountId int64) (dto.BankAccountDto, error)
}

type bankAccountService struct {
    bankAcctRepo store.BankAccountRepo
    currencySvc  CurrencySvc
}

func NewBankAccountService(bankAcctRepo store.BankAccountRepo, currencySvc CurrencySvc) BankAccountSvc {
    return &bankAccountService{
        bankAcctRepo: bankAcctRepo,
        currencySvc:  currencySvc,
    }
}

func (b *bankAccountService) CreateBankAccount(ctx context.Context, newBankAcctDto dto.CreateBankAccountDto) (dto.BankAccountDto, error) {
    var bankAcctDto dto.BankAccountDto

    currency, err := b.currencySvc.GetCurrency(ctx, newBankAcctDto.Currency)
    if err != nil {
        return bankAcctDto, err
    }

    arg := store.CreateBankAccountWithWalletParams{
        UserID:    newBankAcctDto.UserID,
        BankName:  newBankAcctDto.BankName,
        Ifsc:      newBankAcctDto.Ifsc,
        AccountNo: newBankAcctDto.AccountNo,
        Currency:  currency.Code,
    }

    bankAcct, err := b.bankAcctRepo.CreateBankAccountWithWallet(ctx, arg)
    if err != nil {
        return bankAcctDto, err
    }

    bankAcctDto = dto.NewBankAccountDto(bankAcct.BankAccount)
    return bankAcctDto, nil
}

func (b *bankAccountService) GetBankAccount(ctx context.Context, bankAccountId int64) (dto.BankAccountDto, error) {
    var bankAcctDto dto.BankAccountDto

    bankAcct, err := b.bankAcctRepo.GetBankAccount(ctx, bankAccountId)
    if err != nil {
        if err == sql.ErrNoRows {
            return bankAcctDto, common.ErrBankAccountNotFound
        }
        return bankAcctDto, err
    }

    bankAcctDto = dto.NewBankAccountDto(bankAcct)
    return bankAcctDto, nil
}

func (b *bankAccountService) VerificationSuccess(ctx context.Context, verificationDto dto.BankAccountVerificationDto) (dto.BankAccountDto, error) {
    var bankAcctDto dto.BankAccountDto

    bankAcctDto, err := b.GetBankAccount(ctx, verificationDto.BankAccountID)
    if err != nil {
        return bankAcctDto, err
    }

    arg := store.BankAccountVerificationParams{
        BankAccountID: verificationDto.BankAccountID,
    }

    res, err := b.bankAcctRepo.BankAccountVerificationSuccess(ctx, arg)
    if err != nil {
        return bankAcctDto, err
    }

    bankAcctDto = dto.NewBankAccountDto(res.BankAccount)
    return bankAcctDto, nil
}

func (b *bankAccountService) VerificationFailed(ctx context.Context, verificationDto dto.BankAccountVerificationDto) (dto.BankAccountDto, error) {
    var bankAcctDto dto.BankAccountDto

    bankAcctDto, err := b.GetBankAccount(ctx, verificationDto.BankAccountID)
    if err != nil {
        return bankAcctDto, err
    }

    arg := store.BankAccountVerificationParams{
        BankAccountID: verificationDto.BankAccountID,
    }

    res, err := b.bankAcctRepo.BankAccountVerificationFailed(ctx, arg)
    if err != nil {
        return bankAcctDto, err
    }

    bankAcctDto = dto.NewBankAccountDto(res.BankAccount)
    return bankAcctDto, nil
}

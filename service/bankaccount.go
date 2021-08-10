package service

import "context"

type BankAcctSvc interface {
    CreateBankAccount(ctx context.Context, )
    VerifyBankAccount(ctx context.Context)
    GetBankAccount(ctx context.Context)
    ListBankAccounts(ctx context.Context)
}

type BankAccountService struct {

}
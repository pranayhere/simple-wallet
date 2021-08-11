package common

import "errors"

var (
    ErrUserNotFound            = errors.New("user not found")
    ErrIncorrectPassword       = errors.New("incorrect password")
    ErrUserAlreadyExist        = errors.New("user already exist")
    ErrCurrencyNotFound        = errors.New("currency not found")
    ErrBankAccountAlreadyExist = errors.New("bank account already exist")
    ErrBankAccountNotFound     = errors.New("bank account not found")
    ErrSomethingWrong          = errors.New("something went wrong")
    ErrCurrencyMismatch        = errors.New("currency mismatch")
    ErrWalletNotFound          = errors.New("wallet not found")
)

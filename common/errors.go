package common

import "errors"

var (
    ErrUserNotFound     = errors.New("user not found")
    ErrIncorrectPassword = errors.New("incorrect password")
    ErrUserAlreadyExist = errors.New("user already exist")
)

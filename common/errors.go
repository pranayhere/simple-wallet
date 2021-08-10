package common

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
    ErrUnauthorisedUser = errors.New("unauthorised user")
    ErrUserAlreadyExist = errors.New("user already exist")
)

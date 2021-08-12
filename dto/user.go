package dto

import (
    "github.com/pranayhere/simple-wallet/domain"
    "time"
)

type CreateUserDto struct {
    Username string `json:"username" validate:"required,alphanum"`
    Password string `json:"password" validate:"required,min=6"`
    FullName string `json:"full_name" validate:"required"`
    Email    string `json:"email" validate:"required,email"`
}

type UserDto struct {
    ID                int64     `json:"id"`
    Username          string    `json:"username" validate:"required,alphanum"`
    Status            string    `json:"status"`
    FullName          string    `json:"full_name" validate:"required"`
    Email             string    `json:"email" validate:"required,email"`
    PasswordChangedAt time.Time `json:"password_changed_at"`
    CreatedAt         time.Time `json:"created_at"`
}

type LoginCredentialsDto struct {
    Username string `json:"username" validate:"required,alphanum"`
    Password string `json:"password" validate:"required,min=6"`
}

type LoggedInUserDto struct {
    AccessToken string  `json:"access_token"`
    User        UserDto `json:"user"`
}

func NewUserDto(user domain.User) UserDto {
    return UserDto{
        ID:                user.ID,
        Username:          user.Username,
        FullName:          user.FullName,
        Status:            string(user.Status),
        Email:             user.Email,
        PasswordChangedAt: user.PasswordChangedAt,
        CreatedAt:         user.CreatedAt,
    }
}

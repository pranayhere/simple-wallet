package service

import (
    "context"
    "database/sql"
    "github.com/lib/pq"
    "github.com/pkg/errors"
    "github.com/pranayhere/simple-wallet/common"
    "github.com/pranayhere/simple-wallet/domain"
    "github.com/pranayhere/simple-wallet/dto"
    "github.com/pranayhere/simple-wallet/store"
    "github.com/pranayhere/simple-wallet/token"
    "github.com/pranayhere/simple-wallet/util"
)

type UserSvc interface {
    CreateUser(ctx context.Context, createUserDto dto.CreateUserDto) (dto.UserDto, error)
    LoginUser(ctx context.Context, loginCredsDto dto.LoginCredentialsDto) (dto.LoginUserDto, error)
}

type userService struct {
    userRepo store.UserRepo
    tokenMaker token.Maker
}

func NewUserService(userRepo store.UserRepo, tokenMaker token.Maker) UserSvc {
    return &userService{
        userRepo: userRepo,
        tokenMaker: tokenMaker,
    }
}

func (u *userService) CreateUser(ctx context.Context, createUserDto dto.CreateUserDto) (dto.UserDto, error) {
    var userDto dto.UserDto

    hashedPassword, err := util.HashPassword(createUserDto.Password)
    if err != nil {
        return userDto, err
    }

    arg := store.CreateUserParams{
        Username:       createUserDto.Username,
        HashedPassword: hashedPassword,
        FullName:       createUserDto.FullName,
        Email:          createUserDto.Email,
        Status:         domain.UserStatusACTIVE,
    }

    user, err := u.userRepo.CreateUser(ctx, arg)
    if err != nil {
        if pqErr, ok := err.(*pq.Error); ok {
            switch pqErr.Code.Name() {
            case "unique_violation":
                return userDto, errors.Wrap(common.ErrUserAlreadyExist, err.Error())
            }
        }
        return userDto, err
    }

    userDto = dto.NewUserDto(user)
    return userDto, nil
}

func (u userService) LoginUser(ctx context.Context, loginCredentialsDto dto.LoginCredentialsDto) (dto.LoginUserDto, error) {
    var loginDto dto.LoginUserDto

    user, err := u.userRepo.GetUserByUsername(ctx, loginCredentialsDto.Username)
    if err != nil {
        return loginDto, errors.Wrap(common.ErrUserNotFound, err.Error())
    }

    err = util.CheckPassword(loginCredentialsDto.Password, user.HashedPassword)
    if err != nil {
        if err == sql.ErrNoRows {
            return loginDto, errors.Wrap(common.ErrUnauthorisedUser, err.Error())
        }
        return loginDto, err
    }

    accessToken, err := u.tokenMaker.CreateToken(user.Username, common.AccessTokenDuration)
    if err != nil {
        return loginDto, err
    }

    loginDto = dto.LoginUserDto{
        AccessToken: accessToken,
        User:        dto.NewUserDto(user),
    }

    return loginDto, nil
}

package service_test

import (
    "context"
    "fmt"
    "github.com/golang/mock/gomock"
    "github.com/pranayhere/simple-wallet/common"
    "github.com/pranayhere/simple-wallet/domain"
    "github.com/pranayhere/simple-wallet/dto"
    "github.com/pranayhere/simple-wallet/service"
    "github.com/pranayhere/simple-wallet/store"
    mockdb "github.com/pranayhere/simple-wallet/store/mock"
    "github.com/pranayhere/simple-wallet/token"
    "github.com/pranayhere/simple-wallet/util"
    "github.com/stretchr/testify/require"
    "reflect"
    "testing"
)

type eqCreateUserParamsMatcher struct {
    arg      store.CreateUserParams
    password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
    arg, ok := x.(store.CreateUserParams)
    if !ok {
        return false
    }

    err := util.CheckPassword(e.password, arg.HashedPassword)
    if err != nil {
        return false
    }

    e.arg.HashedPassword = arg.HashedPassword
    return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
    return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg store.CreateUserParams, password string) gomock.Matcher {
    return eqCreateUserParamsMatcher{arg, password}
}

func TestCreateUser(t *testing.T) {
    createUserDto := randomCreateUserDto()
    user, password := randomUser(t, createUserDto)

    testcases := []struct{
        name string
        reqDto dto.CreateUserDto
        buildStub func(mockUserRepo *mockdb.MockUserRepo)
        checkResp func(t *testing.T, userDto dto.UserDto, err error)
    }{
        {
            name: "OK",
            buildStub: func(mockUserRepo *mockdb.MockUserRepo) {
                arg := store.CreateUserParams{
                    Username: createUserDto.Username,
                    FullName: createUserDto.FullName,
                    Email: createUserDto.Email,
                    Status: domain.UserStatusACTIVE,
                }

                mockUserRepo.EXPECT().CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).Times(1).Return(user, nil)
            },
            checkResp: func(t *testing.T, userDto dto.UserDto, err error) {
                require.NoError(t, err)
                require.NotEmpty(t, userDto)
                require.Equal(t, createUserDto.Username, userDto.Username)
                require.Equal(t, createUserDto.FullName, userDto.FullName)
                require.Equal(t, createUserDto.Email, userDto.Email)
                require.Equal(t, string(domain.UserStatusACTIVE), userDto.Status)
            },
        },
    }

    for _, tc := range testcases {
        t.Run(tc.name, func(t *testing.T) {
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()

            mockUserRepo := mockdb.NewMockUserRepo(ctrl)
            tc.buildStub(mockUserRepo)

            tokenMaker, err := token.NewJWTMaker(common.SymmetricKey)
            require.NoError(t, err)

            ctx := context.TODO()
            userSvc := service.NewUserService(mockUserRepo, tokenMaker)
            userDto, err := userSvc.CreateUser(ctx, createUserDto)

            fmt.Println(userDto)
            fmt.Println(err)

            tc.checkResp(t, userDto, err)
        })
    }
}

func randomUser(t *testing.T, createUserDto dto.CreateUserDto) (user domain.User, password string) {
    password = createUserDto.Password
    hashedPassword, err := util.HashPassword(password)
    require.NoError(t, err)

    user = domain.User{
        Username:       createUserDto.Username,
        HashedPassword: hashedPassword,
        FullName:       createUserDto.FullName,
        Status:         domain.UserStatusACTIVE,
        Email:          createUserDto.Email,
    }
    return
}

func randomCreateUserDto() dto.CreateUserDto {
    createUserDto := dto.CreateUserDto{
        Username:       util.RandomUser(),
        Password:       util.RandomString(8),
        FullName:       util.RandomUser(),
        Email:          util.RandomEmail(),
    }

    return createUserDto
}
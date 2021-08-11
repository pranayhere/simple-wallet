package service_test

import (
    "context"
    "database/sql"
    "fmt"
    "github.com/golang/mock/gomock"
    "github.com/lib/pq"
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
    testcases := []struct {
        name      string
        reqDto    func() dto.CreateUserDto
        buildStub func(mockUserRepo *mockdb.MockUserRepo, createUserDto dto.CreateUserDto)
        checkResp func(t *testing.T, createUserDto dto.CreateUserDto, userDto dto.UserDto, err error)
    }{
        {
            name: "OK",
            reqDto: func() dto.CreateUserDto {
                return randomCreateUserDto()
            },
            buildStub: func(mockUserRepo *mockdb.MockUserRepo, createUserDto dto.CreateUserDto) {
                user, password := randomUser(t, createUserDto)

                arg := store.CreateUserParams{
                    Username: createUserDto.Username,
                    FullName: createUserDto.FullName,
                    Email:    createUserDto.Email,
                    Status:   domain.UserStatusACTIVE,
                }

                mockUserRepo.EXPECT().CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).Times(1).Return(user, nil)
            },
            checkResp: func(t *testing.T, createUserDto dto.CreateUserDto, userDto dto.UserDto, err error) {
                require.NoError(t, err)
                require.NotEmpty(t, userDto)
                require.Equal(t, createUserDto.Username, userDto.Username)
                require.Equal(t, createUserDto.FullName, userDto.FullName)
                require.Equal(t, createUserDto.Email, userDto.Email)
                require.Equal(t, string(domain.UserStatusACTIVE), userDto.Status)
            },
        },
        {
            name: "DatabaseConnectionClosed",
            reqDto: func() dto.CreateUserDto {
                return randomCreateUserDto()
            },
            buildStub: func(mockUserRepo *mockdb.MockUserRepo, createUserDto dto.CreateUserDto) {
                mockUserRepo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(1).Return(domain.User{}, sql.ErrConnDone)
            },
            checkResp: func(t *testing.T, createUserDto dto.CreateUserDto, userDto dto.UserDto, err error) {
                require.Error(t, err)
                require.EqualError(t, err, sql.ErrConnDone.Error())
            },
        },
        {
            name: "DuplicateUser",
            reqDto: func() dto.CreateUserDto {
                return randomCreateUserDto()
            },
            buildStub: func(mockUserRepo *mockdb.MockUserRepo, createUserDto dto.CreateUserDto) {
                // https://www.postgresql.org/docs/13/errcodes-appendix.html
                mockUserRepo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(1).Return(domain.User{}, &pq.Error{Code: "23505"})
            },
            checkResp: func(t *testing.T, createUserDto dto.CreateUserDto, userDto dto.UserDto, err error) {
                require.Error(t, err)
                require.EqualError(t, err, common.ErrUserAlreadyExist.Error())
            },
        },
    }

    for _, tc := range testcases {
        t.Run(tc.name, func(t *testing.T) {
            createUserDto := tc.reqDto()

            ctrl := gomock.NewController(t)
            defer ctrl.Finish()

            mockUserRepo := mockdb.NewMockUserRepo(ctrl)
            tc.buildStub(mockUserRepo, createUserDto)

            tokenMaker, err := token.NewJWTMaker(common.SymmetricKey)
            require.NoError(t, err)

            ctx := context.TODO()
            userSvc := service.NewUserService(mockUserRepo, tokenMaker)
            userDto, err := userSvc.CreateUser(ctx, createUserDto)

            tc.checkResp(t, createUserDto, userDto, err)
        })
    }
}

func TestLoginUser(t *testing.T) {
    createUserDto := randomCreateUserDto()
    user, password := randomUser(t, createUserDto)

    testcases := []struct {
        name      string
        reqDto    func() dto.LoginCredentialsDto
        buildStub func(mockUserRepo *mockdb.MockUserRepo, username string)
        checkResp func(t *testing.T, createUserDto dto.CreateUserDto, userDto dto.LoggedInUserDto, err error)
    }{
        {
            name: "Ok",
            reqDto: func() dto.LoginCredentialsDto {
                return dto.LoginCredentialsDto{
                    Username: user.Username,
                    Password: password,
                }
            },
            buildStub: func(mockUserRepo *mockdb.MockUserRepo, username string) {
                mockUserRepo.EXPECT().GetUserByUsername(gomock.Any(), username).Times(1).Return(user, nil)
            },
            checkResp: func(t *testing.T, createUserDto dto.CreateUserDto, loggedInUserDto dto.LoggedInUserDto, err error) {
                require.NoError(t, err)
                require.NotEmpty(t, loggedInUserDto)
                require.Equal(t, createUserDto.Username, loggedInUserDto.User.Username)
                require.Equal(t, createUserDto.Email, loggedInUserDto.User.Email)
                require.Equal(t, createUserDto.FullName, loggedInUserDto.User.FullName)
            },
        },
        {
            name: "UserNotFound",
            reqDto: func() dto.LoginCredentialsDto {
                return dto.LoginCredentialsDto{
                    Username: "Not Found",
                    Password: password,
                }
            },
            buildStub: func(mockUserRepo *mockdb.MockUserRepo, username string) {
                mockUserRepo.EXPECT().GetUserByUsername(gomock.Any(), gomock.Any()).Times(1).Return(user, sql.ErrNoRows)
            },
            checkResp: func(t *testing.T, createUserDto dto.CreateUserDto, loggedInUserDto dto.LoggedInUserDto, err error) {
                require.Error(t, err)
                require.EqualError(t, err, common.ErrUserNotFound.Error())
            },
        },
        {
            name: "IncorrectPassword",
            reqDto: func() dto.LoginCredentialsDto {
                return dto.LoginCredentialsDto{
                    Username: user.Username,
                    Password: "invalid",
                }
            },
            buildStub: func(mockUserRepo *mockdb.MockUserRepo, username string) {
                mockUserRepo.EXPECT().GetUserByUsername(gomock.Any(), username).Times(1).Return(user, nil)
            },
            checkResp: func(t *testing.T, createUserDto dto.CreateUserDto, loggedInUserDto dto.LoggedInUserDto, err error) {
                require.Error(t, err)
                require.EqualError(t, err, common.ErrIncorrectPassword.Error())
            },
        },
        {
            name: "DatabaseConnectionClosed",
            reqDto: func() dto.LoginCredentialsDto {
                return dto.LoginCredentialsDto{
                    Username: user.Username,
                    Password: password,
                }
            },
            buildStub: func(mockUserRepo *mockdb.MockUserRepo, username string) {
                mockUserRepo.EXPECT().GetUserByUsername(gomock.Any(), username).Times(1).Return(user, sql.ErrConnDone)
            },
            checkResp: func(t *testing.T, createUserDto dto.CreateUserDto, loggedInUserDto dto.LoggedInUserDto, err error) {
                require.Error(t, err)
                require.EqualError(t, err, sql.ErrConnDone.Error())
            },
        },
    }

    for _, tc := range testcases {
        t.Run(tc.name, func(t *testing.T) {
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()

            mockUserRepo := mockdb.NewMockUserRepo(ctrl)
            tc.buildStub(mockUserRepo, user.Username)

            tokenMaker, err := token.NewJWTMaker(common.SymmetricKey)
            require.NoError(t, err)

            ctx := context.TODO()
            userSvc := service.NewUserService(mockUserRepo, tokenMaker)

            arg := tc.reqDto()
            loggedInUserDto, err := userSvc.LoginUser(ctx, arg)

            tc.checkResp(t, createUserDto, loggedInUserDto, err)
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
        Username: util.RandomUser(),
        Password: util.RandomString(8),
        FullName: util.RandomUser(),
        Email:    util.RandomEmail(),
    }

    return createUserDto
}

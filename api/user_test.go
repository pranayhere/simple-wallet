package api_test

import (
    "bytes"
    "database/sql"
    "encoding/json"
    "fmt"
    "github.com/go-chi/chi"
    "github.com/golang/mock/gomock"
    "github.com/pranayhere/simple-wallet/api"
    "github.com/pranayhere/simple-wallet/dto"
    "github.com/pranayhere/simple-wallet/pkg/errors"
    mocksvc "github.com/pranayhere/simple-wallet/service/mock"
    "github.com/pranayhere/simple-wallet/util"
    "github.com/stretchr/testify/require"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestCreateUser(t *testing.T) {
    createUserDto := util.RandomCreateUserDto()
    userDto := util.RandomUserDto(createUserDto)

    testcases := []struct{
        name string
        body map[string]interface{}
        buildStub func(mockUserSvc *mocksvc.MockUserSvc)
        checkRes func(recorder *httptest.ResponseRecorder)
    }{
        {
            name: "Ok",
            body: map[string]interface{}{
                "username": createUserDto.Username,
                "password": createUserDto.Password,
                "full_name": createUserDto.FullName,
                "email": createUserDto.Email,
            },
            buildStub: func(mockUserSvc *mocksvc.MockUserSvc) {
                mockUserSvc.EXPECT().CreateUser(gomock.Any(), createUserDto).Times(1).Return(userDto, nil)
            },
            checkRes: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusOK, recorder.Code)
            },
        },
        {
            name: "JsonInvalid",
            body: map[string]interface{}{
                "username": 1,
                "password": createUserDto.Password,
                "full_name": createUserDto.FullName,
                "email": createUserDto.Email,
            },
            buildStub: func(mockUserSvc *mocksvc.MockUserSvc) {
                mockUserSvc.EXPECT().CreateUser(gomock.Any(), createUserDto).Times(0)
            },
            checkRes: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusBadRequest, recorder.Code)
                fmt.Println(recorder.Body)
            },
        },
        {
            name: "InternalErr",
            body: map[string]interface{}{
                "username": createUserDto.Username,
                "password": createUserDto.Password,
                "full_name": createUserDto.FullName,
                "email": createUserDto.Email,
            },
            buildStub: func(mockUserSvc *mocksvc.MockUserSvc) {
                mockUserSvc.EXPECT().CreateUser(gomock.Any(), createUserDto).Times(1).Return(dto.UserDto{}, sql.ErrConnDone)
            },
            checkRes: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusInternalServerError, recorder.Code)
            },
        },
        {
            name: "DuplicateUser",
            body: map[string]interface{}{
                "username": createUserDto.Username,
                "password": createUserDto.Password,
                "full_name": createUserDto.FullName,
                "email": createUserDto.Email,
            },
            buildStub: func(mockUserSvc *mocksvc.MockUserSvc) {
                mockUserSvc.EXPECT().CreateUser(gomock.Any(), createUserDto).Times(1).Return(dto.UserDto{}, errors.ErrUserAlreadyExist)
            },
            checkRes: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusForbidden, recorder.Code)
            },
        },
        {
            name: "InvalidUsername",
            body: map[string]interface{}{
                "username": "invalid-user#1",
                "password": createUserDto.Password,
                "full_name": createUserDto.FullName,
                "email": createUserDto.Email,
            },
            buildStub: func(mockUserSvc *mocksvc.MockUserSvc) {
                mockUserSvc.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
            },
            checkRes: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusBadRequest, recorder.Code)
            },
        },
        {
            name: "InvalidEmail",
            body: map[string]interface{}{
                "username": createUserDto.Username,
                "password": createUserDto.Password,
                "full_name": createUserDto.FullName,
                "email": "invalid-email",
            },
            buildStub: func(mockUserSvc *mocksvc.MockUserSvc) {
                mockUserSvc.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
            },
            checkRes: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusBadRequest, recorder.Code)
            },
        },
        {
            name: "TooShortPassword",
            body: map[string]interface{}{
                "username": createUserDto.Username,
                "password": "123",
                "full_name": createUserDto.FullName,
                "email": createUserDto.Email,
            },
            buildStub: func(mockUserSvc *mocksvc.MockUserSvc) {
                mockUserSvc.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)
            },
            checkRes: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusBadRequest, recorder.Code)
            },
        },
    }

    for _, tc := range testcases {
        t.Run(tc.name, func(t *testing.T) {
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()

            mockUserSvc := mocksvc.NewMockUserSvc(ctrl)
            tc.buildStub(mockUserSvc)

            recorder := httptest.NewRecorder()
            router := chi.NewRouter()

            userApi := api.NewUserResource(mockUserSvc)
            userApi.RegisterRoutes(router)

            data, err := json.Marshal(tc.body)
            require.NoError(t, err)

            url := "/users"
            request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
            require.NoError(t, err)

            router.ServeHTTP(recorder, request)

            tc.checkRes(recorder)
        })
    }
}

func TestLoginUser(t *testing.T) {
    createUserDto := util.RandomCreateUserDto()
    userDto := util.RandomUserDto(createUserDto)

    loggedInUserDto := dto.LoggedInUserDto{
        AccessToken: util.RandomString(10),
        User: userDto,
    }


    testcases := []struct{
        name string
        body map[string]interface{}
        buildStub func(mockUserSvc *mocksvc.MockUserSvc)
        checkRes func(recorder *httptest.ResponseRecorder)
    }{
        {
            name: "Ok",
            body: map[string]interface{}{
                "username": createUserDto.Username,
                "password": createUserDto.Password,
            },
            buildStub: func(mockUserSvc *mocksvc.MockUserSvc) {
                loginDto := dto.LoginCredentialsDto{
                    Username: createUserDto.Username,
                    Password: createUserDto.Password,
                }
                mockUserSvc.EXPECT().LoginUser(gomock.Any(), loginDto).Times(1).Return(loggedInUserDto, nil)
            },
            checkRes: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusOK, recorder.Code)
            },
        },
        {
            name: "JsonInvalid",
            body: map[string]interface{}{
                "username": 1,
                "password": createUserDto.Password,
            },
            buildStub: func(mockUserSvc *mocksvc.MockUserSvc) {
                mockUserSvc.EXPECT().LoginUser(gomock.Any(), gomock.Any()).Times(0)
            },
            checkRes: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusBadRequest, recorder.Code)
            },
        },
        {
            name: "UserNotFound",
            body: map[string]interface{}{
                "username": createUserDto.Username,
                "password": createUserDto.Password,
            },
            buildStub: func(mockUserSvc *mocksvc.MockUserSvc) {
                mockUserSvc.EXPECT().LoginUser(gomock.Any(), gomock.Any()).Times(1).Return(dto.LoggedInUserDto{}, errors.ErrUserNotFound)
            },
            checkRes: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusNotFound, recorder.Code)
            },
        },
        {
            name: "InternalError",
            body: map[string]interface{}{
                "username": createUserDto.Username,
                "password": createUserDto.Password,
            },
            buildStub: func(mockUserSvc *mocksvc.MockUserSvc) {
                mockUserSvc.EXPECT().LoginUser(gomock.Any(), gomock.Any()).Times(1).Return(dto.LoggedInUserDto{}, errors.ErrUserNotFound)
            },
            checkRes: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusNotFound, recorder.Code)
            },
        },
        {
            name: "InvalidUsername",
            body: map[string]interface{}{
                "username": "invalid-user#1",
                "password": createUserDto.Password,
            },
            buildStub: func(mockUserSvc *mocksvc.MockUserSvc) {
                mockUserSvc.EXPECT().LoginUser(gomock.Any(), gomock.Any()).Times(0)
            },
            checkRes: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusBadRequest, recorder.Code)
            },
        },
        {
            name: "TooShortPassword",
            body: map[string]interface{}{
                "username": createUserDto.Username,
                "password": "123",
            },
            buildStub: func(mockUserSvc *mocksvc.MockUserSvc) {
                mockUserSvc.EXPECT().LoginUser(gomock.Any(), gomock.Any()).Times(0)
            },
            checkRes: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusBadRequest, recorder.Code)
            },
        },
    }

    for _, tc := range testcases {
        t.Run(tc.name, func(t *testing.T) {
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()

            mockUserSvc := mocksvc.NewMockUserSvc(ctrl)
            tc.buildStub(mockUserSvc)

            recorder := httptest.NewRecorder()
            router := chi.NewRouter()

            userApi := api.NewUserResource(mockUserSvc)
            userApi.RegisterRoutes(router)

            data, err := json.Marshal(tc.body)
            require.NoError(t, err)

            url := "/users/login"
            request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
            require.NoError(t, err)

            router.ServeHTTP(recorder, request)
            tc.checkRes(recorder)
        })
    }
}
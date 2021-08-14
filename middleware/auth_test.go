package middleware_test

import (
    "fmt"
    "github.com/go-chi/chi"
    "github.com/go-chi/render"
    "github.com/pranayhere/simple-wallet/middleware"
    "github.com/pranayhere/simple-wallet/pkg/constant"
    "github.com/pranayhere/simple-wallet/token"
    "github.com/stretchr/testify/require"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"
)

func AddAuthorization(t *testing.T, request *http.Request, tokenMaker token.Maker, authorizationType string, userID int64, duration time.Duration) {
    accessToken, err := tokenMaker.CreateToken(userID, duration)
    require.NoError(t, err)

    authorizationHeader := fmt.Sprintf("%s %s", authorizationType, accessToken)
    request.Header.Set(constant.AuthorizationHeaderKey, authorizationHeader)
}

func TestAuthMiddleware(t *testing.T) {
    testCases := []struct {
        name          string
        setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
        checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
    }{
        {
            name: "OK",
            setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
                AddAuthorization(t, request, tokenMaker, constant.AuthorizationTypeBearer, 1, time.Minute)
            },
            checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusOK, recorder.Code)
            },
        },
        {
            name: "NoAuthorization",
            setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
            },
            checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusUnauthorized, recorder.Code)
            },
        },
        {
            name: "UnsupportedAuthorization",
            setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
                AddAuthorization(t, request, tokenMaker, "unsupported", 1, time.Minute)
            },
            checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusUnauthorized, recorder.Code)
            },
        },
        {
            name: "InvalidAuthorizationFormat",
            setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
                AddAuthorization(t, request, tokenMaker, "", 1, time.Minute)
            },
            checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusUnauthorized, recorder.Code)
            },
        },
        {
            name: "ExpiredToken",
            setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
                AddAuthorization(t, request, tokenMaker, constant.AuthorizationTypeBearer, 1, -time.Minute)
            },
            checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusUnauthorized, recorder.Code)
            },
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            r := chi.NewRouter()

            tokenMaker, err := token.NewJWTMaker(constant.SymmetricKey)
            require.NoError(t, err)
            require.NotEmpty(t, tokenMaker)

            authPath := "/auth"
            r.With(middleware.Auth(tokenMaker)).Get(
                authPath,
                func(w http.ResponseWriter, r *http.Request) {
                    render.JSON(w, r, "Ok")
                })

            recorder := httptest.NewRecorder()
            request, err := http.NewRequest(http.MethodGet, authPath, nil)
            require.NoError(t, err)

            tc.setupAuth(t, request, tokenMaker)
            r.ServeHTTP(recorder, request)
            tc.checkResponse(t, recorder)
        })
    }
}

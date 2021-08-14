package main

import (
    "database/sql"
    "github.com/go-chi/chi"
    "github.com/go-chi/chi/middleware"
    "github.com/go-chi/httprate"
    "github.com/go-chi/render"
    "github.com/pranayhere/simple-wallet/api"
    middleware2 "github.com/pranayhere/simple-wallet/middleware"
    "github.com/pranayhere/simple-wallet/pkg/constant"
    "github.com/pranayhere/simple-wallet/service"
    "github.com/pranayhere/simple-wallet/store"
    "github.com/pranayhere/simple-wallet/token"
    "net/http"
    "time"
)

func createRouter() *chi.Mux {
    r := chi.NewRouter()

    r.Use(middleware.RequestID)
    r.Use(middleware.RealIP)
    r.Use(middleware.Recoverer)
    r.Use(middleware.Logger)
    r.Use(httprate.LimitByIP(100, 1*time.Minute))

    return r
}

func initRoutes(db *sql.DB, r *chi.Mux) *chi.Mux {
    currencyRepo := store.NewCurrencyRepo(db)
    currencySvc := service.NewCurrencyService(currencyRepo)
    currencyApi := api.NewCurrencyResource(currencySvc)

    tokenMaker, err := token.NewJWTMaker(constant.SymmetricKey)
    if err != nil {
        panic(err)
    }

    userRepo := store.NewUserRepo(db)
    userSvc := service.NewUserService(userRepo, tokenMaker)
    userApi := api.NewUserResource(userSvc)

    transferRepo := store.NewTransferRepo(db)
    entryRepo := store.NewEntryRepo(db)
    walletRepo := store.NewWalletRepo(db, transferRepo, entryRepo)
    bankAccountRepo := store.NewBankAccountRepo(db, walletRepo, userRepo)
    bankAcctSvc := service.NewBankAccountService(bankAccountRepo, currencySvc)
    bankAcctApi := api.NewBankAccountResource(bankAcctSvc)

    walletSvc := service.NewWalletService(walletRepo)
    walletApi := api.NewWalletResource(walletSvc)

    // Routes
    // public
    userApi.RegisterRoutes(r.With(httprate.LimitByIP(10, 1*time.Minute)))

    // authorized
    r.Group(func(r chi.Router) {
        r.Use(middleware2.Auth(tokenMaker))
        bankAcctApi.RegisterRoutes(r)
        currencyApi.RegisterRoutes(r)
        walletApi.RegisterRoutes(r)
    })

    r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
        render.JSON(w, r, "ok")
    })

    return r
}

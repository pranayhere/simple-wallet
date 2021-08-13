package main

import (
    "database/sql"
    "github.com/go-chi/chi"
    "github.com/go-chi/render"
    "github.com/pranayhere/simple-wallet/api"
    m "github.com/pranayhere/simple-wallet/middleware"
    "github.com/pranayhere/simple-wallet/pkg/constant"
    "github.com/pranayhere/simple-wallet/service"
    "github.com/pranayhere/simple-wallet/store"
    "github.com/pranayhere/simple-wallet/token"
    "net/http"
)

func inject(db *sql.DB, r *chi.Mux) *chi.Mux {
    // currency
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

    r.Mount("/users", userApi.RegisterRoutes(r))

    r.Group(func(r chi.Router) {
        r.Use(m.Auth(tokenMaker))
        r.Mount("/bank-accounts", bankAcctApi.RegisterRoutes(r))
        r.Mount("/currencies", currencyApi.RegisterRoutes(r))
        r.Mount("/wallets", walletApi.RegisterRoutes(r))
    })

    r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
        render.JSON(w, r, "ok")
    })

    return r
}

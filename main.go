package main

import (
    "context"
    "database/sql"
    "fmt"
    "github.com/go-chi/chi"
    "github.com/go-chi/chi/middleware"
    "github.com/pranayhere/simple-wallet/api"
    "github.com/pranayhere/simple-wallet/common"
    "github.com/pranayhere/simple-wallet/service"
    "github.com/pranayhere/simple-wallet/store"
    "github.com/pranayhere/simple-wallet/token"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)

const (
    dbDriver      = "postgres"
    dbSource      = "postgresql://root:secret@localhost:5555/simple_wallet?sslmode=disable"
    serverAddress = "localhost:8080"
)

func main() {
    log.Println("starting wallet service")

    db := NewStore()
    defer db.Close()

    r := CreateRouter()

    // currency
    currencyRepo := store.NewCurrencyRepo(db)
    currencySvc := service.NewCurrencyService(currencyRepo)
    currencyApi := api.NewCurrencyResource(currencySvc)
    currencyApi.RegisterRoutes(r)

    tokenMaker, err := token.NewJWTMaker(common.SymmetricKey)
    if err != nil {
        panic(err)
    }

    userRepo := store.NewUserRepo(db)
    userSvc := service.NewUserService(userRepo, tokenMaker)
    userApi := api.NewUserResource(userSvc)
    userApi.RegisterRoutes(r)

    transferRepo := store.NewTransferRepo(db)
    entryRepo := store.NewEntryRepo(db)
    walletRepo := store.NewWalletRepo(db, transferRepo, entryRepo)
    bankAccountRepo := store.NewBankAccountRepo(db, walletRepo, userRepo)
    bankAcctSvc := service.NewBankAccountService(bankAccountRepo, currencySvc)
    bankAcctApi := api.NewBankAccountResource(bankAcctSvc)
    bankAcctApi.RegisterRoutes(r)

    walletSvc := service.NewWalletService(walletRepo)
    walletApi := api.NewWalletResource(walletSvc)
    walletApi.RegisterRoutes(r)

    r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte(fmt.Sprintf("Sup!!!")))
    })

    server := &http.Server{Addr: serverAddress, Handler: r}
    serverCtx, serverStopCtx := context.WithCancel(context.Background())

    // Listen for syscall signals for process to interrupt/quit
    sig := make(chan os.Signal, 1)
    signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
    go func() {
        <-sig

        // Shutdown signal with grace period of 30 seconds
        shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

        go func() {
            <-shutdownCtx.Done()
            if shutdownCtx.Err() == context.DeadlineExceeded {
                log.Fatal("graceful shutdown timed out.. forcing exit.")
            }
        }()

        // Trigger graceful shutdown
        err := server.Shutdown(shutdownCtx)
        if err != nil {
            log.Fatal(err)
        }
        serverStopCtx()
    }()

    // Run the server
    err = server.ListenAndServe()
    if err != nil && err != http.ErrServerClosed {
        log.Fatal(err)
    }

    // Wait for server context to be stopped
    <-serverCtx.Done()
}

func NewStore() *sql.DB {
    conn, err := sql.Open(dbDriver, dbSource)
    if err != nil {
        log.Fatal("cannot connect to db", err)
    }

    return conn
}

func CreateRouter() *chi.Mux {
    r := chi.NewRouter()

    r.Use(middleware.Recoverer)
    r.Use(middleware.Logger)

    return r
}


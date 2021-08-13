package main

import (
    "context"
    "database/sql"
    "github.com/go-chi/chi"
    "github.com/go-chi/chi/middleware"
    "go.uber.org/zap"
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
    r = inject(db, r)

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
            log.Fatal("", zap.Error(err))
        }
        serverStopCtx()
    }()

    // Run the server
    err := server.ListenAndServe()
    if err != nil && err != http.ErrServerClosed {
        log.Fatal("failed to start server", zap.Error(err))
        os.Exit(1)
    }

    // Wait for server context to be stopped
    <-serverCtx.Done()
}

func NewStore() *sql.DB {
    conn, err := sql.Open(dbDriver, dbSource)
    if err != nil {
        log.Fatal("cannot connect to db", zap.Error(err))
    }

    return conn
}

func CreateRouter() *chi.Mux {
    r := chi.NewRouter()

    r.Use(middleware.Recoverer)
    r.Use(middleware.Logger)

    return r
}
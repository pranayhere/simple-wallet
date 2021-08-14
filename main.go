package main

import (
    "context"
    "database/sql"
    log "github.com/sirupsen/logrus"
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
    log.SetFormatter(&log.TextFormatter{})

    log.Println("starting wallet service")

    db := newStore()
    defer db.Close()

    r := createRouter()
    r = initRoutes(db, r)

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
            log.Fatal("", err)
        }
        serverStopCtx()
    }()

    // Run the server
    err := server.ListenAndServe()
    if err != nil && err != http.ErrServerClosed {
        log.Fatal("failed to start server", err)
        os.Exit(1)
    }

    // Wait for server context to be stopped
    <-serverCtx.Done()
}

func newStore() *sql.DB {
    conn, err := sql.Open(dbDriver, dbSource)
    if err != nil {
        log.Fatal("cannot connect to db", err)
    }

    return conn
}

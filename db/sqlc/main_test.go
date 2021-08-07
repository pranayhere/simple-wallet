package db_test

import (
    "database/sql"
    db "github.com/pranayhere/simple-wallet/db/sqlc"
    "log"
    "os"
    "testing"

    _ "github.com/lib/pq"
)

const (
    dbDriver = "postgres"
    dbSource = "postgresql://root:secret@localhost:5555/simple_wallet?sslmode=disable"
)

var testQueries *db.Queries
var testDb *sql.DB

func TestMain(m *testing.M) {
    var err error

    testDb, err = sql.Open(dbDriver, dbSource)
    if err != nil {
        log.Fatal("cannot connect to db")
    }

    testQueries = db.New(testDb)

    os.Exit(m.Run())
}
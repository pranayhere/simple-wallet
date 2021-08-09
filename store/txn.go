package store

import "database/sql"

type Tx interface {
    Exec(query string, args ...interface{}) (sql.Result, error)
    Prepare(query string) (*sql.Stmt, error)
    Query(query string, args ...interface{}) (*sql.Rows, error)
    QueryRow(query string, args ...interface{}) *sql.Row
}

var TxKey = struct{}{}

// TxFn is a function that will be called with an initialized `Tx` object
// that can be used for executing statements and queries against a database.
type TxFn func(Tx) error

// ExecTx creates a new transaction and handles rollback/commit based on the
// error object returned by the `TxFn`
func ExecTx(db *sql.DB, fn TxFn) (err error) {
    tx, err := db.Begin()
    if err != nil {
        return
    }

    defer func() {
        if p := recover(); p != nil {
            // a panic occurred, rollback and repanic
            tx.Rollback()
            panic(p)
        } else if err != nil {
            // something went wrong, rollback
            tx.Rollback()
        } else {
            // all good, commit
            err = tx.Commit()
        }
    }()

    err = fn(tx)
    return err
}

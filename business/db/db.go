// Package db provides support for access the database.
package db

import (
	"database/sql"
	"embed"
)

//go:embed *
var DB embed.FS

//go:generate mockery --name Transaction

// Transaction is the transaction interface for database handler.
type Transaction interface {
	Rollback() error
	Commit() error
	TxEnd(txFunc func() error) error
	TxBegin() (SQLGbc, error)
}

//go:generate mockery --name SQLGbc

// SQLGbc (SQL Go database connection) is a wrapper for SQL database handler ( can be *sql.DB or *sql.Tx)
// It SQLGbc be able to work with all SQL data that follows SQL standard.
type SQLGbc interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Transaction
}

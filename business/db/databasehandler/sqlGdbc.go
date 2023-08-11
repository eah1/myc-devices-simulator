// Package databasehandler low level database access including transaction through *sql.Tx or *sql.DB
package databasehandler

import (
	"database/sql"
	"fmt"
)

// SQLDBTx is the concrete implementation of sqlGdbc by using *sql.DB.
type SQLDBTx struct {
	DB *sql.DB
}

// SQLConnTx is the concrete implementation of sqlGdbc by using *sql.Tx.
type SQLConnTx struct {
	DB *sql.Tx
}

// Exec executes a query without returning any rows.
func (sdt *SQLDBTx) Exec(query string, args ...interface{}) (sql.Result, error) {
	exec, err := sdt.DB.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("databasehandler.sqlGdbc.Exec.Exec(%v, %v): %w", query, args, err)
	}

	return exec, nil
}

// Prepare creates a prepared statement for later queries or executions.
func (sdt *SQLDBTx) Prepare(query string) (*sql.Stmt, error) {
	prepare, err := sdt.DB.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("databasehandler.sqlGdbc.Prepare.Prepare(%v): %w", query, err)
	}

	return prepare, nil
}

// Exec executes a query without returning any rows.
func (sdb *SQLConnTx) Exec(query string, args ...interface{}) (sql.Result, error) {
	exec, err := sdb.DB.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("databasehandler.sqlGdbc.Exec.Exec(%v, %v): %w", query, args, err)
	}

	return exec, nil
}

// Prepare creates a prepared statement for later queries or executions.
func (sdb *SQLConnTx) Prepare(query string) (*sql.Stmt, error) {
	prepare, err := sdb.DB.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("databasehandler.sqlGdbc.Prepare.Prepare(%v): %w", query, err)
	}

	return prepare, nil
}

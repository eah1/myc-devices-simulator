package databasehandler

import (
	"fmt"
	"myc-devices-simulator/business/db"
)

// Rollback transaction rollback.
func (sdt *SQLDBTx) Rollback() error {
	return nil
}

// Commit transaction commit.
func (sdt *SQLDBTx) Commit() error {
	return nil
}

// TxBegin transaction begin.
func (sdt *SQLDBTx) TxBegin() (db.SQLGbc, error) {
	tx, err := sdt.DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("databasehandler.transactor.TxBegin.Begin(-): %w", err)
	}

	return &SQLConnTx{DB: tx}, nil
}

// TxEnd transaction end.
func (sdt *SQLDBTx) TxEnd(_ func() error) error {
	return nil
}

// Rollback transaction rollback.
func (sdb *SQLConnTx) Rollback() error {
	if err := sdb.DB.Rollback(); err != nil {
		return fmt.Errorf("databasehandler.transactor.Rollback.Rollback(-): %w", err)
	}

	return nil
}

// Commit transaction commit.
func (sdb *SQLConnTx) Commit() error {
	if err := sdb.DB.Commit(); err != nil {
		return fmt.Errorf("databasehandler.transactor.Commit.Commit(-): %w", err)
	}

	return nil
}

// TxBegin transaction begin.
func (sdb *SQLConnTx) TxBegin() (db.SQLGbc, error) {
	return nil, nil
}

// TxEnd transaction end.
func (sdb *SQLConnTx) TxEnd(txFunc func() error) error {
	var err error

	tx := sdb.DB

	defer func() {
		if p := recover(); p != nil {
			if err := tx.Rollback(); err != nil {
				return
			}

			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			if err := tx.Rollback(); err != nil {
				return
			} // err is non-nil; don't change it
		} else {
			err = tx.Commit() // if Commit returns error update err with commit err
		}
	}()

	err = txFunc()

	return err
}

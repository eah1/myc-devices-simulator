// Package config.go provides support for access the database connection.
package config

import (
	"database/sql"
	"fmt"
	"myc-devices-simulator/business/sys/errors"
	"myc-devices-simulator/cmd/config"
	"time"

	_ "github.com/lib/pq"
)

const maxConnected = 500

// Config is the required properties to use the database.
type Config struct {
	DBPostgres  string
	MaxIdleCons int
	MaxOpenCons int
	ShowSQL     bool
	Timezone    string
}

// NewConfig create a Config Data Base.
func NewConfig(config config.Config) Config {
	configDB := new(Config)
	configDB.DBPostgres = config.DBPostgres
	configDB.MaxIdleCons = config.DBMaxIdleConns
	configDB.MaxOpenCons = config.DBMaxOpenConns
	configDB.ShowSQL = config.DBLogger
	configDB.Timezone = "UTC"

	return *configDB
}

// Open knows how to open a DB connection based on the configuration.
func Open(cfg Config, retries int) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DBPostgres)
	if err != nil {
		return nil, fmt.Errorf("database.sys.OpenDB.NewEngine: %w", errors.ErrPsqlConnection)
	}

	for i := 1; i <= retries; i++ {
		err = db.Ping()
		if err == nil {
			break
		}

		//nolint:forbidigo
		fmt.Printf("Failed to connect to database (attempt %d)\n", i)
		time.Sleep(time.Duration(maxConnected*i) * time.Millisecond)
	}

	if err != nil {
		return nil, fmt.Errorf("database.sys.OpenDB.Ping: %w", errors.ErrPsqlConnection)
	}

	if locUTC, err := time.LoadLocation("Etc/UTC"); err == nil {
		time.Local = locUTC
	} else {
		return nil, fmt.Errorf("database.sys.OpenDB.LoadLocation: %w", errors.ErrPsqlConnection)
	}

	locUTC, err := time.LoadLocation("Etc/UTC")
	if err != nil {
		return nil, fmt.Errorf("database.sys.Open.LoadLocation(-) mycError: {%w}", errors.ErrPsqlConnection)
	}

	time.Local = locUTC

	db.SetMaxIdleConns(cfg.MaxIdleCons)
	db.SetMaxOpenConns(cfg.MaxOpenCons)

	return db, nil
}

// Package config content struct Config. The function will parse the environment variables to run the service.
package config

type Config struct {
	DBPostgres     string
	DBMaxIdleConns int
	DBMaxOpenConns int
	DBLogger       bool
}
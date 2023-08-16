package test

import (
	"database/sql"
	"io"
	"myc-devices-simulator/business/db"
	dbconfig "myc-devices-simulator/business/db/config"
	"myc-devices-simulator/business/sys/logger"
	"myc-devices-simulator/cmd/config"
	"myc-devices-simulator/foundation/docker"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var (
	lock     = &sync.Mutex{}
	database *sql.DB
)

const retryDBMax = 5

// InitLogger create logger.
func InitLogger(t *testing.T, service string) *zap.SugaredLogger {
	t.Helper()

	log, err := logger.InitLogger(service, "test")
	require.NoError(t, err)

	return log
}

// InitDockerContainerDatabase create container postgres.
func InitDockerContainerDatabase(t *testing.T, testName string) *docker.Container {
	image := "postgres"
	port := "5432"
	dockerArgs := []string{"--name", testName, "-e", "POSTGRES_PASSWORD=postgres"}
	appArgs := []string{"-c", "log_statement=all"}

	container, err := docker.StartContainer(image, port, dockerArgs, appArgs)
	require.NoError(t, err)

	t.Logf("Image:       %s\n", image)
	t.Logf("ContainerID: %s\n", container.ID)
	t.Logf("Host:        %s\n", container.Host)

	// Give Vault time to initialize.
	time.Sleep(time.Second)

	return container
}

// InitConfig create a config.go env.
func InitConfig(host string) config.Config {
	configENV := new(config.Config)
	configENV.DBPostgres = "postgres://postgres:postgres@" + host + "/postgres?sslmode=disable"
	configENV.DBMaxOpenConns = 25
	configENV.DBMaxIdleConns = 25
	configENV.DBLogger = false

	return *configENV
}

// InitDatabase create database.
func InitDatabase(t *testing.T, config config.Config, log *zap.SugaredLogger) *sql.DB {
	t.Helper()

	var err error

	if database == nil {
		lock.Lock()
		defer lock.Unlock()

		if database == nil {
			database, err = dbconfig.Open(dbconfig.NewConfig(config), retryDBMax)
			require.NoError(t, err)

			dir, err := db.DB.ReadDir("schema")
			require.NoError(t, err)

			err = os.Mkdir("new_schema", os.ModePerm)
			require.NoError(t, err)

			for _, files := range dir {
				srcFile, err := db.DB.Open("schema/" + files.Name())
				require.NoError(t, err)
				defer srcFile.Close()

				dstFile, err := os.Create(filepath.Join("new_schema", files.Name()))
				require.NoError(t, err)
				defer dstFile.Close()

				_, err = io.Copy(dstFile, srcFile)
				require.NoError(t, err)
			}

			err = goose.Up(database, "new_schema")
			require.NoError(t, err)
		}
	}

	return database
}

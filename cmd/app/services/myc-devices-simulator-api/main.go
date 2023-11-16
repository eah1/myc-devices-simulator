package main

import (
	"database/sql"
	"fmt"
	"myc-devices-simulator/business/core/email"
	userCore "myc-devices-simulator/business/core/user"
	dbConfig "myc-devices-simulator/business/db/config"
	"myc-devices-simulator/business/db/databasehandler"
	v1 "myc-devices-simulator/business/infra/handlers/v1"
	userHandler "myc-devices-simulator/business/infra/handlers/v1/user"
	userStore "myc-devices-simulator/business/repository/store/user"
	emailsender "myc-devices-simulator/business/sys/email_sender"
	"myc-devices-simulator/business/sys/logger"
	userUC "myc-devices-simulator/business/usecase/user"
	"myc-devices-simulator/cmd/config"
	"os"
	"strings"

	"github.com/jhillyerd/enmime"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	retryDBMax   = 5
	echoLogLevel = 2
)

func main() {
	// Create config env vars.
	cfg, err := config.Load()
	if err != nil {
		//nolint: forbidigo
		fmt.Println("Cannot crate LoadConfig: " + err.Error())

		os.Exit(1)
	}

	// Construct the application logger.
	log, err := logger.InitLogger("MYC-DEVICES-SIMULATOR-API", cfg.Environment)
	if err != nil {
		//nolint: forbidigo
		fmt.Println("Cannot crate Logger: " + err.Error())

		os.Exit(1)
	}

	log.Infow("starting environments status",
		"host", cfg.Host, "hostName", cfg.HostName, "port", cfg.Port,
		"base url", cfg.BaseURL, "server URI", cfg.ServerURI, "environment", cfg.Environment)

	// Perform the startup and shutdown sequence.
	if err := run(log, cfg); err != nil {
		log.Errorw("startup", "ERROR", err)

		os.Exit(1)
	}

	defer func(log *zap.SugaredLogger) {
		if err := log.Sync(); err != nil {
			log.Error(err)
		}
	}(log)
}

// run init app.
func run(log *zap.SugaredLogger, cfg config.Config) error {
	log.Infow("startup")

	// Create connectivity to the database.
	host := cfg.DBPostgres[strings.Index(cfg.DBPostgres, "@")+1 : strings.LastIndex(cfg.DBPostgres, "/")]

	database, err := dbConfig.Open(dbConfig.NewConfig(cfg), retryDBMax)
	if err != nil {
		log.Errorf("main.run.db.open:: %s", err)

		return fmt.Errorf("main.run.db.open: %w", err)
	}

	log.Infow("starting database status", "host", host)

	defer func() {
		log.Infow("shutdown - stopping database support", "host", host)

		if err := database.Close(); err != nil {
			log.Error(err)
		}
	}()

	// Created a configuration email.
	emailConfig, err := emailsender.InnitEmailConfig(cfg)
	if err != nil {
		log.Errorf("email configuration: %w", err)

		return fmt.Errorf("email config: %w", err)
	}

	log.Infow("starting email sender status")

	// start services.
	log.Errorf("%s", startEcho(log, cfg, database, emailConfig))

	return nil
}

// startEcho start server.
func startEcho(log *zap.SugaredLogger, cfg config.Config, database *sql.DB, emailConfig *enmime.SMTPSender) error {
	// Start App
	app := echo.New()

	// hide echo banner.
	app.HideBanner = false

	// Set logging level to INFO.
	app.Logger.SetLevel(echoLogLevel)

	// start root group api.
	root := app.Group("/api/")
	groupV1 := v1.CreateGroupV1(root)

	// set binder custom.
	// app.Binder = &binder.CustomBinder{}

	// start database.

	// start stores.
	storeUser := userStore.NewUserStore(&databasehandler.SQLDBTx{DB: database}, log)

	// start cores.
	coreEmail := email.NewEmailCore(emailConfig, cfg.SMTPFrom, log, cfg)
	coreUser := userCore.NewUserCore(&storeUser)

	// start use case.
	ucUserRegister := userUC.NewUCUserRegister(&coreUser, &coreEmail)

	// Handlers components.
	handlerUser := userHandler.NewHandlerUser(ucUserRegister)

	// Group components.
	userHandler.GroupUser(groupV1, "/users", handlerUser)

	return errors.Wrap(app.Start(cfg.Host+":"+cfg.Port), "error in server start")
}

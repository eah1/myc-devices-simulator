package user

import (
	"context"
	"fmt"
	"myc-devices-simulator/business/db"
	"myc-devices-simulator/business/db/errors"
	errorssys "myc-devices-simulator/business/sys/errors"

	"go.uber.org/zap"
)

const InsertUser = "INSERT INTO users (id, first_name, last_name, email, password, language, company) VALUES " +
	"($1, $2, $3, $4, $5, $6, $7)"

type UserStore struct {
	db  db.SQLGbc
	log *zap.SugaredLogger
}

func NewUserStore(database db.SQLGbc, log *zap.SugaredLogger) UserStore {
	return UserStore{
		db:  database,
		log: log,
	}
}

func (store *UserStore) InsertUser(_ context.Context, user User) error {
	prepare, err := store.db.Prepare(InsertUser)
	if err != nil {
		return fmt.Errorf("store.user.InsertUser.Prepare: %w", errors.WrapperError(store.log, err))
	}

	res, err := prepare.Exec(user.ID, user.FirstName, user.LastName,
		user.Email, user.Password, user.Language, user.Company)
	if err != nil {
		return fmt.Errorf("store.user.InsertUser.Exec: %w", errors.WrapperError(store.log, err))
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("store.user.InsertUser.RowsAffected: %w", errors.WrapperError(store.log, err))
	}

	if affected == 0 {
		return fmt.Errorf("store.user.InsertUser: %w", errorssys.ErrRowAffected)
	}

	return nil
}

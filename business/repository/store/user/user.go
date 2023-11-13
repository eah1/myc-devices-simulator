// Package user call database layer.
package user

import (
	"context"
	"fmt"
	"myc-devices-simulator/business/db"
	errorssys "myc-devices-simulator/business/sys/errors"

	"github.com/lib/pq"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const InsertUser = "INSERT INTO users (id, first_name, last_name, email, password, language, company) VALUES " +
	"($1, $2, $3, $4, $5, $6, $7)"

// StoreUser store user to call database.
type StoreUser struct {
	db  db.SQLGbc
	log *zap.SugaredLogger
}

// NewUserStore construct a store user group.
func NewUserStore(database db.SQLGbc, log *zap.SugaredLogger) StoreUser {
	return StoreUser{
		db:  database,
		log: log,
	}
}

// InsertUser insert new user into database.
func (store *StoreUser) InsertUser(ctx context.Context, user User) error {
	prepare, err := store.db.Prepare(InsertUser)
	if err != nil {
		return fmt.Errorf("store.user.InsertUser.Prepare(%v) err: %w: - mycError: %w",
			InsertUser, err, errorssys.ErrPsqlPrepare)
	}

	res, err := prepare.ExecContext(ctx, user.ID, user.FirstName, user.LastName,
		user.Email, user.Password, user.Language, user.Company)
	if err != nil {
		return fmt.Errorf("store.user.InsertUser.Exec(%v): %w", user, store.translateSQLError(err))
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("store.user.InsertUser.RowsAffected: %w", store.translateSQLError(err))
	}

	if affected == 0 {
		return fmt.Errorf("store.user.InsertUser: %w", errorssys.ErrPsqlRowAffected)
	}

	return nil
}

// translateSQLError translate error in SQL object to error go.
//
//nolint:cyclop
func (store *StoreUser) translateSQLError(err error) error {
	var errorPQ *pq.Error

	switch {
	case errors.As(err, &errorPQ):
		switch errorPQ.Code {
		case "22P02":
			return fmt.Errorf("db.translateSqlError(%w) - mycError: {%w}", err, errorssys.ErrUserInvalidInputSyntax)
		case "22021":
			return fmt.Errorf("db.translateSqlError(%w) - mycError: {%w}", err, errorssys.ErrUserInvalidEncoding)
		case "23503":
			return fmt.Errorf("db.translateSqlError(%w) - mycError: {%w}", err, errorssys.ErrUserForeignKeyViolation)
		case "23505":
			switch errorPQ.Constraint {
			case "users_email_key":
				return fmt.Errorf("db.translateSqlError(%w) - mycError: {%w}", err, errorssys.ErrUserDupKeyEmail)
			case "users_pkey":
				return fmt.Errorf("db.translateSqlError(%w) - mycError: {%w}", err, errorssys.ErrUserDupKeyID)
			default:
				return fmt.Errorf("db.translateSqlError(%w) - mycError: {%w}", err, errorssys.ErrUserDupKeyOther)
			}
		case "42703":
			return fmt.Errorf("db.translateSqlError(%w) - mycError: {%w}", err, errorssys.ErrUserUndefinedColumn)
		default:
			return fmt.Errorf("db.translateSqlError(%w) - mycError: {%w}", err, errorssys.ErrUserQuery)
		}
	default:
		return fmt.Errorf("db.translateSqlError(%w) - mycError: {%w}", err, errorssys.ErrUserQuery)
	}
}

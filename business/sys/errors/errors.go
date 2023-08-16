// Package errors to the system.
package errors

import "errors"

var (
	ErrPsql           = errors.New("error database")
	ErrPsqlConnection = errors.New("errors connection database")

	ErrRowAffected               = errors.New("no rows affected")
	ErrGeneratePassHash          = errors.New("generate password hash failed")
	ErrValidatorInvalidCoreModel = errors.New("invalid core model")
)

// Package errors to the system.
package errors

import "errors"

var (
	ErrPsql            = errors.New("error database")
	ErrPsqlConnection  = errors.New("errors connection database")
	ErrPsqlPrepare     = errors.New("prepare syntax error")
	ErrPsqlRowAffected = errors.New("no rows affected")

	ErrEmailFromMailServer   = errors.New("error from mail server")
	ErrConfigEmail           = errors.New("error configuration email")
	ErrEmailRenderTemplate   = errors.New("error while rendering template")
	ErrEmailReadFileTemplate = errors.New("error read file template file")
	ErrEmailSend             = errors.New("error send email")

	ErrGeneratePassHash          = errors.New("generate password hash failed")
	ErrTokenGenerating           = errors.New("error generating token")
	ErrTokenNull                 = errors.New("token generate is null")
	ErrValidatorInvalidCoreModel = errors.New("invalid core model")
)

package errors

import "errors"

var (
	ErrUserDupKeyEmail         = errors.New("a user with that email already exists in the system")
	ErrUserDupKeyID            = errors.New("there is a user that exists for that id")
	ErrUserDupKeyOther         = errors.New("user error duplicated key other")
	ErrUserInvalidInputSyntax  = errors.New("wrong input data from user")
	ErrUserInvalidEncoding     = errors.New("wrong data formatting for user")
	ErrUserUndefinedColumn     = errors.New("column does not exists on users")
	ErrUserForeignKeyViolation = errors.New("the inserted user does not meet a foreign key constraint")
	ErrUserQuery               = errors.New("error performing user query")
)

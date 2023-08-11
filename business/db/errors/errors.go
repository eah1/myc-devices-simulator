// Package errors provides support for access the database.
package errors

import (
	"errors"
	"fmt"
	errorssys "myc-devices-simulator/business/sys/errors"

	"github.com/lib/pq"
	"go.uber.org/zap"
)

// PsqlError errors from sql.
type PsqlError struct {
	CodeSQL        string
	TableName      string
	ConstraintName string
	Err            string
}

// Error redefined error from struct PsqlError.
func (err *PsqlError) Error() string {
	return fmt.Sprintf("code sql %s, table %s, constraints %s : %s",
		err.CodeSQL, err.TableName, err.ConstraintName, err.Err)
}

// WrapperError parse error PgError to custom error.
func WrapperError(log *zap.SugaredLogger, err error) error {
	var pgError *pq.Error

	if errors.As(err, &pgError) {
		return &PsqlError{
			CodeSQL:        pgError.Code.Name(),
			TableName:      pgError.Table,
			ConstraintName: pgError.Constraint,
			Err:            pgError.Message,
		}
	}

	log.Error(err)

	return fmt.Errorf("db.PsqlError: %w", errorssys.ErrPsql)
}

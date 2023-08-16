// Package user to implement use case.
package user

import (
	"context"
	"myc-devices-simulator/business/core/user"
)

//go:generate mockery --name CoreUser

// CoreUser methods core user to store.
type CoreUser interface {
	Create(ctx context.Context, user user.User) (user.User, error)
}

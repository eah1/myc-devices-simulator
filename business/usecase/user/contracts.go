package user

import (
	"context"
	"myc-devices-simulator/business/core/user"
)

// CoreUser methods core user to store.
//
//go:generate mockery --name CoreUser
type CoreUser interface {
	Create(ctx context.Context, user user.User) (user.User, error)
}

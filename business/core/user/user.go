package user

import (
	"context"
	"myc-devices-simulator/business/repository/store/user"
)

//go:generate mockery --name StoreUser

// StoreUser methods store user to call database.
type StoreUser interface {
	InsertUser(ctx context.Context, user user.User) error
}

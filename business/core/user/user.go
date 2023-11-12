// Package user to validation entity and call layer store.
package user

import (
	"context"
	"fmt"
	"myc-devices-simulator/business/core/validator"
	"myc-devices-simulator/business/repository/store/user"
	errorssys "myc-devices-simulator/business/sys/errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

//go:generate mockery --name StoreUser

// StoreUser methods store user to call database.
type StoreUser interface {
	InsertUser(ctx context.Context, user user.User) error
}

// CoreUser core user to call store layer.
type CoreUser struct {
	storeUser StoreUser
}

// NewCoreUser construct a core user group.
func NewCoreUser(storeUser StoreUser) CoreUser {
	return CoreUser{
		storeUser: storeUser,
	}
}

// Create insert a new user into the system.
func (core *CoreUser) Create(ctx context.Context, user User) (User, error) {
	if _, err := validator.IsValid(user); err != nil {
		return User{}, fmt.Errorf("core.user.Create.IsValid: %w", err)
	}

	if err := generatePasswordHash(&user); err != nil {
		return User{}, fmt.Errorf("core.user.Create.generatePasswordHash: %w", err)
	}

	user.ID = uuid.New()

	if err := core.storeUser.InsertUser(ctx, user.ToStore()); err != nil {
		return User{}, fmt.Errorf("core.user.Create.InsertUser: %w", err)
	}

	return user, nil
}

// generatePassword generate password hash.
func generatePasswordHash(user *User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("core.user.generatePasswordHash: %w - errCustom %w", err, errorssys.ErrGeneratePassHash)
	}

	user.Password = string(hash)

	return nil
}

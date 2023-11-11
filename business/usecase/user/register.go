// Package user to implement use case.
package user

import (
	"context"
	"fmt"
)

// UCUserRegister use case user register to call core layers.
type UCUserRegister struct {
	coreUser CoreUser
}

// NewUCUserRegister construct a use case user group.
func NewUCUserRegister(coreUser CoreUser) UCUserRegister {
	return UCUserRegister{
		coreUser: coreUser,
	}
}

// Execute use case layer.
func (uc UCUserRegister) Execute(ctx context.Context, registerUser RegisterUseCase) error {
	if _, err := uc.coreUser.Create(ctx, registerUser.toCoreModel()); err != nil {
		return fmt.Errorf("usecase.user.Execute.Create: %w", err)
	}

	return nil
}

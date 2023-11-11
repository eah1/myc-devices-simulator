// Package user to implement handlers web.
package user

import (
	"context"
	"myc-devices-simulator/business/usecase/user"
)

// UseCaseUserRegister user use case interface.
type UseCaseUserRegister interface {
	Execute(ctx context.Context, registerUser user.RegisterUseCase) error
}

// Package user to implement handlers web.
package user

import (
	"context"
	"myc-devices-simulator/business/usecase/user"

	"github.com/labstack/echo/v4"
)

// UseCaseUserRegister user use case interface.
type UseCaseUserRegister interface {
	Execute(ctx context.Context, registerUser user.RegisterUseCase) error
}

func (h HandlerUser) RegisterUser(ctx echo.Context) error {
	return nil
}

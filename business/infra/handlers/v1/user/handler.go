package user

import "github.com/labstack/echo/v4"

// HandlerUser define the dependencies need end points relations from user.
type HandlerUser struct {
	useCaseUserRegister UseCaseUserRegister
}

// NewHandlerUser initializes a new user handler with its dependencies.
func NewHandlerUser(useCaseUserRegister UseCaseUserRegister) HandlerUser {
	return HandlerUser{
		useCaseUserRegister: useCaseUserRegister,
	}
}

// GroupUser endpoints for user.
func GroupUser(group *echo.Group, prefix string, userHandler HandlerUser) {
	groupUsers := group.Group(prefix)

	groupUsers.POST("/register", userHandler.RegisterUser)
}

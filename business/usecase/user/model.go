package user

import "myc-devices-simulator/business/core/user"

// RegisterUseCase register model from use case layer.
type RegisterUseCase struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
	Language  string
	Company   string
}

// toCoreModel transform data use case to data core model.
func (register RegisterUseCase) toCoreModel() user.User {
	return user.User{
		FirstName: register.FirstName,
		LastName:  register.LastName,
		Email:     register.Email,
		Password:  register.Password,
		Language:  register.Language,
		Company:   register.Company,
	}
}

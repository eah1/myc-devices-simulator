package user

import (
	"myc-devices-simulator/business/core/user"

	"github.com/google/uuid"
)

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

// UserUseCase user model from use case layer.
type UserUseCase struct {
	ID        uuid.UUID
	FirstName string
	LastName  string
	Email     string
	Password  string
	Language  string
	Company   string
}

// toUseCaseModel transform data model to data use case model.
func (user UserUseCase) toUseCaseModel(userCore user.User) UserUseCase {
	return UserUseCase{
		ID:        userCore.ID,
		FirstName: userCore.FirstName,
		LastName:  userCore.LastName,
		Email:     userCore.Email,
		Password:  userCore.Password,
		Language:  userCore.Language,
		Company:   userCore.Company,
	}
}

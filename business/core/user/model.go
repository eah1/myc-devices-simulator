// Package user to validation entity and call layer store.
package user

import (
	store "myc-devices-simulator/business/repository/store/user"

	"github.com/google/uuid"
)

// User represents the structure of the core.
type User struct {
	ID              uuid.UUID
	FirstName       string
	LastName        string
	Email           string
	Password        string
	Language        string `validate:"required,oneof=es en fr pt,max=2"`
	Company         string
	ValidationToken string
	Validated       bool
}

// ToStore transform model core to store.
func (user User) ToStore() store.User {
	return store.User{
		ID:              user.ID,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		Email:           user.Email,
		Password:        user.Password,
		Language:        user.Language,
		Company:         user.Company,
		ValidationToken: user.ValidationToken,
		Validated:       user.Validated,
	}
}

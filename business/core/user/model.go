// Package user to validation entity and call layer store.
package user

import (
	"myc-devices-simulator/business/repository/store/user"

	"github.com/google/uuid"
)

// User represents the structure of the core.
type User struct {
	ID        uuid.UUID
	FirstName string
	LastName  string
	Email     string
	Password  string
	Language  string `validate:"required,oneof=es en fr pt,max=2"`
	Company   string
}

// ToStore transform model core to store.
func (u User) ToStore() user.User {
	return user.User{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Password:  u.Password,
		Language:  u.Language,
		Company:   u.Company,
	}
}

// Package user call database layer.
package user

import (
	"time"

	"github.com/google/uuid"
)

// User represents the structure of the database.
type User struct {
	ID              uuid.UUID `db:"pk not null 'id'"`
	FirstName       string    `db:"first_name"`
	LastName        string    `db:"last_name"`
	Email           string    `db:"not null unique email"`
	Password        string    `db:"password"`
	Language        string    `db:"language"`
	Company         string    `db:"company"`
	ValidationToken string    `db:"validation_token"`
	Validated       bool      `db:"validated"`
	CreatedAt       time.Time `db:"created"`
	UpdatedAt       time.Time `db:"updated"`
}

func (user User) SetID(id uuid.UUID) User {
	user.ID = id

	return user
}

func (user User) SetFirstName(firstName string) User {
	user.FirstName = firstName

	return user
}

func (user User) SetLasName(lastName string) User {
	user.LastName = lastName

	return user
}

func (user User) SetLanguage(language string) User {
	user.Language = language

	return user
}

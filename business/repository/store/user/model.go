// Package user call database layer.
package user

import (
	"time"

	"github.com/google/uuid"
)

// User represents the structure of the database.
type User struct {
	ID        uuid.UUID `db:"pk not null 'id'"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	Email     string    `db:"not null unique email"`
	Password  string    `db:"password"`
	Language  string    `db:"language"`
	Company   string    `db:"company"`
	CreatedAt time.Time `db:"created"`
	UpdatedAt time.Time `db:"updated"`
}

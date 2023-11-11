// Package user to implement use case.
package user

// RegisterUser represents the structure of the use case.
type RegisterUser struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
	Language  string
	Company   string
}

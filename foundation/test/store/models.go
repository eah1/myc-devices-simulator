package store

import (
	"github.com/stretchr/testify/require"
	"myc-devices-simulator/business/repository/store/user"
	"testing"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"syreclabs.com/go/faker"
)

// NewUser new user test layer store.
func NewUser(t *testing.T, testName string) user.User {
	hash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	require.NoError(t, err)

	return user.User{
		ID:        uuid.New(),
		FirstName: faker.Name().FirstName() + "_" + testName,
		LastName:  faker.Name().LastName() + "_" + testName,
		Email:     faker.Internet().Email(),
		Password:  string(hash),
		Language:  faker.RandomChoice([]string{"en", "es", "fr", "pt"}),
		Company:   faker.Company().Name(),
	}
}

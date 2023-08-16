package user_test

import (
	"context"
	"errors"
	"myc-devices-simulator/business/db"
	"myc-devices-simulator/business/db/databasehandler"
	"myc-devices-simulator/business/db/mocks"
	"myc-devices-simulator/business/repository/store/user"
	"myc-devices-simulator/foundation/docker"
	"myc-devices-simulator/foundation/test"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"syreclabs.com/go/faker"
)

func TestUserStore_InsertUser(t *testing.T) {
	t.Parallel()

	testName := "store-user-insert"

	// Create logger.
	newLog := test.InitLogger(t, "t-"+testName)

	// Create DB container
	container := test.InitDockerContainerDatabase(t, testName)

	// Create session database.
	database := test.InitDatabase(t, test.InitConfig(container.Host), newLog)

	// Create a user store.
	storeUser := user.NewUserStore(&databasehandler.SQLDBTx{DB: database}, newLog)

	// Create a password hash to create a user.
	hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	tests := []struct {
		name     string
		database db.SQLGbc
		init     func(store *mocks.SQLGbc)
		user     user.User
		isError  bool
	}{
		{
			name:     "success creating user",
			database: &databasehandler.SQLDBTx{DB: database},
			user: user.User{
				ID:        uuid.NewString(),
				FirstName: faker.Name().FirstName() + "_" + testName,
				LastName:  faker.Name().LastName() + "_" + testName,
				Email:     faker.Internet().Email(),
				Password:  string(hash),
				Language:  faker.RandomChoice([]string{"en", "es", "fr", "pt"}),
				Company:   faker.Company().Name(),
			},
			isError: false,
		},
		{
			name:     "duplicate user",
			database: &databasehandler.SQLDBTx{DB: database},
			user: func(t *testing.T) user.User {
				newUser := user.User{
					ID:        uuid.NewString(),
					FirstName: faker.Name().FirstName() + "_" + testName,
					LastName:  faker.Name().LastName() + "_" + testName,
					Email:     faker.Internet().Email(),
					Password:  string(hash),
					Language:  faker.RandomChoice([]string{"en", "es", "fr", "pt"}),
					Company:   faker.Company().Name(),
				}

				t.Run("duplicate user --> success creating user", func(t *testing.T) {
					assert.Equal(t, nil, storeUser.InsertUser(context.Background(), newUser))
				})

				return newUser
			}(t),
			isError: true,
		},
		{
			name:     "error creating user invalid id",
			database: &databasehandler.SQLDBTx{DB: database},
			user: user.User{
				ID:        faker.RandomString(36),
				FirstName: faker.Name().FirstName() + "_" + testName,
				LastName:  faker.Name().LastName() + "_" + testName,
				Email:     faker.Internet().Email(),
				Password:  string(hash),
				Language:  faker.RandomChoice([]string{"en", "es", "fr", "pt"}),
				Company:   faker.Company().Name(),
			},
			isError: true,
		},
		{
			name:     "error creating user invalid language",
			database: &databasehandler.SQLDBTx{DB: database},
			user: user.User{
				ID:        uuid.NewString(),
				FirstName: faker.Name().FirstName() + "_" + testName,
				LastName:  faker.Name().LastName() + "_" + testName,
				Email:     faker.Internet().Email(),
				Password:  string(hash),
				Language:  faker.RandomString(2),
				Company:   faker.Company().Name(),
			},
			isError: true,
		},
		{
			name:     "mock error in prepare method",
			database: nil,
			init: func(store *mocks.SQLGbc) {
				store.On("Prepare", mock.AnythingOfType("string")).
					Return(nil, errors.New("error Prepare"))
			},
			user: user.User{
				ID:        uuid.NewString(),
				FirstName: faker.Name().FirstName() + "_" + testName,
				LastName:  faker.Name().LastName() + "_" + testName,
				Email:     faker.Internet().Email(),
				Password:  string(hash),
				Language:  faker.RandomChoice([]string{"en", "es", "fr", "pt"}),
				Company:   faker.Company().Name(),
			},
			isError: true,
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var storeUser user.StoreUser

			if test.database == nil {
				database := mocks.NewSQLGbc(t)
				test.init(database)

				storeUser = user.NewUserStore(database, newLog)
			}

			if test.database != nil {
				storeUser = user.NewUserStore(test.database, newLog)
			}

			err := storeUser.InsertUser(context.Background(), test.user)

			if test.isError {
				assert.Error(t, err)
			}

			if !test.isError {
				assert.Equal(t, nil, err)
			}
		})
	}

	t.Cleanup(func() {
		defer docker.StopContainer(container.ID)
		defer os.RemoveAll("new_schema")
	})
}

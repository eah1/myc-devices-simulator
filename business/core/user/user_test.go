package user_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	coremodel "myc-devices-simulator/business/core/user"
	"myc-devices-simulator/business/core/user/mocks"
	"myc-devices-simulator/business/db/databasehandler"
	"myc-devices-simulator/business/repository/store/user"
	"myc-devices-simulator/business/sys/errors"
	"myc-devices-simulator/foundation/docker"
	"myc-devices-simulator/foundation/test"
	"os"
	"syreclabs.com/go/faker"
	"testing"
)

func TestCoreUser_Create(t *testing.T) {
	t.Parallel()

	testName := "user-core-create"

	// Create logger.
	newLog := test.InitLogger(t, "t-"+testName)

	// Create DB container
	container := test.InitDockerContainerDatabase(t, testName)

	// Create session database.
	database := test.InitDatabase(t, test.InitConfig(container.Host), newLog)

	// Create a user store.
	storeUser := user.NewUserStore(&databasehandler.SQLDBTx{DB: database}, newLog)

	// Create a user core.
	coreUser := coremodel.NewCoreUser(&storeUser)

	tests := []struct {
		name      string
		storeUser user.UserStore
		init      func(store *mocks.StoreUser)
		user      coremodel.User
		isError   bool
	}{
		{
			name:      "success creating user",
			storeUser: user.NewUserStore(&databasehandler.SQLDBTx{DB: database}, newLog),
			user: coremodel.User{
				FirstName: faker.Name().FirstName() + "_" + testName,
				LastName:  faker.Name().LastName() + "_" + testName,
				Email:     faker.Internet().Email(),
				Password:  faker.Internet().Password(8, 64),
				Language:  faker.RandomChoice([]string{"en", "es", "fr", "pt"}),
				Company:   faker.Company().Name(),
			},
			isError: false,
		},
		{
			name:      "duplicate create user",
			storeUser: user.NewUserStore(&databasehandler.SQLDBTx{DB: database}, newLog),
			user: func(t *testing.T) coremodel.User {
				newUser := coremodel.User{
					FirstName: faker.Name().FirstName() + "_" + testName,
					LastName:  faker.Name().LastName() + "_" + testName,
					Email:     faker.Internet().Email(),
					Password:  faker.Internet().Password(8, 64),
					Language:  faker.RandomChoice([]string{"en", "es", "fr", "pt"}),
					Company:   faker.Company().Name(),
				}

				t.Run("duplicate create user --> success creating user", func(t *testing.T) {
					_, err := coreUser.Create(context.Background(), newUser)
					assert.Equal(t, nil, err)
				})

				return newUser
			}(t),
			isError: true,
		},
		{
			name:      "error creating user invalid language",
			storeUser: user.NewUserStore(&databasehandler.SQLDBTx{DB: database}, newLog),
			user: coremodel.User{
				FirstName: faker.Name().FirstName() + "_" + testName,
				LastName:  faker.Name().LastName() + "_" + testName,
				Email:     faker.Internet().Email(),
				Password:  faker.Internet().Password(8, 64),
				Language:  faker.RandomString(2),
				Company:   faker.Company().Name(),
			},
			isError: true,
		},
		{
			name:      "error creating user invalid field null",
			storeUser: user.NewUserStore(&databasehandler.SQLDBTx{DB: database}, newLog),
			user: coremodel.User{
				FirstName: "FistName\000",
				LastName:  faker.Name().LastName() + "_" + testName,
				Email:     faker.Internet().Email(),
				Password:  faker.Internet().Password(8, 64),
				Language:  faker.RandomChoice([]string{"en", "es", "fr", "pt"}),
				Company:   faker.Company().Name(),
			},
			isError: true,
		},
		{
			name:      "mock error in row affected",
			storeUser: user.UserStore{},
			init: func(store *mocks.StoreUser) {
				store.On("InsertUser", context.TODO(), mock.AnythingOfType("user.User")).Return(errors.ErrRowAffected)
			},
			user: coremodel.User{
				FirstName: faker.Name().FirstName() + "_" + testName,
				LastName:  faker.Name().LastName() + "_" + testName,
				Email:     faker.Internet().Email(),
				Password:  faker.Internet().Password(8, 64),
				Language:  faker.RandomChoice([]string{"en", "es", "fr", "pt"}),
				Company:   faker.Company().Name(),
			},
			isError: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var coreUser coremodel.CoreUser

			if tt.storeUser == (user.UserStore{}) {
				storeUser := mocks.NewStoreUser(t)
				tt.init(storeUser)

				coreUser = coremodel.NewCoreUser(storeUser)
			}

			if tt.storeUser != (user.UserStore{}) {
				coreUser = coremodel.NewCoreUser(&tt.storeUser)
			}

			newUser, err := coreUser.Create(context.TODO(), tt.user)

			if tt.isError {
				assert.Error(t, err)
			}

			if !tt.isError {
				assert.Equal(t, nil, err)
				assert.NotEmpty(t, newUser)
			}
		})
	}

	t.Cleanup(func() {
		defer docker.StopContainer(container.ID)
		defer os.RemoveAll("new_schema")
	})
}

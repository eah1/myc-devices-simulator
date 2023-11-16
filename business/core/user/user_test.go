package user_test

import (
	"context"
	"errors"
	"fmt"
	"myc-devices-simulator/business/core/user"
	"myc-devices-simulator/business/core/user/mocks"
	"myc-devices-simulator/business/db/databasehandler"
	store "myc-devices-simulator/business/repository/store/user"
	errorssys "myc-devices-simulator/business/sys/errors"
	"myc-devices-simulator/foundation/docker"
	"myc-devices-simulator/foundation/test"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"syreclabs.com/go/faker"
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
	storeUser := store.NewUserStore(&databasehandler.SQLDBTx{DB: database}, newLog)

	// Create a user core.
	coreUser := user.NewUserCore(&storeUser)

	tests := []struct {
		name          string
		user          user.User
		expectedError error
	}{
		{
			name: testName + " success creating user",
			user: user.User{
				FirstName: faker.Name().FirstName() + "_" + testName,
				LastName:  faker.Name().LastName() + "_" + testName,
				Email:     faker.Internet().Email(),
				Password:  faker.Internet().Password(8, 64),
				Language:  faker.RandomChoice([]string{"en", "es", "fr", "pt"}),
				Company:   faker.Company().Name(),
			},
			expectedError: nil,
		},
		{
			name: testName + " duplicate create user",
			user: func(t *testing.T) user.User {
				newUser := user.User{
					FirstName: faker.Name().FirstName() + "_" + testName,
					LastName:  faker.Name().LastName() + "_" + testName,
					Email:     faker.Internet().Email(),
					Password:  faker.Internet().Password(8, 64),
					Language:  faker.RandomChoice([]string{"en", "es", "fr", "pt"}),
					Company:   faker.Company().Name(),
				}

				t.Run(testName+" duplicate create user --> success creating user", func(t *testing.T) {
					userModel, err := coreUser.Create(context.TODO(), newUser)
					assert.Equal(t, nil, err)
					assert.NotEmpty(t, userModel)
				})

				fmt.Println(newUser)

				return newUser
			}(t),
			expectedError: errorssys.ErrUserDupKeyEmail,
		},
		{
			name: testName + " error creating user invalid language",
			user: user.User{
				FirstName: faker.Name().FirstName() + "_" + testName,
				LastName:  faker.Name().LastName() + "_" + testName,
				Email:     faker.Internet().Email(),
				Password:  faker.Internet().Password(8, 64),
				Language:  faker.RandomString(2),
				Company:   faker.Company().Name(),
			},
			expectedError: errorssys.ErrValidatorInvalidCoreModel,
		},
		{
			name: testName + " error creating user invalid field null",
			user: user.User{
				FirstName: "FistName\000",
				LastName:  faker.Name().LastName() + "_" + testName,
				Email:     faker.Internet().Email(),
				Password:  faker.Internet().Password(8, 64),
				Language:  faker.RandomChoice([]string{"en", "es", "fr", "pt"}),
				Company:   faker.Company().Name(),
			},
			expectedError: errorssys.ErrUserInvalidEncoding,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userModel, err := coreUser.Create(context.TODO(), tt.user)
			assert.Equal(t, true, errors.Is(err, tt.expectedError))

			if tt.expectedError == nil {
				assert.NotEmpty(t, userModel)
			}

			if tt.expectedError != nil {
				assert.Empty(t, userModel)
			}
		})
	}

	testsMock := []struct {
		name          string
		init          func(storeUser *mocks.StoreUser)
		user          user.User
		expectedError error
	}{
		{
			name: testName + " mock error in row affected",
			init: func(storeUser *mocks.StoreUser) {
				storeUser.On("InsertUser", context.TODO(), mock.AnythingOfType("user.User")).Return(errorssys.ErrPsqlRowAffected)
			},
			user: user.User{
				FirstName: faker.Name().FirstName() + "_" + testName,
				LastName:  faker.Name().LastName() + "_" + testName,
				Email:     faker.Internet().Email(),
				Password:  faker.Internet().Password(8, 64),
				Language:  faker.RandomChoice([]string{"en", "es", "fr", "pt"}),
				Company:   faker.Company().Name(),
			},
			expectedError: errorssys.ErrPsqlRowAffected,
		},
		{
			name: testName + " mock error in unified column",
			init: func(storeUser *mocks.StoreUser) {
				storeUser.On("InsertUser", context.TODO(), mock.AnythingOfType("user.User")).Return(errorssys.ErrUserUndefinedColumn)
			},
			user: user.User{
				FirstName: faker.Name().FirstName() + "_" + testName,
				LastName:  faker.Name().LastName() + "_" + testName,
				Email:     faker.Internet().Email(),
				Password:  faker.Internet().Password(8, 64),
				Language:  faker.RandomChoice([]string{"en", "es", "fr", "pt"}),
				Company:   faker.Company().Name(),
			},
			expectedError: errorssys.ErrUserUndefinedColumn,
		},
		{
			name: testName + " mock error in row affected in update user",
			init: func(storeUser *mocks.StoreUser) {
				storeUser.On("InsertUser", context.TODO(), mock.AnythingOfType("user.User")).Return(nil)
				storeUser.On("UpdateUser", context.TODO(), mock.AnythingOfType("user.User")).Return(errorssys.ErrPsqlRowAffected)
			},
			user: user.User{
				FirstName: faker.Name().FirstName() + "_" + testName,
				LastName:  faker.Name().LastName() + "_" + testName,
				Email:     faker.Internet().Email(),
				Password:  faker.Internet().Password(8, 64),
				Language:  faker.RandomChoice([]string{"en", "es", "fr", "pt"}),
				Company:   faker.Company().Name(),
			},
			expectedError: errorssys.ErrPsqlRowAffected,
		},
		{
			name: testName + " mock error in unified column in update user",
			init: func(storeUser *mocks.StoreUser) {
				storeUser.On("InsertUser", context.TODO(), mock.AnythingOfType("user.User")).Return(nil)
				storeUser.On("UpdateUser", context.TODO(), mock.AnythingOfType("user.User")).Return(errorssys.ErrUserUndefinedColumn)
			},
			user: user.User{
				FirstName: faker.Name().FirstName() + "_" + testName,
				LastName:  faker.Name().LastName() + "_" + testName,
				Email:     faker.Internet().Email(),
				Password:  faker.Internet().Password(8, 64),
				Language:  faker.RandomChoice([]string{"en", "es", "fr", "pt"}),
				Company:   faker.Company().Name(),
			},
			expectedError: errorssys.ErrUserUndefinedColumn,
		},
	}

	for _, tt := range testsMock {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			storeUser := mocks.NewStoreUser(t)
			tt.init(storeUser)

			coreUser := user.NewUserCore(storeUser)

			userModel, err := coreUser.Create(context.TODO(), tt.user)
			assert.Equal(t, true, errors.Is(err, tt.expectedError))

			if tt.expectedError == nil {
				assert.NotEmpty(t, userModel)
			}
		})
	}

	t.Cleanup(func() {
		defer func(id string) {
			require.NoError(t, docker.StopContainer(id))
		}(container.ID)

		defer func() {
			require.NoError(t, os.RemoveAll("new_schema"))
		}()
	})
}

package user_test

import (
	"context"
	core "myc-devices-simulator/business/core/user"
	"myc-devices-simulator/business/db/databasehandler"
	store "myc-devices-simulator/business/repository/store/user"
	"myc-devices-simulator/business/usecase/user"
	"myc-devices-simulator/business/usecase/user/mocks"
	"myc-devices-simulator/foundation/docker"
	"myc-devices-simulator/foundation/test"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"syreclabs.com/go/faker"
)

func TestUCUser_Execute(t *testing.T) {
	t.Parallel()

	testName := "user-uc-register"

	// Create logger.
	newLog := test.InitLogger(t, "t-"+testName)

	// Create DB container
	container := test.InitDockerContainerDatabase(t, testName)

	// Create session database.
	database := test.InitDatabase(t, test.InitConfig(container.Host), newLog)

	// Create a user store.
	storeUser := store.NewUserStore(&databasehandler.SQLDBTx{DB: database}, newLog)

	// Create a user core.
	coreUser := core.NewCoreUser(&storeUser)

	// Create a use case register user.
	ucUserRegister := user.NewUCUserRegister(&coreUser)

	tests := []struct {
		name          string
		init          func(coreUser *mocks.CoreUser)
		userRegister  user.RegisterUseCase
		expectedError bool
	}{
		{
			name: testName + " success user register",
			init: func(coreUser *mocks.CoreUser) {},
			userRegister: user.RegisterUseCase{
				FirstName: faker.Name().FirstName() + "_" + testName,
				LastName:  faker.Name().LastName() + "_" + testName,
				Email:     faker.Internet().Email(),
				Password:  faker.Internet().Password(8, 64),
				Language:  faker.RandomChoice([]string{"en", "es", "fr", "pt"}),
				Company:   faker.Company().Name(),
			},
			expectedError: false,
		},
		{
			name: testName + " duplicate user register",
			init: func(coreUser *mocks.CoreUser) {},
			userRegister: func(t *testing.T) user.RegisterUseCase {
				newUserRegister := user.RegisterUseCase{
					FirstName: faker.Name().FirstName() + "_" + testName,
					LastName:  faker.Name().LastName() + "_" + testName,
					Email:     faker.Internet().Email(),
					Password:  faker.Internet().Password(8, 64),
					Language:  faker.RandomChoice([]string{"en", "es", "fr", "pt"}),
					Company:   faker.Company().Name(),
				}

				t.Run(" duplicate create user register --> success user register", func(t *testing.T) {
					assert.Equal(t, nil, ucUserRegister.Execute(context.TODO(), newUserRegister))
				})

				return newUserRegister
			}(t),
			expectedError: true,
		},
		{
			name: testName + " wrong language is invalid",
			init: func(coreUser *mocks.CoreUser) {},
			userRegister: user.RegisterUseCase{
				FirstName: faker.Name().FirstName() + "_" + testName,
				LastName:  faker.Name().LastName() + "_" + testName,
				Email:     faker.Internet().Email(),
				Password:  faker.Internet().Password(8, 64),
				Language:  faker.RandomString(2),
				Company:   faker.Company().Name(),
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.init(nil)

			err := ucUserRegister.Execute(context.TODO(), tt.userRegister)

			if tt.expectedError {
				assert.Error(t, err)
			}

			if !tt.expectedError {
				assert.Equal(t, nil, err)
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

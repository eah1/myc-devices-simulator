package user_test

import (
	"context"
	"errors"
	"github.com/stretchr/testify/mock"
	"myc-devices-simulator/business/db/databasehandler"
	"myc-devices-simulator/business/db/mocks"
	"myc-devices-simulator/business/repository/store/user"
	errorssys "myc-devices-simulator/business/sys/errors"
	"myc-devices-simulator/foundation/docker"
	"myc-devices-simulator/foundation/test"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	hash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	require.NoError(t, err)

	tests := []struct {
		name          string
		user          user.User
		expectedError error
	}{
		{
			name: testName + " success creating user",
			user: user.User{
				ID:        uuid.New(),
				FirstName: faker.Name().FirstName() + "_" + testName,
				LastName:  faker.Name().LastName() + "_" + testName,
				Email:     faker.Internet().Email(),
				Password:  string(hash),
				Language:  faker.RandomChoice([]string{"en", "es", "fr", "pt"}),
				Company:   faker.Company().Name(),
			},
			expectedError: nil,
		},
		{
			name: testName + " duplicate user, from id",
			user: func(t *testing.T) user.User {
				newUser := user.User{
					ID:        uuid.New(),
					FirstName: faker.Name().FirstName() + "_" + testName,
					LastName:  faker.Name().LastName() + "_" + testName,
					Email:     faker.Internet().Email(),
					Password:  string(hash),
					Language:  faker.RandomChoice([]string{"en", "es", "fr", "pt"}),
					Company:   faker.Company().Name(),
				}

				t.Run(testName+" duplicate user --> success creating user", func(t *testing.T) {
					assert.Equal(t, nil, storeUser.InsertUser(context.TODO(), newUser))
				})

				return newUser
			}(t),
			expectedError: errorssys.ErrUserDupKeyID,
		},
		{
			name: testName + " duplicate user, from email",
			user: func(t *testing.T) user.User {
				newUser := user.User{
					ID:        uuid.New(),
					FirstName: faker.Name().FirstName() + "_" + testName,
					LastName:  faker.Name().LastName() + "_" + testName,
					Email:     faker.Internet().Email(),
					Password:  string(hash),
					Language:  faker.RandomChoice([]string{"en", "es", "fr", "pt"}),
					Company:   faker.Company().Name(),
				}

				t.Run(testName+" duplicate user --> success creating user", func(t *testing.T) {
					assert.Equal(t, nil, storeUser.InsertUser(context.Background(), newUser))
				})

				newUser.ID = uuid.New()

				return newUser
			}(t),
			expectedError: errorssys.ErrUserDupKeyEmail,
		},
		{
			name: testName + " error creating user invalid language",
			user: user.User{
				ID:        uuid.New(),
				FirstName: faker.Name().FirstName() + "_" + testName,
				LastName:  faker.Name().LastName() + "_" + testName,
				Email:     faker.Internet().Email(),
				Password:  string(hash),
				Language:  faker.RandomString(2),
				Company:   faker.Company().Name(),
			},
			expectedError: errorssys.ErrUserInvalidInputSyntax,
		},
		{
			name: testName + " null (\000) in field",
			user: user.User{
				ID:        uuid.New(),
				FirstName: "\000",
				LastName:  faker.Name().LastName() + "_" + testName,
				Email:     faker.Internet().Email(),
				Password:  string(hash),
				Language:  faker.RandomString(2),
				Company:   faker.Company().Name(),
			},
			expectedError: errorssys.ErrUserInvalidEncoding,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := storeUser.InsertUser(context.TODO(), tt.user)
			assert.Equal(t, true, errors.Is(err, tt.expectedError))
		})
	}

	testsMock := []struct {
		name          string
		init          func(db *mocks.SQLGbc)
		user          user.User
		expectedError error
	}{
		{
			name: testName + " error in prepare method",
			init: func(db *mocks.SQLGbc) {
				db.On("Prepare", mock.AnythingOfType("string")).Return(nil, errors.New("error Prepare"))
			},
			user: user.User{
				ID:        uuid.New(),
				FirstName: faker.Name().FirstName() + "_" + testName,
				LastName:  faker.Name().LastName() + "_" + testName,
				Email:     faker.Internet().Email(),
				Password:  string(hash),
				Language:  faker.RandomChoice([]string{"en", "es", "fr", "pt"}),
				Company:   faker.Company().Name(),
			},
			expectedError: errorssys.ErrPsqlPrepare,
		},
	}

	for _, tt := range testsMock {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db := mocks.NewSQLGbc(t)
			tt.init(db)

			storeUser := user.NewUserStore(db, newLog)

			err := storeUser.InsertUser(context.TODO(), tt.user)
			assert.Equal(t, true, errors.Is(err, tt.expectedError))
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

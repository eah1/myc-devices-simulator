package user_test

import (
	"context"
	"errors"
	"myc-devices-simulator/business/core/email"
	core "myc-devices-simulator/business/core/user"
	"myc-devices-simulator/business/db/databasehandler"
	store "myc-devices-simulator/business/repository/store/user"
	errorssys "myc-devices-simulator/business/sys/errors"
	"myc-devices-simulator/business/usecase/user"
	"myc-devices-simulator/business/usecase/user/mocks"
	"myc-devices-simulator/foundation/docker"
	"myc-devices-simulator/foundation/test"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

	// Create config.
	config := test.InitConfig(container.Host)

	// Create session database.
	database := test.InitDatabase(t, config, newLog)

	// Create email config.
	emailConfig := test.InitEmailConfig(t, config)

	// Create a user store.
	storeUser := store.NewUserStore(&databasehandler.SQLDBTx{DB: database}, newLog)

	// Create a user core.
	coreUser := core.NewUserCore(&storeUser)

	coreEmail := email.NewEmailCore(emailConfig, config.SMTPFrom, newLog, config)

	// Create a use case register user.
	ucUserRegister := user.NewUCUserRegister(&coreUser, &coreEmail)

	tests := []struct {
		name          string
		userRegister  user.RegisterUseCase
		expectedError error
	}{
		{
			name: testName + " success user register",
			userRegister: user.RegisterUseCase{
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
			name: testName + " duplicate user register",
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
			expectedError: errorssys.ErrUserDupKeyEmail,
		},
		{
			name: testName + " wrong language is invalid",
			userRegister: user.RegisterUseCase{
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
			name: testName + " wrong field is null",
			userRegister: user.RegisterUseCase{
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

			err := ucUserRegister.Execute(context.TODO(), tt.userRegister)
			assert.Equal(t, true, errors.Is(err, tt.expectedError))
		})
	}

	testsMock := []struct {
		name          string
		init          func(coreUser *mocks.CoreUser, emailManager *mocks.EmailManager)
		userRegister  user.RegisterUseCase
		expectedError error
	}{
		{
			name: testName + " mock error in generate hash password",
			init: func(coreUser *mocks.CoreUser, emailManager *mocks.EmailManager) {
				coreUser.On("Create", context.TODO(), mock.AnythingOfType("user.User")).Return(core.User{}, errorssys.ErrGeneratePassHash)
			},
			userRegister: user.RegisterUseCase{
				FirstName: faker.Name().FirstName() + "_" + testName,
				LastName:  faker.Name().LastName() + "_" + testName,
				Email:     faker.Internet().Email(),
				Password:  faker.Internet().Password(8, 64),
				Language:  faker.RandomChoice([]string{"en", "es", "fr", "pt"}),
				Company:   faker.Company().Name(),
			},
			expectedError: errorssys.ErrGeneratePassHash,
		},
		{
			name: testName + " mock error in generate token activation",
			init: func(coreUser *mocks.CoreUser, emailManager *mocks.EmailManager) {
				coreUser.On("Create", context.TODO(), mock.AnythingOfType("user.User")).Return(core.User{}, errorssys.ErrTokenGenerating)
			},
			userRegister: user.RegisterUseCase{
				FirstName: faker.Name().FirstName() + "_" + testName,
				LastName:  faker.Name().LastName() + "_" + testName,
				Email:     faker.Internet().Email(),
				Password:  faker.Internet().Password(8, 64),
				Language:  faker.RandomChoice([]string{"en", "es", "fr", "pt"}),
				Company:   faker.Company().Name(),
			},
			expectedError: errorssys.ErrTokenGenerating,
		},
		{
			name: testName + " mock error open template",
			init: func(coreUser *mocks.CoreUser, emailManager *mocks.EmailManager) {
				coreUser.On("Create", context.TODO(), mock.AnythingOfType("user.User")).Return(core.User{
					Email:    faker.Internet().Email(),
					Language: faker.RandomString(2)}, nil)
			},
			userRegister: user.RegisterUseCase{
				FirstName: faker.Name().FirstName() + "_" + testName,
				LastName:  faker.Name().LastName() + "_" + testName,
				Email:     faker.Internet().Email(),
				Password:  faker.Internet().Password(8, 64),
				Language:  faker.RandomChoice([]string{"en", "es", "fr", "pt"}),
				Company:   faker.Company().Name(),
			},
			expectedError: errorssys.ErrEmailReadFileTemplate,
		},
		{
			name: testName + " mock error in send email",
			init: func(coreUser *mocks.CoreUser, emailManager *mocks.EmailManager) {
				coreUser.On("Create", context.TODO(), mock.AnythingOfType("user.User")).Return(core.User{
					Email:    faker.Internet().Email(),
					Language: faker.RandomChoice([]string{"en", "es", "fr", "pt"})}, nil)
				emailManager.On("SendEmailBody", mock.AnythingOfType("bytes.Buffer"), mock.AnythingOfType("string"),
					mock.AnythingOfType("string"), mock.AnythingOfType("[]string")).Return(errorssys.ErrEmailSend)
			},
			userRegister: user.RegisterUseCase{
				FirstName: faker.Name().FirstName() + "_" + testName,
				LastName:  faker.Name().LastName() + "_" + testName,
				Email:     faker.Internet().Email(),
				Password:  faker.Internet().Password(8, 64),
				Language:  faker.RandomChoice([]string{"en", "es", "fr", "pt"}),
				Company:   faker.Company().Name(),
			},
			expectedError: errorssys.ErrEmailSend,
		},
	}

	for _, tt := range testsMock {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			coreUser := mocks.NewCoreUser(t)
			emailManager := mocks.NewEmailManager(t)
			tt.init(coreUser, emailManager)

			ucUserRegister := user.NewUCUserRegister(coreUser, emailManager)

			err := ucUserRegister.Execute(context.TODO(), tt.userRegister)
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

// Package user to implement use case.
package user

import (
	"context"
	"fmt"
	"myc-devices-simulator/business/core/email/template"
)

// UCUserRegister use case user register to call core layers.
type UCUserRegister struct {
	coreUser     CoreUser
	emailManager EmailManager
}

// NewUCUserRegister construct a use case user group.
func NewUCUserRegister(coreUser CoreUser, emailManager EmailManager) UCUserRegister {
	return UCUserRegister{
		coreUser:     coreUser,
		emailManager: emailManager,
	}
}

// Execute use case layer.
func (uc UCUserRegister) Execute(ctx context.Context, registerUser RegisterUseCase) error {
	userCore, err := uc.coreUser.Create(ctx, registerUser.toCoreModel())
	if err != nil {
		return fmt.Errorf("usecase.user.Execute.Create: %w", err)
	}

	user := UserUseCase{}.toUseCaseModel(userCore)

	body, err := template.Render(user.Language, "account-validation.html", struct {
		Email, ValidationURI string
	}{Email: user.Email, ValidationURI: ""})
	if err != nil {
		return fmt.Errorf("usecase.user.Execute.Render(-) - error: {%w}", err)
	}

	if err := uc.emailManager.SendEmailBody(*body, "Welcome", "account-validation", []string{user.Email}); err != nil {
		return fmt.Errorf("usecase.user.Execute.SendEmailBody - error: {%w}", err)
	}

	return nil
}

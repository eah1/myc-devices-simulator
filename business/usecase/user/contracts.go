package user

import (
	"bytes"
	"context"
	"myc-devices-simulator/business/core/user"
)

// CoreUser methods core user to store.
//
//go:generate mockery --name CoreUser
type CoreUser interface {
	Create(ctx context.Context, user user.User) (user.User, error)
}

// EmailManager interface to send emails.
//
//go:generate mockery --name EmailManager
type EmailManager interface {
	SendEmailBody(body bytes.Buffer, subject, tag string, recipient []string) error
}

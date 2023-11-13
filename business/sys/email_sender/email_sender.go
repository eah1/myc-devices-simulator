// Package emailsender contains system email sender library.
package emailsender

import (
	"crypto/tls"
	"fmt"
	errorssys "myc-devices-simulator/business/sys/errors"
	"myc-devices-simulator/cmd/config"
	"net"
	"net/smtp"
	"strconv"

	"github.com/jhillyerd/enmime"
)

const portDefault = 587

// loginAuth struct login email config.
type loginAuth struct {
	username, password string
}

// authorization auth email config.
func authorization(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

// Start return login command from auth server SMTP.
func (a *loginAuth) Start(*smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(a.username), nil
}

// Next handles SMTP server responses in the authentication process.
func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, fmt.Errorf("email_sender.Next.Default(-) mycError: {%w}", errorssys.ErrEmailFromMailServer)
		}
	}

	return nil, nil
}

// InnitEmailConfig create a basic configuration email service.
func InnitEmailConfig(config config.Config) (*enmime.SMTPSender, error) {
	// connection tcp.
	conn, err := net.Dial("tcp", config.SMTPHost+":"+config.SMTPPort)
	if err != nil {
		return nil, fmt.Errorf(
			"mail_sender.InnitEmailConfig.Dial(-) - error: {%w} - mycError: {%w}", err, errorssys.ErrConfigEmail)
	}

	client, err := smtp.NewClient(conn, config.SMTPHost)
	if err != nil {
		return nil, fmt.Errorf(
			"mail_sender.InnitEmailConfig.NewClient(-) - error: {%w} - mycError: {%w}", err, errorssys.ErrConfigEmail)
	}

	port, err := strconv.Atoi(config.SMTPPort)
	if err != nil {
		return nil, fmt.Errorf(
			"mail_sender.InnitEmailConfig.Atoi(-) - error: {%w} - mycError: {%w}", err, errorssys.ErrConfigEmail)
	}

	if port == portDefault {
		// TLS config.
		tlsConfig := new(tls.Config)
		tlsConfig.ServerName = config.SMTPHost

		if err = client.StartTLS(tlsConfig); err != nil {
			return nil, fmt.Errorf("mail_sender.InnitEmailConfig.StartTLS(-)"+
				" - error: {%v} - mycError: {%w}", err, errorssys.ErrConfigEmail)
		}
	}

	// authentication configuration.
	auth := authorization(config.PostmarkToken, config.PostmarkToken)

	if err = client.Auth(auth); err != nil {
		return nil, fmt.Errorf("mail_sender.InnitEmailConfig.Auth(-) "+
			"- error: {%v} -  mycError: {%w}", err, errorssys.ErrConfigEmail)
	}

	sender := enmime.NewSMTP(config.SMTPHost+":"+config.SMTPPort, auth)

	return sender, nil
}

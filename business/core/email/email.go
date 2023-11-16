// Package email send email.
package email

import (
	"bytes"
	"fmt"
	errorssys "myc-devices-simulator/business/sys/errors"
	"myc-devices-simulator/cmd/config"
	"regexp"

	"github.com/jhillyerd/enmime"
	"go.uber.org/zap"
)

var regex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// EmailCore struct to send emails.
type EmailCore struct {
	sender   *enmime.SMTPSender
	smtpFrom string
	log      *zap.SugaredLogger
	conf     config.Config
}

// NewEmailCore create a new email service struct.
func NewEmailCore(emailSender *enmime.SMTPSender, smtpFrom string, log *zap.SugaredLogger, conf config.Config) EmailCore {
	return EmailCore{sender: emailSender, smtpFrom: smtpFrom, log: log, conf: conf}
}

// SendEmailBody email not contend subjects fields.
func (core EmailCore) SendEmailBody(body bytes.Buffer, subject, tag string, recipient []string) error {
	if body.Len() == 0 {
		return fmt.Errorf("core.email.SendEmailBody(-) - error: invalid body - mycError: {%w}", errorssys.ErrEmailSend)
	}

	master := enmime.Builder().
		From("CIRCUTOR", core.smtpFrom).
		Subject(subject).
		HTML(body.Bytes()).Header("X-PM-Tag", tag).Header("X-PM-Metadata-env", core.conf.Environment)

	for _, emailAddress := range recipient {
		if !isValidEmail(emailAddress) {
			core.log.Warnw(fmt.Sprintf(
				"core.email.SendEmailBody.isValidEmail(-) - warning: invalid email - mycError: {%s}", errorssys.ErrEmailSend))

			continue
		}

		msg := master.To(emailAddress, emailAddress)

		if err := msg.Send(core.sender); err != nil {
			return fmt.Errorf("core.email.SendEmailBody.Send() - error: {%w} mycError: {%w}", err, errorssys.ErrEmailSend)
		}
	}

	return nil
}

func isValidEmail(email string) bool {
	return regex.MatchString(email)
}

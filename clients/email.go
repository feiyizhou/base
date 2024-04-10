package clients

import (
	"context"
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

type EmailConf struct {
	FromEmail      string `json:"fromEmail" mapstructure:"fromEmail"`
	SMTPServer     string `json:"smtpServer" mapstructure:"smtpServer"`
	SMTPPort       int    `json:"smtpPort" mapstructure:"smtpPort"`
	SMTPVerifyCode string `json:"smtpVerifyCode" mapstructure:"smtpVerifyCode"`
}

type EmailClient struct {
	ctx  context.Context
	conf EmailConf
}

func NewEmailClient(ctx context.Context, conf EmailConf) *EmailClient {
	return &EmailClient{ctx: ctx, conf: conf}
}

func (ec *EmailClient) Send(desEmails []string, subject string, content []byte) error {
	em := email.NewEmail()
	em.From = ec.conf.FromEmail
	em.To = desEmails
	em.Subject = subject
	em.Text = content
	return em.Send(fmt.Sprintf("%s:%d", ec.conf.SMTPServer, ec.conf.SMTPPort), smtp.PlainAuth(
		"",
		ec.conf.FromEmail,
		ec.conf.SMTPVerifyCode,
		ec.conf.SMTPServer,
	))
}

package email

import (
	"fmt"
	"net/smtp"

	"github.com/lithammer/dedent"
)

type SMTPCreds struct {
	User string
	Pass string
	Host string
	Port string
}

type Client interface {
	SendEmail(to, subject, body string) error
}

type client struct {
	creds SMTPCreds
	auth  smtp.Auth
}

func NewClient(creds SMTPCreds) Client {
	auth := smtp.PlainAuth("", creds.User, creds.Pass, creds.Host)
	return &client{
		creds: creds,
		auth:  auth,
	}
}

func (c *client) SendEmail(to, subject, body string) error {
	message := []byte(dedent.Dedent(
		fmt.Sprintf(`
			From: %s\r\n
			To: %s\r\n
			Subject: %s\r\n\r\n
			%s
		`, c.creds.User, to, subject, body),
	))
	return smtp.SendMail(
		c.creds.Host+":"+c.creds.Port,
		c.auth,
		c.creds.User,
		[]string{to},
		message,
	)
}

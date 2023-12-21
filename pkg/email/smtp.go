package email

import (
	"fmt"
	"net/smtp"
	"time"
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
	message := []byte(
		fmt.Sprintf("From: %s\n", c.creds.User) +
			fmt.Sprintf("To: %s\n", to) +
			fmt.Sprintf("Subject: %s\n", subject) +
			fmt.Sprintf("Date: %s\n", time.Now().Format(time.RFC1123Z)) +
			"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" +
			"<html><body>" +
			body +
			"</body></html>",
	)
	return smtp.SendMail(
		c.creds.Host+":"+c.creds.Port,
		c.auth,
		c.creds.User,
		[]string{to},
		message,
	)
}

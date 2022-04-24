package email

type Address string

type SMTPCreds struct {
	User string
	Pass string
	Host string
	Port string
}

type Client interface {
	SendEmail(to Address)
}

type client struct {
	creds SMTPCreds
}

func NewClient(creds SMTPCreds) Client {
	return &client{
		creds: creds,
	}
}

func (*client) SendEmail(_ Address) {}

package email

type EmailAddress string

type Client struct {
	SendgridKey string
}

func NewClient(sendgridKey string) *Client {
	return &Client{
		SendgridKey: sendgridKey,
	}
}

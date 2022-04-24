package services

import (
	"database/sql"

	"github.com/azuline/presage/pkg/email"
	_ "modernc.org/sqlite"
)

type Services struct {
	DB    *sql.DB
	Email email.Client
}

func Initialize(
	databasePath string,
	smtpCreds email.SMTPCreds,
) (*Services, error) {
	// Initialize services.
	db, err := sql.Open("sqlite", databasePath)
	if err != nil {
		return nil, err
	}

	emailClient := email.NewClient(smtpCreds)

	return &Services{
		DB:    db,
		Email: emailClient,
	}, nil
}

package services

import (
	"database/sql"

	"github.com/azuline/presage/pkg/email"
)

type Services struct {
	DB    *sql.DB
	Email email.Client
}

func Initialize(
	sendgridKey string,
	databaseURI string,
) (*Services, error) {
	// Initialize services.
	db, err := sql.Open("sqlite", databaseURI)
	if err != nil {
		return nil, err
	}

	emailClient := email.NewClient(sendgridKey)

	return &Services{
		DB:    db,
		Email: emailClient,
	}, nil
}

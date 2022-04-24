package services

import (
	"database/sql"

	// Import database driver.
	_ "modernc.org/sqlite"

	"github.com/azuline/presage/pkg/email"
)

type Services struct {
	DB    *sql.DB
	Email email.Client
}

// Initialize initializes services used across the application via a dependency
// injection.
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

package services

import (
	"database/sql"

	// Import database driver.
	_ "modernc.org/sqlite"

	"github.com/azuline/presage/pkg/email"
	"github.com/jmoiron/sqlx"
)

type Services struct {
	DB      *sqlx.DB
	PlainDB *sql.DB
	Email   email.Client
}

// Initialize initializes services used across the application via a dependency
// injection.
func Initialize(
	databasePath string,
	smtpCreds email.SMTPCreds,
) (*Services, error) {
	// Initialize services.
	plainDB, err := sql.Open("sqlite", databasePath)
	if err != nil {
		return nil, err
	}

	db := sqlx.NewDb(plainDB, "sqlite")
	emailClient := email.NewClient(smtpCreds)

	return &Services{
		DB:      db,
		PlainDB: plainDB,
		Email:   emailClient,
	}, nil
}

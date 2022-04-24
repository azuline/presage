package fixtures

import (
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/azuline/presage/pkg/email"
	"github.com/azuline/presage/pkg/migrate"
	"github.com/azuline/presage/pkg/services"
)

func Services(t *testing.T) *services.Services {
	// Database.
	testdir := t.TempDir()
	databasePath := filepath.Join(testdir, "presage.sqlite")

	plainDB, err := sql.Open("sqlite", databasePath)
	require.NoError(t, err)

	db := sqlx.NewDb(plainDB, "sqlite")

	// Email.
	smtpCreds := email.SMTPCreds{
		User: "username@ema.il",
		Pass: "hopeyouknowit",
		Host: "lol.find.your.own",
		Port: "465",
	}
	emailClient := email.NewClient(smtpCreds)

	// Migrate database.
	srv := &services.Services{
		PlainDB: plainDB,
		DB:      db,
		Email:   emailClient,
	}
	err = migrate.Migrate(srv)
	require.NoError(t, err)

	return srv
}

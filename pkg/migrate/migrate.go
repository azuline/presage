package migrate

import (
	"embed"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/pkg/errors"

	"github.com/azuline/presage/pkg/services"
)

//go:embed sql/*.sql
var fs embed.FS

// Migrate migrates the SQLite database to the most recent version, applying
// the migrations in the `sql/` directory.
func Migrate(srv *services.Services) error {
	migs, err := iofs.New(fs, "sql")
	if err != nil {
		return errors.Wrap(err, "failed to create migration iofs")
	}

	instance, err := sqlite.WithInstance(srv.PlainDB, &sqlite.Config{})
	if err != nil {
		return errors.Wrap(err, "failed to create sqlite instance")
	}

	m, err := migrate.NewWithInstance("iofs", migs, "sqlite", instance)
	if err != nil {
		return errors.Wrap(err, "failed to create migrate instance")
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return errors.Wrap(err, "failed to migrate database")
	}

	return nil
}

package migrations
// Package migrations applies the embedded SQL schema at startup.
//
// The .sql files are compiled into the binary, so the container image needs no
// extra files and the schema can never drift from the code that uses it.
package migrations

import (
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // registers the "postgres" scheme
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed *.sql
var files embed.FS

// Run applies every pending migration. It is a no-op when the schema is current.
func Run(dsn string) error {
	source, err := iofs.New(files, ".")
	if err != nil {
		return fmt.Errorf("migrations: read embedded files: %w", err)
	}

	migrator, err := migrate.NewWithSourceInstance("iofs", source, dsn)
	if err != nil {
		return fmt.Errorf("migrations: open: %w", err)
	}
	defer func() {
		_, _ = migrator.Close()
	}()

	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrations: apply: %w", err)
	}
	return nil
}

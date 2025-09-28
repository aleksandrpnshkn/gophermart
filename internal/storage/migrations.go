package storage

import (
	"embed"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrationsFiles embed.FS

func RunMigrations(databaseDSN string) error {
	sourceDriver, err := iofs.New(migrationsFiles, "migrations")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithSourceInstance("iofs", sourceDriver, databaseDSN)
	if err != nil {
		return err
	}

	return m.Up()
}

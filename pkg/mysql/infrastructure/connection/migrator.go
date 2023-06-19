package connection

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	migrate "github.com/rubenv/sql-migrate"
)

type Migrator interface {
	MigrateUp() error
	MigrateDown() error
}

func NewMigrator(db *sqlx.DB, fileSystem http.FileSystem) Migrator {
	return &migrator{
		db:         db,
		fileSystem: fileSystem,
	}
}

type migrator struct {
	db         *sqlx.DB
	fileSystem http.FileSystem
}

func (m *migrator) MigrateUp() error {
	_, err := migrate.Exec(m.db.DB, driverName, migrate.HttpFileSystemMigrationSource{FileSystem: m.fileSystem}, migrate.Up)
	return errors.Wrap(err, "failed to migrate")
}

func (m *migrator) MigrateDown() error {
	_, err := migrate.Exec(m.db.DB, driverName, migrate.HttpFileSystemMigrationSource{FileSystem: m.fileSystem}, migrate.Down)
	return errors.Wrap(err, "failed to migrate")
}

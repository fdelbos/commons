package pg

import (
	"errors"
	"io/fs"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

// Migrate migrates the database to the latest version.
// You can create a migration file with the following command:
//
//	migrate create -ext sql -dir <sql_fs_dir> -seq <migration_name>
func Migrate(destURL string, sqlFS fs.FS) error {
	source, err := iofs.New(sqlFS, ".")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithSourceInstance("iofs", source, destURL)
	if err != nil {
		return err
	}
	defer m.Close()

	if err = m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}
		return err
	}

	return nil
}

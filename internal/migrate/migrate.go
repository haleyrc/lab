package migrate

import (
	"context"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/jmoiron/sqlx"
)

type Migration struct {
	Name string
	SQL  string
}

func Down(ctx context.Context, db *sqlx.DB) error {
	dir := filepath.Join(".", "migrations")
	return DownFrom(ctx, db, dir)
}

func DownFrom(ctx context.Context, db *sqlx.DB, dir string) error {
	downMigrations, err := loadMigrationsFrom(ctx, dir, "*.down.sql")
	if err != nil {
		return err
	}

	if len(downMigrations) == 0 {
		log.Println("migrate: no down migrations found")
		return nil
	}

	return runInTransaction(db, func(tx *sqlx.Tx) error {
		for _, migration := range downMigrations {
			log.Printf("migrate: running down migration %s...\n", migration.Name)
			if err := runMigration(ctx, tx, migration.SQL); err != nil {
				return err
			}
		}
		return nil
	})
}

func Up(ctx context.Context, db *sqlx.DB) error {
	dir := filepath.Join(".", "migrations")
	return UpFrom(ctx, db, dir)
}

func UpFrom(ctx context.Context, db *sqlx.DB, dir string) error {
	upMigrations, err := loadMigrationsFrom(ctx, dir, "*.up.sql")
	if err != nil {
		return err
	}

	if len(upMigrations) == 0 {
		log.Println("migrate: no up migrations found")
		return nil
	}

	return runInTransaction(db, func(tx *sqlx.Tx) error {
		for _, migration := range upMigrations {
			log.Printf("migrate: running up migration %s...\n", migration.Name)
			if err := runMigration(ctx, tx, migration.SQL); err != nil {
				return err
			}
		}
		return nil
	})
}

func loadMigrationsFrom(ctx context.Context, dir, glob string) ([]Migration, error) {
	paths, err := filepath.Glob(filepath.Join(dir, glob))
	if err != nil {
		return nil, err
	}

	migrations := make([]Migration, 0, len(paths))
	for _, path := range paths {
		migration, err := loadMigration(ctx, path)
		if err != nil {
			return nil, err
		}
		migrations = append(migrations, migration)
	}

	return migrations, nil
}

func loadMigration(ctx context.Context, path string) (Migration, error) {
	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return Migration{}, err
	}
	return Migration{
		Name: filepath.Base(path),
		SQL:  string(fileBytes),
	}, nil
}

func runInTransaction(db *sqlx.DB, f func(tx *sqlx.Tx) error) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	if err := f(tx); err != nil {
		if rberr := tx.Rollback(); rberr != nil {
			log.Println("migrate: rollback failed:", err)
		}
		return err
	}

	return tx.Commit()
}

func runMigration(ctx context.Context, tx *sqlx.Tx, script string) error {
	if _, err := tx.ExecContext(ctx, script); err != nil {
		return err
	}
	return nil
}

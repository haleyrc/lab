package migrate_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/haleyrc/cheevos-simple/api/internal/migrate"
	"github.com/haleyrc/cheevos-simple/api/internal/testutil"
	"github.com/haleyrc/cheevos-simple/api/lib/check"
)

func TestMigrations(t *testing.T) {
	ctx := context.Background()
	check := check.New(t)

	db, cleanup, err := testutil.ConnectToTestDatabase()
	check.OK(err).Fatal()
	defer cleanup()

	dir := filepath.Join("..", "..", "migrations")

	err = migrate.UpFrom(ctx, db, dir)
	check.OK(err)
	check.Equals(countTables(ctx, db), 7)
	if err := migrate.Up(ctx, db); err != nil {
		t.Error("up migrations are not idempotent. got error:", err)
	}

	err = migrate.DownFrom(ctx, db, dir)
	check.OK(err)
	check.Equals(countTables(ctx, db), 0)
	if err := migrate.Down(ctx, db); err != nil {
		t.Error("down migrations are not idempotent. got error:", err)
	}
}

func countTables(ctx context.Context, db *sqlx.DB) int {
	var count int
	err := db.GetContext(
		ctx,
		&count,
		`
		SELECT
			count(*)
		FROM
			pg_catalog.pg_tables
		WHERE
			schemaname != 'pg_catalog' AND
			schemaname != 'information_schema';
	`)
	if err != nil {
		return -1
	}
	return count
}

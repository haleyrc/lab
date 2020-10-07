package db

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/haleyrc/cheevos-simple/api/internal/migrate"
	"github.com/haleyrc/cheevos-simple/api/internal/testutil"
)

var store *Store

func TestMain(m *testing.M) {
	db, cleanup, err := testutil.ConnectToTestDatabase()
	if err != nil {
		fmt.Println("failed to setup tests:", err)
		os.Exit(1)
	}
	store = &Store{db: db}

	dir := filepath.Join("..", "..", "migrations")
	if err := migrate.UpFrom(context.Background(), db, dir); err != nil {
		fmt.Println("failed to initialize database:", err)
		cleanup()
		os.Exit(1)
	}

	code := m.Run()

	cleanup()

	os.Exit(code)
}

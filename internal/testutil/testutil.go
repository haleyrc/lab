package testutil

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest"
)

func ConnectToTestDatabase() (*sqlx.DB, func(), error) {
	var db *sqlx.DB

	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, func() {}, err
	}

	log.Println("testutil: creating pool...")
	resource, err := pool.Run("postgres", "11.6", nil)
	if err != nil {
		return nil, func() {}, err
	}

	url := fmt.Sprintf(
		"postgres://postgres@localhost:%s/postgres?sslmode=disable",
		resource.GetPort("5432/tcp"),
	)
	log.Printf("testutil: attempting to connect to %s...\n", url)

	if err := pool.Retry(func() error {
		var err error
		db, err = sqlx.Connect("postgres", url)
		return err
	}); err != nil {
		pool.Purge(resource)
		return nil, func() {}, err
	}

	cleanup := func() {
		if err := db.Close(); err != nil {
			log.Println("testutil: could not close database:", err)
		}
		if err := pool.Purge(resource); err != nil {
			log.Println("testutil: could not purge resource:", err)
		}
	}

	return db, cleanup, nil
}

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/haleyrc/cheevos-simple/api/internal/db"
	"github.com/haleyrc/cheevos-simple/api/internal/migrate"
	"github.com/haleyrc/cheevos-simple/api/internal/web"
	"github.com/haleyrc/cheevos-simple/api/pkg/domain"
)

func main() {
	ctx := context.Background()

	cfg, err := loadConfig()
	if err != nil {
		log.Println("main: failed to load config:", err)
		os.Exit(1)
	}

	pgdb, err := connectWithRetry(cfg.DatabaseURL)
	if err != nil {
		log.Println("main: failed to connect to database:", err)
		os.Exit(1)
	}
	defer pgdb.Close()

	if err := migrate.Up(ctx, pgdb); err != nil {
		log.Println("main: migration failed:", err)
		os.Exit(1)
	}

	store := db.New(db.Config{
		DB: pgdb,
	})

	dom := domain.New(domain.Config{
		Database: store,
	})

	router := mux.NewRouter()
	router.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
		migrate.Down(ctx, pgdb)
		migrate.Up(ctx, pgdb)
	})

	web.NewServer(web.Config{
		Domain:     dom,
		Router:     router.PathPrefix("/api").Subrouter(),
		SigningKey: cfg.SigningKey,
	})

	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: cors(logRequests(router)),
	}
	log.Printf("main: listening on %s...\n", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Println("main: server quit unexpectedly:", err)
		os.Exit(1)
	}
}

func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := httputil.DumpRequest(r, true)
		fmt.Println(string(b))
		next.ServeHTTP(w, r)
	})
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == http.MethodOptions {
			return
		}
		next.ServeHTTP(w, r)
	})
}

func connectWithRetry(url string) (*sqlx.DB, error) {
	var db *sqlx.DB
	var err error
	for i := 0; i < 5; i++ {
		db, err = sqlx.Connect("postgres", url)
		if err == nil {
			return db, nil
		}
		dur := time.Duration(i) * time.Second
		log.Printf("database connection failed: trying again in %s...\n", dur)
		<-time.After(dur)
	}

	return nil, err
}

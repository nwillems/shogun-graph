package main

import (
	"context"
	gosql "database/sql"
	_ "embed"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nwillems/shogun-graph/shogun"
	"github.com/nwillems/shogun-graph/shogun/store/sql"
)

func run() error {
	ctx := context.Background()

	db, err := gosql.Open("sqlite3", ":memory:")
	if err != nil {
		return err
	}
	store := sql.NewStore(db)

	// Run migrations
	err = store.RunMigration(ctx)
	if err != nil {
		return err
	}
	log.Println("Migration success")

	// start server
	mux := http.NewServeMux()
	httpServer := &http.Server{
		Addr:    ":9000",
		Handler: mux,
	}
	server := &shogun.Server{
		ShutdownTimeout: 10 * time.Second,
		Log:             log.Default(),
		Server:          httpServer,
	}

	api := shogun.NewAPI(store)
	err = api.Register(mux)
	if err != nil {
		return err
	}

	go server.WaitForExit(ctx)
	if err := server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

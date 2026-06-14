package main

import (
	"log"
	"net/http"
	"time"

	"basalam-feedback/internal/config"
	"basalam-feedback/internal/repository"
	"basalam-feedback/internal/server"
)

func main() {
	cfg := config.Load()
	db, err := repository.OpenWithRetry(cfg.DatabaseURL, 3, 2*time.Second)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer db.Close()

	if err := repository.EnsureSchema(db); err != nil {
		log.Fatalf("schema initialization failed: %v", err)
	}

	handler := server.New(cfg, db)
	log.Printf("listening on %s", cfg.Addr())
	if err := http.ListenAndServe(cfg.Addr(), handler); err != nil {
		log.Fatalf("server exited: %v", err)
	}
}

package main

import (
	"log"
	"lumiiam/internal/config"
	"lumiiam/internal/db"
	"lumiiam/internal/router"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	db_conn, err := db.Open(cfg)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}

	r := router.New(cfg, db_conn)
	if err := r.Run(cfg.Addr()); err != nil {
		log.Fatalf("server run: %v", err)
	}
}

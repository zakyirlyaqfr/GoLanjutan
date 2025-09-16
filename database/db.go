package database

import (
	"database/sql"
	"fmt"
	"log"

	"golanjutan/config"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {
	cfg := config.AppEnv
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed open db: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("failed ping db: %v", err)
	}
}

func ConnectDB() {
	cfg := config.AppEnv
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed open db: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("failed ping db: %v", err)
	}
}

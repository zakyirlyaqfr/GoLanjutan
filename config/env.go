package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	ServerPort string
	JWTSecret  string
}

var AppEnv *Env

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, will use system env")
	}

	AppEnv = &Env{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "dbminggu4"),
		ServerPort: getEnv("SERVER_PORT", "3000"),
		JWTSecret:  getEnv("JWT_SECRET", "golanjutan"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

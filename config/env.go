package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	// DBHost     string // Tidak dipakai
	// DBPort     string // Tidak dipakai
	// DBUser     string // Tidak dipakai
	// DBPassword string // Tidak dipakai
	DBName     string
	MongoURI   string
	ServerPort string
	JWTSecret  string
}

var AppEnv *Env

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, will use system env")
	}

	AppEnv = &Env{
		MongoURI:   getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		DBName:     getEnv("DB_NAME", "dbgolanjutan"),
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
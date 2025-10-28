package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"golanjutan/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database
var Client *mongo.Client // Menyimpan client untuk session (transaksi)

func Connect() {
	cfg := config.AppEnv
	mongoURI := cfg.MongoURI
	dbName := cfg.DBName

	clientOptions := options.Client().ApplyURI(mongoURI)

	// Membuat konteks dengan batas waktu
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	Client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Koneksi ke MongoDB gagal: %v", err)
	}

	// Cek koneksi (Ping)
	err = Client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Ping ke MongoDB gagal: %v", err)
	}

	fmt.Println("Berhasil terhubung ke MongoDB!")
	DB = Client.Database(dbName)
}

// ConnectDB tidak diperlukan lagi karena Connect() sudah melakukannya
// func ConnectDB() { ... }
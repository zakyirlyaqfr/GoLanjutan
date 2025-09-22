package main

import (
	"golanjutan/config"
	"golanjutan/utils"
	"golanjutan/database"
	"golanjutan/route"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// 1. Load environment variables terlebih dahulu
	config.LoadEnv()

	// 2. Inisialisasi logger (opsional, tergantung kebutuhan)
	config.InitLogger()

	// 3. Inisialisasi JWT secret setelah env diload
	utils.InitJWT(config.AppEnv.JWTSecret)

	// 4. Koneksi ke database
	database.Connect()

	// 5. Buat Fiber app
	app := fiber.New(config.NewFiberConfig())

	// 6. Setup semua route
	route.Setup(app)

	// 7. Jalankan server
	port := config.AppEnv.ServerPort
	if port == "" {
		port = "8080" // fallback default port
	}
	log.Fatal(app.Listen(":" + port))
}

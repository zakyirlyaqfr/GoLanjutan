package main

import (
	"golanjutan/config"
	"golanjutan/database"
	"golanjutan/route"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// load env & init
	config.LoadEnv()
	config.InitLogger()

	// connect db
	database.Connect()

	// create fiber app
	app := fiber.New(config.NewFiberConfig())

	// setup routes
	route.Setup(app)

	// run
	port := config.AppEnv.ServerPort
	log.Fatal(app.Listen(":" + port))
}

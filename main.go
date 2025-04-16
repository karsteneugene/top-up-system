package main

import (
	"log"

	api "github.com/karsteneugene/top-up-system/api/routes"
	_ "github.com/karsteneugene/top-up-system/docs"

	"github.com/gofiber/swagger"
)

// @title Top Up System API
// @version 1.0
// @description This is a sample top-up system server.
// @host localhost:3000
// @BasePath /api
func main() {
	app := api.Api()

	// Swagger documentation
	app.Get("/swagger/*", swagger.HandlerDefault)

	log.Fatal(app.Listen(":3000"))
}

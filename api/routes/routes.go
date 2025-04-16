package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/karsteneugene/top-up-system/api/handlers"
)

func Api() *fiber.App {
	app := fiber.New()

	api := app.Group("/api")

	api.Get("/users", handlers.GetUsers)

	return app
}

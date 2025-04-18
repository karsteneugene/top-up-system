package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/karsteneugene/top-up-system/api/handlers"
)

func Api() *fiber.App {
	app := fiber.New()

	// Define base API route
	api := app.Group("/api")

	// User routes
	api.Get("/users", handlers.GetAllUsers)
	api.Get("/users/:id", handlers.GetUserByID)

	// Wallet routes
	api.Get("/wallets", handlers.GetAllWallets)
	api.Get("/wallets/:id", handlers.GetWalletByID)
	api.Get("/wallets/user/:id", handlers.GetWalletByUserID)
	api.Get("/wallets/va/:id", handlers.GetVirtualAccountByWalletID)

	// Transaction routes
	api.Get("/transactions/wallet/:id", handlers.GetTransactionsByWalletID)
	api.Post("/transactions/topup/direct/:id", handlers.TopUpDirect)
	api.Post("/transactions/topup/bank/:va", handlers.TopUpBank)

	return app
}

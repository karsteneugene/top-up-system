package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/karsteneugene/top-up-system/api/models"
)

// GetAllWallets godoc
// @Summary Get all wallets
// @Description Get all wallets
// @Tags wallets
// @Produce json
// @Success 200 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /wallets [get]
func GetAllWallets(c *fiber.Ctx) error {
	var wallets []models.Wallet

	// Get all wallets and check if there are any problems with the database
	if err := db.Find(&wallets).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Error retrieving wallets"})
	}

	return c.JSON(fiber.Map{"success": true, "payload": fiber.Map{"wallets": wallets}})
}

// GetWalletByID godoc
// @Summary Get wallet by ID
// @Description Get wallet by ID
// @Tags wallets
// @Produce json
// @Param id path int true "Wallet ID"
// @Success 200 {string} string
// @Failure 404 {string} string
// @Router /wallets/{id} [get]
func GetWalletByID(c *fiber.Ctx) error {
	var wallet models.Wallet
	id := c.Params("id")

	// Query wallet by ID
	if err := db.First(&wallet, id).Error; err != nil { // Check is wallet exists
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Wallet not found"})
	}
	return c.JSON(fiber.Map{"success": true, "payload": fiber.Map{"wallet": wallet}})
}

// GetWalletByUserID godoc
// @Summary Get wallet by User ID
// @Description Get wallet by User ID
// @Tags wallets
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {string} string
// @Failure 404 {string} string
// @Router /wallets/user/{id} [get]
func GetWalletByUserID(c *fiber.Ctx) error {
	var wallet models.Wallet
	id := c.Params("id")

	// Query wallet by User ID
	if err := db.Where("user_id = ?", id).First(&wallet).Error; err != nil { // Check if wallet exists
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Wallet not found"})
	}
	return c.JSON(fiber.Map{"success": true, "payload": fiber.Map{"wallet": wallet}})
}

// GetVirtualAccountByWalletID godoc
// @Summary Get virtual account by Wallet ID
// @Description Get virtual account by Wallet ID
// @Tags wallets
// @Produce json
// @Param id path int true "Wallet ID"
// @Success 200 {string} string
// @Failure 404 {string} string
// @Router /wallets/va/{id} [get]
func GetVirtualAccountByWalletID(c *fiber.Ctx) error {
	var wallet models.Wallet
	id := c.Params("id")

	// Query virtual account number by Wallet ID
	if err := db.First(&wallet, id).Error; err != nil { // Check if wallet exists
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Wallet not found"})
	}
	return c.JSON(fiber.Map{"success": true, "payload": fiber.Map{"virtual_account": wallet.VirtualAccount}})
}

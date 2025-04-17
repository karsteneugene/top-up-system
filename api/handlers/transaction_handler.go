package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/karsteneugene/top-up-system/api/models"
	"github.com/karsteneugene/top-up-system/utils"
)

// GetTransactionsByWalletID godoc
// @Summary Get transactions by wallet ID
// @Description Get transactions by wallet ID
// @Tags transactions
// @Produce json
// @Param id path int true "Wallet ID"
// @Success 200 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /transactions/wallet/{id} [get]
func GetTransactionsByWalletID(c *fiber.Ctx) error {
	id := c.Params("id")
	var transactions []models.Transaction

	// Check if there are transactions for the given wallet ID
	if err := db.Where("wallet_id = ?", id).Find(&transactions).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "No transactions found"})
	}
	return c.JSON(fiber.Map{"success": true, "payload": fiber.Map{"transactions": transactions}})
}

// TopUpDirect godoc
// @Summary Top up wallet directly
// @Description Top up wallet directly
// @Tags transactions
// @Accept json
// @Produce json
// @Param id path int true "Wallet ID"
// @Param amount body models.TransactionAmount true "Amount to top up"
// @Success 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /transactions/topup/direct/{id} [post]
func TopUpDirect(c *fiber.Ctx) error {

	amount := new(models.TransactionAmount)

	transaction := new(models.Transaction)

	// Parse request body to get amount
	if err := c.BodyParser(amount); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Invalid request"})
	}

	// Get wallet ID from URL parameter
	transaction.WalletID, _ = strconv.Atoi(c.Params("id"))

	// Set transaction amount
	transaction.Amount = amount.Amount

	// Set transaction type to DIRECT
	transaction.Type = models.TransactionTopUpDirect

	// Check if amount is more than minimum top up amount
	validMinMax, err := utils.CheckMinMaxTopUp(transaction.Amount)
	if !validMinMax {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": err})
	}

	// Check if amount is less than daily limit
	validDaily, err := utils.CheckDailyLimit(transaction.Amount, transaction.WalletID)
	if !validDaily {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": err})
	}

	// Check if amount is less than monthly limit
	validMonthly, err := utils.CheckMonthlyLimit(transaction.Amount, transaction.WalletID)
	if !validMonthly {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": err})
	}

	// Create transaction
	if err := db.Create(&transaction).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Error creating transaction"})
	}

	// Update wallet balance
	var wallet models.Wallet

	// Check if wallet exists
	if err := db.First(&wallet, transaction.WalletID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Wallet not found"})
	}

	// Add amount to wallet balance
	wallet.Balance += transaction.Amount

	// Save updated balance
	if err := db.Save(&wallet).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Error updating wallet"})
	}
	return c.JSON(fiber.Map{"success": true, "payload": fiber.Map{"transaction": transaction}})
}

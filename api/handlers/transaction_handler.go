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
	var transactions []models.Transaction
	id := c.Params("id")

	// Query transactions by wallet ID
	if err := db.Where("wallet_id = ?", id).Find(&transactions).Error; err != nil { // Checks if wallet exists
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Wallet not found"})
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
	wallet := new(models.Wallet)

	transaction.WalletID, _ = strconv.Atoi(c.Params("id"))

	// Get amount from request body
	if err := c.BodyParser(amount); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Invalid request"})
	}

	transaction.Amount = amount.Amount
	transaction.Type = models.TransactionTopUpDirect

	// Query wallet by given ID
	if err := db.First(&wallet, transaction.WalletID).Error; err != nil { // Check if wallet exists
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Wallet not found"})
	}

	// Check amount >= minimum or < maximum top up amount
	validMinMax, err := utils.CheckMinMaxTopUp(transaction.Amount)
	if !validMinMax {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": err})
	}

	// Check amount < daily limit
	validDaily, err := utils.CheckDailyLimit(transaction.Amount, transaction.WalletID)
	if !validDaily {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": err})
	}

	// Check amount < monthly limit
	validMonthly, err := utils.CheckMonthlyLimit(transaction.Amount, transaction.WalletID)
	if !validMonthly {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": err})
	}

	// Create transaction
	if err := db.Create(&transaction).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Error creating transaction"})
	}

	// Add amount to wallet balance
	wallet.Balance += transaction.Amount

	// Update balance
	if err := db.Save(&wallet).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Error updating wallet"})
	}
	return c.JSON(fiber.Map{"success": true, "payload": fiber.Map{"transaction": transaction}})
}

// TopUpBank godoc
// @Summary Top up wallet via bank transfer
// @Description Top up wallet via bank transfer
// @Tags transactions
// @Accept json
// @Produce json
// @Param va path int true "Virtual Account Number"
// @Param bank body models.BankTransactionRequest true "Bank transaction request"
// @Success 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /transactions/topup/bank/{va} [post]
func TopUpBank(c *fiber.Ctx) error {

	request := new(models.BankTransactionRequest)
	wallet := new(models.Wallet)
	transaction := new(models.Transaction)

	va := c.Params("va")

	// Get request body
	if err := c.BodyParser(request); err != nil {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": "Invalid request"})
	}

	// Query wallet by virtual account number
	if err := db.Where("virtual_account = ?", va).First(&wallet).Error; err != nil { // Check if virtual account exists
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "Virtual account not found"})
	}

	// Validate bank code
	valid, resBank := utils.ValidateBank(request.BankCode)
	if !valid {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": resBank})
	}

	// Validate account number
	valid, resAccount := utils.ValidateAccountNumber(request.AccountNumber)
	if !valid {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": resAccount})
	}

	transaction.Amount = request.Amount
	transaction.Type = models.TransactionTopUpBank
	transaction.RecipientBank = resBank
	transaction.RecipientName = resAccount
	transaction.WalletID = wallet.ID
	transaction.Description = request.Description

	// Check amount >= minimum or < maximum top up amount
	validMinMax, err := utils.CheckMinMaxTopUp(transaction.Amount)
	if !validMinMax {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": err})
	}

	// Check amount < daily limit
	validDaily, err := utils.CheckDailyLimit(transaction.Amount, transaction.WalletID)
	if !validDaily {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": err})
	}

	// Check amount < monthly limit
	validMonthly, err := utils.CheckMonthlyLimit(transaction.Amount, transaction.WalletID)
	if !validMonthly {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": err})
	}

	// Create transaction
	if err := db.Create(&transaction).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Error creating transaction"})
	}

	// Add amount to wallet balance
	wallet.Balance += transaction.Amount

	// Update balance
	if err := db.Save(&wallet).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"success": false, "message": "Error updating wallet"})
	}
	return c.JSON(fiber.Map{"success": true, "payload": fiber.Map{"transaction": transaction}})

}

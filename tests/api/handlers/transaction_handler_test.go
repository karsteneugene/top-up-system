package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/karsteneugene/top-up-system/api/handlers"
	"github.com/karsteneugene/top-up-system/api/models"
	"github.com/karsteneugene/top-up-system/setting"
	"github.com/karsteneugene/top-up-system/utils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var testDB *gorm.DB

func setupTestDB() (*gorm.DB, error) {
	// Initialize in-memory SQLite database
	var err error
	testDB, err = setting.Connect(":memory:")
	if err != nil {
		panic("failed to connect to database")
	}

	// Migrate the schema
	err = testDB.AutoMigrate(&models.User{}, &models.Wallet{}, &models.Transaction{})
	if err != nil {
		return nil, err
	}

	return testDB, nil
}

func TestMain(m *testing.M) {
	// Setup test database
	var err error
	testDB, err = setupTestDB()
	if err != nil {
		panic("failed to connect to database")
	}

	// Set the database connection for handlers
	handlers.SetDB(testDB)
	utils.SetDB(testDB)

	// Run tests
	code := m.Run()

	// Teardown database
	sqlDB, err := testDB.DB()
	if err != nil {
		panic("failed to get database connection")
	}
	sqlDB.Close()

	os.Exit(code)
}

func TestTopUpDirect(t *testing.T) {
	// Initialize Fiber for testing
	app := fiber.New()
	api := app.Group("/api")
	api.Post("/transactions/topup/direct/:id", handlers.TopUpDirect)

	// Seed users, wallets, and transactions for testing
	user1 := models.User{ // To test: successful top up, wallet not found, minimum and maximum top up amount
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
	}
	user2 := models.User{ // To test: daily limit
		ID:        2,
		FirstName: "Jane",
		LastName:  "Smith",
	}
	user3 := models.User{ // To test: monthly limit
		ID:        3,
		FirstName: "Alice",
		LastName:  "Johnson",
	}

	wallet1 := models.Wallet{
		ID:             1,
		Balance:        0,
		VirtualAccount: 1234567890,
		UserID:         1,
	}
	wallet2 := models.Wallet{
		ID:             2,
		Balance:        5000000,
		VirtualAccount: 9876543210,
		UserID:         2,
	}
	wallet3 := models.Wallet{
		ID:             3,
		Balance:        20000000,
		VirtualAccount: 1122334455,
		UserID:         3,
	}

	transactionDailyLimit := models.Transaction{
		Amount:   5000000,
		Type:     models.TransactionTopUpDirect,
		WalletID: 2,
	}
	transactionMonthlyLimit := models.Transaction{
		Amount:    20000000,
		CreatedAt: time.Now().AddDate(0, 0, -1),
		Type:      models.TransactionTopUpDirect,
		WalletID:  3,
	}

	testDB.Create(&user1)
	testDB.Create(&user2)
	testDB.Create(&user3)
	testDB.Create(&wallet1)
	testDB.Create(&wallet2)
	testDB.Create(&wallet3)
	testDB.Create(&transactionDailyLimit)
	testDB.Create(&transactionMonthlyLimit)

	// Test case: Successful direct top up
	t.Run("Successful direct top up", func(t *testing.T) {
		// Request body
		requestBody := `{"amount": 50000}`

		// Create a new request
		req := httptest.NewRequest("POST", "/api/transactions/topup/direct/1", bytes.NewBufferString(requestBody))
		req.Header.Set("Content-Type", "application/json")

		// Perform the request
		resp, _ := app.Test(req)

		// Assert status code
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Parse response body
		var responseBody map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&responseBody)

		// Assert response body
		assert.True(t, responseBody["success"].(bool))
		assert.NotNil(t, responseBody["payload"])

		// Assert balance update
		var updatedWallet models.Wallet
		testDB.First(&updatedWallet, 1)
		assert.Equal(t, 50000, updatedWallet.Balance)

		// Assert transaction creation
		var transaction models.Transaction
		testDB.Where("wallet_id = ?", 1).First(&transaction)
		assert.Equal(t, 50000, transaction.Amount)
		assert.Equal(t, models.TransactionTopUpDirect, transaction.Type)
	})

	// Test case: Wallet not found
	t.Run("Wallet not found", func(t *testing.T) {
		// Request body
		requestBody := `{"amount": 50000}`

		// Create a new request
		req := httptest.NewRequest("POST", "/api/transactions/topup/direct/999", bytes.NewBufferString(requestBody))
		req.Header.Set("Content-Type", "application/json")

		// Perform the request
		resp, _ := app.Test(req)

		// Assert status code
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		// Parse response body
		var responseBody map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&responseBody)

		// Assert response body
		assert.False(t, responseBody["success"].(bool))
		assert.Equal(t, "Wallet not found", responseBody["message"].(string))
	})

	// Test case: Less than minimum top up amount
	t.Run("Less than minimum top up amount", func(t *testing.T) {
		// Request body
		requestBody := `{"amount": 500}`

		// Create a new request
		req := httptest.NewRequest("POST", "/api/transactions/topup/direct/1", bytes.NewBufferString(requestBody))
		req.Header.Set("Content-Type", "application/json")

		// Perform the request
		resp, _ := app.Test(req)

		// Assert status code
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		// Parse response body
		var responseBody map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&responseBody)

		// Assert response body
		assert.False(t, responseBody["success"].(bool))
		assert.Equal(t, "Amount is less than the minimum top up limit of Rp 1,000", responseBody["message"].(string))
	})

	// Test case: Exceeds maximum top up amount
	t.Run("Exceeds maximum top up amount", func(t *testing.T) {
		// Request body
		requestBody := `{"amount": 3000000}`

		// Create a new request
		req := httptest.NewRequest("POST", "/api/transactions/topup/direct/1", bytes.NewBufferString(requestBody))
		req.Header.Set("Content-Type", "application/json")

		// Perform the request
		resp, _ := app.Test(req)

		// Assert status code
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		// Parse response body
		var responseBody map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&responseBody)

		// Assert response body
		assert.False(t, responseBody["success"].(bool))
		assert.Equal(t, "Amount exceeds the maximum top up limit of Rp 2,000,000", responseBody["message"].(string))
	})

	// Test case: Daily limit exceeded
	t.Run("Daily limit exceeded", func(t *testing.T) {
		// Request body
		requestBody := `{"amount": 1000000}`

		// Create a new request
		req := httptest.NewRequest("POST", "/api/transactions/topup/direct/2", bytes.NewBufferString(requestBody))
		req.Header.Set("Content-Type", "application/json")

		// Perform the request
		resp, _ := app.Test(req)

		// Assert status code
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		// Parse response body
		var responseBody map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&responseBody)

		// Assert response body
		assert.False(t, responseBody["success"].(bool))
		assert.Equal(t, "Total daily transaction exceeds/will exceed the daily limit of Rp 5,000,000", responseBody["message"].(string))
	})

	// Test case: Monthly limit exceeded
	t.Run("Monthly limit exceeded", func(t *testing.T) {
		// Request body
		requestBody := `{"amount": 1000000}`

		// Create a new request
		req := httptest.NewRequest("POST", "/api/transactions/topup/direct/3", bytes.NewBufferString(requestBody))
		req.Header.Set("Content-Type", "application/json")

		// Perform the request
		resp, _ := app.Test(req)

		// Assert status code
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		// Parse response body
		var responseBody map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&responseBody)

		// Assert response body
		assert.False(t, responseBody["success"].(bool))
		assert.Equal(t, "Total monthly transaction exceeds/will exceed the monthly limit of Rp 20,000,000", responseBody["message"].(string))
	})

	// Cleanup
	testDB.Unscoped().Delete(&user1)
	testDB.Unscoped().Delete(&user2)
	testDB.Unscoped().Delete(&user3)
	testDB.Unscoped().Delete(&wallet1)
	testDB.Unscoped().Delete(&wallet2)
	testDB.Unscoped().Delete(&wallet3)
	testDB.Unscoped().Delete(&transactionDailyLimit)
	testDB.Unscoped().Delete(&transactionMonthlyLimit)

	// Teardown database
	sqlDB, err := testDB.DB()
	if err != nil {
		panic("failed to get database connection")
	}
	sqlDB.Close()
}

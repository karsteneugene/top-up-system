package utils

import (
	"time"

	"github.com/karsteneugene/top-up-system/api/models"
	"github.com/karsteneugene/top-up-system/setting"
	"gorm.io/gorm"
)

var (
	db *gorm.DB

	// Used for queries
	today               = time.Now().Format("2006-01-02")
	tomorrow            = time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	startOfCurrentMonth = time.Now().AddDate(0, 0, -time.Now().Day()+1).Format("2006-01-02")
	endOfCurrentMonth   = time.Now().AddDate(0, 1, -time.Now().Day()).Format("2006-01-02")
)

// Initialize database connection for rate limit functions
func init() {
	var err error
	db, err = setting.Connect("ewallet.db")
	if err != nil {
		panic("failed to connect to database")
	}
}

// Sets which database to use for rate limit functions (used for unit testing to switch database into memory)
func SetDB(database *gorm.DB) {
	db = database
}

// Per transaction limit
func CheckMinMaxTopUp(amount int) (bool, string) {

	// Check amount <= minimum top-up amount
	if amount <= 1000 {
		return false, "Amount is less than the minimum top up limit of Rp 1,000"
	}

	// Check amount > maximum top-up amount
	if amount > 2000000 {
		return false, "Amount exceeds the maximum top up limit of Rp 2,000,000"
	}
	return true, ""
}

// Daily limit (resets at start of day)
func CheckDailyLimit(amount int, walletId int) (bool, string) {
	transactionTotal := new(models.TransactionTotal)

	// Calculate amount transacted today
	if err := db.Table("transactions").Select("wallet_id", "sum(amount) as amount").Where("wallet_id = ? AND created_at >= ? AND created_at < ?", walletId, today, tomorrow).First(&transactionTotal).Error; err != nil {
		return false, "Failed to calculate daily limit"
	}

	// Calculate total daily transaction when adding the new transaction
	transactionTotal.Amount += amount

	// Check total transaction > daily limit
	if transactionTotal.Amount > 5000000 {
		return false, "Total daily transaction exceeds/will exceed the daily limit of Rp 5,000,000"
	}
	return true, ""
}

// MOnthly limit (resets at start of month)
func CheckMonthlyLimit(amount int, walletId int) (bool, string) {
	transactionTotal := new(models.TransactionTotal)

	// Calculate amount transacted this month
	if err := db.Table("transactions").Select("wallet_id", "sum(amount) as amount").Where("wallet_id = ? AND created_at >= ? AND created_at < ?", walletId, startOfCurrentMonth, endOfCurrentMonth).First(&transactionTotal).Error; err != nil {
		return false, "Failed to calculate monthly limit"
	}

	// Calculate total monthly transaction when adding the new transaction
	transactionTotal.Amount += amount

	// Check total transaction > monthly limit
	if transactionTotal.Amount > 20000000 {
		return false, "Total monthly transaction exceeds/will exceed the monthly limit of Rp 20,000,000"
	}
	return true, ""
}

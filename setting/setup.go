package setting

import (
	"github.com/karsteneugene/top-up-system/api/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var dbConn *gorm.DB

func Connect(dsn string) (*gorm.DB, error) {
	// Connect to the database
	dbConn, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return dbConn, nil
}

func AutoMigrate() (*gorm.DB, error) {
	// Migrate the schema
	err := dbConn.AutoMigrate(&models.User{}, &models.Transaction{}, &models.Wallet{})
	if err != nil {
		return nil, err
	}

	return dbConn, nil
}

package setting

import (
	"github.com/karsteneugene/top-up-system/api/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Database() *gorm.DB {
	// Connect to the database
	db, err := gorm.Open(sqlite.Open("ewallet.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(&models.User{}, &models.Transaction{}, &models.Wallet{})
	if err != nil {
		panic("failed to migrate database")
	}

	return db
}

package handlers

import (
	"github.com/karsteneugene/top-up-system/setting"
	"gorm.io/gorm"
)

var db *gorm.DB

// Initialize database connection for handlers
func init() {
	var err error
	db, err = setting.Connect("ewallet.db")
	if err != nil {
		panic("failed to connect to database")
	}
}

// Sets which database to use for handlers (used for unit testing to switch database into memory)
func SetDB(database *gorm.DB) {
	db = database
}

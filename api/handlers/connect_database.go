package handlers

import (
	"github.com/karsteneugene/top-up-system/setting"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	var err error
	db, err = setting.Connect("ewallet.db")
	if err != nil {
		panic("failed to connect to database")
	}
}

func SetDB(database *gorm.DB) {
	db = database
}

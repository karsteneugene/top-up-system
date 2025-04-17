package handlers

import (
	"log"

	"github.com/karsteneugene/top-up-system/setting"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	var err error
	db, err = setting.Connect()
	if err != nil {
		log.Panicln("Failed to connect to database:", err.Error())
	}
}

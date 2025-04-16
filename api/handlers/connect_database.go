package handlers

import (
	"github.com/karsteneugene/top-up-system/setting"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	db = setting.Database()
}

package models

import "time"

type Wallet struct {
	ID             int       `json:"id" gorm:"primaryKey"`
	Balance        int       `json:"balance" gorm:"not null"`
	VirtualAccount int       `json:"virtual_account" gorm:"not null"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	UserID         int       `json:"user_id" gorm:"not null"`
	User           User      `json:"-" gorm:"foreignKey:UserID;references:ID"`
}

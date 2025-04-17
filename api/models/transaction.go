package models

import "time"

type TransactionType string

const (
	TransactionTopUpDirect TransactionType = "DIRECT"
	TransactionTopUpBank   TransactionType = "BANK"
)

type Transaction struct {
	ID            int             `json:"id" gorm:"primaryKey;autoIncrement"`
	Amount        int             `json:"amount" gorm:"not null"`
	CreatedAt     time.Time       `json:"created_at" gorm:"autoCreateTime"`
	Type          TransactionType `json:"type" gorm:"not null"`
	RecipientBank string          `json:"recipient_bank"`
	RecipientName string          `json:"recipient_name"`
	Description   string          `json:"description"`
	WalletID      int             `json:"wallet_id" gorm:"not null"`
	Wallet        Wallet          `json:"-" gorm:"foreignKey:WalletID;references:ID"`
}

// Used for the body of direct top-up
type TransactionAmount struct {
	Amount int `json:"amount" gorm:"not null"`
}

// Used for getting the total transactions amount
type TransactionTotal struct {
	Amount   int `json:"amount" gorm:"not null"`
	WalletID int `json:"wallet_id" gorm:"not null"`
}

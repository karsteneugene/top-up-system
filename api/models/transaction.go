package models

import "time"

type TransactionType string

const (
	TransactionTopUpDirect TransactionType = "DIRECT"
	TransactionTopUpBank   TransactionType = "BANK"
)

type Transaction struct {
	ID          int             `json:"id" gorm:"primaryKey"`
	Amount      int             `json:"amount" gorm:"not null"`
	CreatedAt   time.Time       `json:"created_at" gorm:"autoCreateTime"`
	Type        TransactionType `json:"type" gorm:"not null"`
	Description string          `json:"description"`
	WalletID    int             `json:"wallet_id" gorm:"not null"`
	Wallet      Wallet          `json:"wallet" gorm:"foreignKey:WalletID;references:ID"`
}

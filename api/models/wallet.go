package models

type Wallet struct {
	ID      int  `json:"id" gorm:"primaryKey"`
	Balance int  `json:"balance" gorm:"not null"`
	UserID  int  `json:"user_id" gorm:"not null"`
	User    User `json:"user" gorm:"foreignKey:UserID;references:ID"`
}

package models

type User struct {
	ID        int    `json:"id" gorm:"primaryKey"`
	FirstName string `json:"first_name" gorm:"not null"`
	LastName  string `json:"last_name" gorm:"not null"`
}

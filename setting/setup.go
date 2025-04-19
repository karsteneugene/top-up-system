package setting

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var dbConn *gorm.DB

// Main function to connect database
func Connect(dsn string) (*gorm.DB, error) { // Takes in the SQLite database file path

	// Connect to the database
	dbConn, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Automigrate the models (Uncomment if tables don't exist)
	// err = dbConn.AutoMigrate(&models.User{}, &models.Wallet{}, &models.Transaction{})
	// if err != nil {
	// 	return nil, err
	// }

	// Seed users and wallets (Uncomment if needed and add more users or wallets if necessary))
	// user1 := models.User{
	// 	ID:        1,
	// 	FirstName: "John",
	// 	LastName:  "Doe",
	// }
	// user2 := models.User{
	// 	ID:        2,
	// 	FirstName: "Jane",
	// 	LastName:  "Smith",
	// }
	// user3 := models.User{
	// 	ID:        3,
	// 	FirstName: "Alice",
	// 	LastName:  "Johnson",
	// }

	// wallet1 := models.Wallet{
	// 	ID:             1,
	// 	Balance:        0,
	// 	VirtualAccount: 1234567890,
	// 	UserID:         1,
	// }
	// wallet2 := models.Wallet{
	// 	ID:             2,
	// 	Balance:        5000000,
	// 	VirtualAccount: 9876543210,
	// 	UserID:         2,
	// }
	// wallet3 := models.Wallet{
	// 	ID:             3,
	// 	Balance:        20000000,
	// 	VirtualAccount: 1122334455,
	// 	UserID:         3,
	// }

	// dbConn.Create(&user1)
	// dbConn.Create(&user2)
	// dbConn.Create(&user3)
	// dbConn.Create(&wallet1)
	// dbConn.Create(&wallet2)
	// dbConn.Create(&wallet3)

	return dbConn, nil
}

package utils

// Bank validation
func ValidateBank(bankCode string) (bool, string) {

	// List of valid bank codes and names
	validBanks := map[string]string{
		"014": "BCA",
		"009": "BNI",
		"008": "MANDIRI",
		"002": "BRI",
	}

	// Check if the provided bank code is in the list of valid banks
	bank, exists := validBanks[bankCode]
	if !exists {
		return false, "Invalid bank code"
	}
	return true, bank
}

// Bank account validation
func ValidateAccountNumber(accountNumber string) (bool, string) {

	// Check if the digits on account number is valid
	if len(accountNumber) < 10 || len(accountNumber) > 15 {
		return false, "Invalid account number length"
	}

	// List of valid bank accounts
	validAccounts := map[string]string{
		"1234567890": "JOHN DOE",
		"9876543210": "JANE DOE",
	}

	// Check if the provided account number is in the list of valid accounts
	name, exists := validAccounts[accountNumber]
	if !exists {
		return false, "Account number not found"
	}
	return true, name
}

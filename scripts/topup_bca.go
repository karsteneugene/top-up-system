package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func main() {
	var va int
	var amount int
	var description string

	fmt.Print("Enter virtual account number: ")
	fmt.Scan(&va)

	fmt.Print("Enter amount to transfer: ")
	fmt.Scan(&amount)

	fmt.Print("Enter description (leave blank for empty): ")
	fmt.Scanln(&description)

	body := map[string]interface{}{
		"account_number": "1234567890",
		"amount":         amount,
		"bank_code":      "014",
		"description":    description,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}
	endpoint := "http://localhost:3000/api/transactions/topup/bank/" + fmt.Sprint(va)

	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println("Error making request:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println("Error decoding response:", err)
		return
	}

	fmt.Println("Successfully transferred")
}

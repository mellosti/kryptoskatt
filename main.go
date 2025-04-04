package main

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func process(api ExchangeApi, startTime int64, endTime int64) {
	// Get withdraw history
	withdrawals, err := api.GetWithdrawHistory(startTime, endTime)
	if err != nil {
		fmt.Printf("Error fetching withdraw history: %v\n", err)
	} else {
		fmt.Printf("Found %d withdrawals\n", len(withdrawals))
	}

	// Get deposit history
	deposits, err := api.GetDepositHistory(startTime, endTime)
	if err != nil {
		fmt.Printf("Error fetching deposit history: %v\n", err)
	} else {
		fmt.Printf("Found %d deposits\n", len(deposits))
	}

	// Get order history
	orders, err := api.GetOrderHistory(startTime, endTime)
	if err != nil {
		fmt.Printf("Error fetching order history: %v\n", err)
	} else {
		fmt.Printf("Found %d orders\n", len(orders))
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file" + err.Error())
	}

	// Define time range for the past 30 days
	endTime := time.Now().Unix()
	startTime := endTime - (30 * 24 * 60 * 60) // 30 days in seconds

	// Create a new OkxApi instance
	okxApi := &OkxApi{
		ApiKey:     os.Getenv("OKX_API_KEY"),
		SecretKey:  os.Getenv("OKX_API_SECRET"),
		Passphrase: os.Getenv("OKX_API_PASSPHRASE"),
		BaseUrl:    "https://my.okx.com",
	}

	process(okxApi, startTime, endTime)
}

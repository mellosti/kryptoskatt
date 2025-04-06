package main

import (
	"os"
	"time"

	"crypto-skatt-go/exchange"
	"crypto-skatt-go/okx"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file" + err.Error())
	}

	// Create a new OkxApi instance
	okxApi := &okx.OkxApiAdapter{
		ApiKey:     os.Getenv("OKX_API_KEY"),
		SecretKey:  os.Getenv("OKX_API_SECRET"),
		Passphrase: os.Getenv("OKX_API_PASSPHRASE"),
		BaseUrl:    "https://my.okx.com",
	}

	ExchangeService := &exchange.ExchangeService{
		Api: okxApi,
	}

	// Define time range for the past 30 days
	endTime := time.Now().Unix()
	startTime := endTime - (30 * 24 * 60 * 60) // 30 days in seconds

	ExchangeService.Process(startTime, endTime)
}

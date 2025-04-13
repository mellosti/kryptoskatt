package main

import (
	"time"

	"crypto-skatt-go/excel"
	"crypto-skatt-go/exchange"
	"crypto-skatt-go/okx"
)

func main() {

	// Create a new OkxApi instance
	okxApiAdapter := okx.NewOkxApiAdapter()
	excelService := excel.NewExcelService()
	exchangeService := &exchange.ExchangeService{
		Api:               okxApiAdapter,
		FileExportService: excelService,
	}

	// Define time range for the past 30 days
	endTime := time.Now().Unix()
	startTime := endTime - (30 * 24 * 60 * 60) // 30 days in seconds

	exchangeService.Process(startTime, endTime)
}

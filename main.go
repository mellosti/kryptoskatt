package main

import (
	"time"

	"crypto-skatt-go/excel"
	"crypto-skatt-go/exchange"
	"crypto-skatt-go/okx"
)

func main() {
	okxApiAdapter := okx.NewOkxApiAdapter()
	excelService := excel.NewExcelService(".")
	exchangeDataCollector := &exchange.ExchangeDataCollector{
		ExchangeAdapter:   okxApiAdapter,
		FileExportService: excelService,
	}

	endTime := time.Now().Unix()
	startTime := endTime - (30 * 24 * 60 * 60) // 30 days in seconds

	exchangeDataCollector.FetchAndExportExchangeData(startTime, endTime)
}

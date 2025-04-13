package exchange

import "fmt"

type ExchangeService struct {
	Api               ExchangeAdapter
	FileExportService FileExportService
}

type TransferHistory struct {
	Coin   string
	Amount string
}

type OrderHistory struct {
	BoughtCoin   string
	BoughtAmount float64
	SoldCoin     string
	SoldAmount   float64
	FeeAmount    float32
	FeeCurrency  string
	Timestamp    string
	Exchange     string
}

type ExchangeAdapter interface {
	GetWithdrawHistory(startTime int64, endTime int64) ([]TransferHistory, error)
	GetDepositHistory(startTime int64, endTime int64) ([]TransferHistory, error)
	GetOrderHistory(startTime int64, endTime int64) ([]OrderHistory, error)
}

type FileExportService interface {
	ExportToFile(withdrawals []TransferHistory, deposits []TransferHistory, orders []OrderHistory) error
}

func (es *ExchangeService) Process(startTime int64, endTime int64) error {
	// Get withdraw history
	withdrawals, err := es.Api.GetWithdrawHistory(startTime, endTime)
	if err != nil {
		return fmt.Errorf("error fetching withdraw history: %v", err)
	} else {
		fmt.Printf("Found %d withdrawals\n", len(withdrawals))
	}

	// Get deposit history
	deposits, err := es.Api.GetDepositHistory(startTime, endTime)
	if err != nil {
		return fmt.Errorf("error fetching deposit history: %w", err)
	} else {
		fmt.Printf("Found %d deposits\n", len(deposits))
	}

	// Get order history
	orders, err := es.Api.GetOrderHistory(startTime, endTime)
	if err != nil {
		return fmt.Errorf("error fetching order history: %w", err)
	} else {
		fmt.Printf("Found %d orders\n", len(orders))
	}

	if es.FileExportService.ExportToFile(withdrawals, deposits, orders) != nil {
		return fmt.Errorf("error exporting data to file")
	} else {
		fmt.Println("Data exported successfully")
	}

	return nil
}

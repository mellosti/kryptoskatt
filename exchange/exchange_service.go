package exchange

import "fmt"

type ExchangeService struct {
	Api ExchangeApi
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
}

type ExchangeApi interface {
	GetWithdrawHistory(startTime int64, endTime int64) ([]TransferHistory, error)
	GetDepositHistory(startTime int64, endTime int64) ([]TransferHistory, error)
	GetOrderHistory(startTime int64, endTime int64) ([]OrderHistory, error)
}

func (es *ExchangeService) Process(startTime int64, endTime int64) {
	// Get withdraw history
	withdrawals, err := es.Api.GetWithdrawHistory(startTime, endTime)
	if err != nil {
		fmt.Printf("Error fetching withdraw history: %v\n", err)
	} else {
		fmt.Printf("Found %d withdrawals\n", len(withdrawals))
	}

	// Get deposit history
	deposits, err := es.Api.GetDepositHistory(startTime, endTime)
	if err != nil {
		fmt.Printf("Error fetching deposit history: %v\n", err)
	} else {
		fmt.Printf("Found %d deposits\n", len(deposits))
	}

	// Get order history
	orders, err := es.Api.GetOrderHistory(startTime, endTime)
	if err != nil {
		fmt.Printf("Error fetching order history: %v\n", err)
	} else {
		fmt.Printf("Found %d orders\n", len(orders))
	}
}

package main

type TransferHistory struct {
	Coin   string
	Amount string
}

type OrderHistory struct {
	BoughtCoin    string
	BoughtAmount  float64
	SoldCoin      string
	SoldAmount    float64
	FeeAmount     float32
	FeeCurrency   string
	UnixTimestamp int64
}

type ExchangeApi interface {
	GetWithdrawHistory(startTime int64, endTime int64) ([]TransferHistory, error)
	GetDepositHistory(startTime int64, endTime int64) ([]TransferHistory, error)
	GetOrderHistory(startTime int64, endTime int64) ([]OrderHistory, error)
}

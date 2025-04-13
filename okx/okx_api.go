package okx

import (
	"crypto-skatt-go/exchange"
	"crypto-skatt-go/http"
	"crypto-skatt-go/util"
	"encoding/json"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type OkxApiAdapter struct {
	httpClient  *http.Client
	credentials Credentials
}

func NewOkxApiAdapter() *OkxApiAdapter {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file" + err.Error())
	}

	apiKey := os.Getenv("OKX_API_KEY")
	secretKey := os.Getenv("OKX_API_SECRET")
	passphrase := os.Getenv("OKX_API_PASSPHRASE")

	return &OkxApiAdapter{
		httpClient: http.NewClient("https://my.okx.com"),
		credentials: Credentials{
			ApiKey:     apiKey,
			SecretKey:  secretKey,
			Passphrase: passphrase,
		},
	}
}

type OkxDepositHistoryResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		Amt   string `json:"amt"`   // Amount
		Ccy   string `json:"ccy"`   // Currency
		State string `json:"state"` // Status (2 = completed)
		OrdId string `json:"ordId"` // Order ID
		Ts    string `json:"ts"`    // Timestamp
	} `json:"data"`
}

type OkxWithdrawalHistoryResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		Amt   string `json:"amt"`   // Amount
		Ccy   string `json:"ccy"`   // Currency
		Fee   string `json:"fee"`   // Fee amount
		TxId  string `json:"txId"`  // Transaction ID
		State string `json:"state"` // Status (2 = completed)
		Ts    string `json:"ts"`    // Timestamp
	} `json:"data"`
}

func (o *OkxApiAdapter) GetWithdrawHistory(startTime int64, endTime int64) ([]exchange.TransferHistory, error) {
	endpoint := "/api/v5/asset/withdrawal-history"
	queryParams := map[string]string{
		"limit": "100",
	}

	headers := o.credentials.GenerateHeaders(endpoint, GET, queryParams, nil)
	_, body, err := o.httpClient.Get(endpoint, queryParams, headers)
	if err != nil {
		return nil, err
	}

	var parsedResponse OkxWithdrawalHistoryResponse
	if err := json.Unmarshal(body, &parsedResponse); err != nil {
		return nil, err
	}

	withdrawals := []exchange.TransferHistory{}
	for _, withdrawal := range parsedResponse.Data {
		// Only include successful withdrawals (state 2)
		if withdrawal.State == "2" {
			withdrawals = append(withdrawals, exchange.TransferHistory{
				Coin:          withdrawal.Ccy,
				Amount:        withdrawal.Amt,
				Timestamp:     util.FormatTimestampToISO(withdrawal.Ts),
				TransactionID: withdrawal.TxId,
				FeeAmount:     util.ParseFloat32(withdrawal.Fee),
				FeeCoin:       withdrawal.Ccy,
				Exchange:      "okx",
				State:         withdrawal.State,
			})
		}
	}

	return withdrawals, nil
}

func (o *OkxApiAdapter) GetDepositHistory(startTime int64, endTime int64) ([]exchange.TransferHistory, error) {
	endpoint := "/api/v5/asset/deposit-history"
	queryParams := map[string]string{
		"limit": "100",
	}

	headers := o.credentials.GenerateHeaders(endpoint, GET, queryParams, nil)
	_, body, err := o.httpClient.Get(endpoint, queryParams, headers)
	if err != nil {
		return nil, err
	}

	var parsedResponse OkxDepositHistoryResponse
	if err := json.Unmarshal(body, &parsedResponse); err != nil {
		return nil, err
	}

	deposits := []exchange.TransferHistory{}
	for _, deposit := range parsedResponse.Data {
		// Only include successful deposits (state 2)
		if deposit.State == "2" {
			deposits = append(deposits, exchange.TransferHistory{
				Coin:          deposit.Ccy,
				Amount:        deposit.Amt,
				Timestamp:     util.FormatTimestampToISO(deposit.Ts),
				TransactionID: deposit.OrdId,
				FeeAmount:     0, // Deposits typically have no fees
				FeeCoin:       "",
				Exchange:      "okx",
				State:         deposit.State,
			})
		}
	}

	return deposits, nil
}

type OkxOrderHistoryResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		OrdId     string `json:"ordId"`
		State     string `json:"state"`
		Side      string `json:"side"`
		FillTime  string `json:"fillTime"`
		InstId    string `json:"instId"`
		AccFillSz string `json:"accFillSz"`
		AvgPx     string `json:"avgPx"`
		Fee       string `json:"fee"`
		FeeCcy    string `json:"feeCcy"`
	} `json:"data"`
}

func (o *OkxApiAdapter) GetOrderHistory(startTime int64, endTime int64) ([]exchange.OrderHistory, error) {
	endpoint := "/api/v5/trade/orders-history-archive"
	queryParams := map[string]string{
		// Must be this order for signature to work
		"instType": "SPOT",
		"limit":    "100",
	}

	headers := o.credentials.GenerateHeaders(endpoint, GET, queryParams, nil)
	_, body, err := o.httpClient.Get(endpoint, queryParams, headers)
	if err != nil {
		return nil, err
	}

	var parsedResponse OkxOrderHistoryResponse
	if err := json.Unmarshal(body, &parsedResponse); err != nil {
		return nil, err
	}

	orderHistory := []exchange.OrderHistory{}
	for _, order := range parsedResponse.Data {
		if order.State == "filled" || order.State == "partially_filled" {
			var boughtCoin string
			var soldCoin string
			var boughtAmount float64
			var soldAmount float64

			if order.Side == "buy" {
				boughtCoin = strings.Split(order.InstId, "-")[0]
				soldCoin = strings.Split(order.InstId, "-")[1]
				boughtAmount = util.ParseFloat64(order.AccFillSz)
				soldAmount = util.ParseFloat64(order.AccFillSz) * util.ParseFloat64(order.AvgPx)
			} else {
				boughtCoin = strings.Split(order.InstId, "-")[1]
				soldCoin = strings.Split(order.InstId, "-")[0]
				soldAmount = util.ParseFloat64(order.AccFillSz)
				boughtAmount = util.ParseFloat64(order.AccFillSz) * util.ParseFloat64(order.AvgPx)
			}

			orderHistory = append(orderHistory, exchange.OrderHistory{
				BoughtCoin:   boughtCoin,
				BoughtAmount: boughtAmount,
				SoldCoin:     soldCoin,
				SoldAmount:   soldAmount,
				FeeAmount:    util.ParseFloat32(order.Fee) * -1,
				FeeCurrency:  order.FeeCcy,
				Timestamp:    util.FormatTimestampToISO(order.FillTime),
				Exchange:     "okx",
			})
		}
	}

	return orderHistory, nil
}

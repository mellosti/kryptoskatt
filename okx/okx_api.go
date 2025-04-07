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
	ApiKey     string
	SecretKey  string
	Passphrase string
	httpClient *http.Client
	headers    Headers
}

func NewOkxApiAdapter() *OkxApiAdapter {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file" + err.Error())
	}

	apiKey := os.Getenv("OKX_API_KEY")
	secretKey := os.Getenv("OKX_API_SECRET")
	passphrase := os.Getenv("OKX_API_PASSPHRASE")

	return &OkxApiAdapter{
		ApiKey:     apiKey,
		SecretKey:  secretKey,
		Passphrase: passphrase,
		httpClient: http.NewClient("https://my.okx.com"),
		headers: Headers{
			ApiKey:     apiKey,
			SecretKey:  secretKey,
			Passphrase: passphrase,
		},
	}
}

func (o *OkxApiAdapter) GetWithdrawHistory(startTime int64, endTime int64) ([]exchange.TransferHistory, error) {
	// Implement the logic to fetch withdraw history from Okx API
	return []exchange.TransferHistory{}, nil
}

func (o *OkxApiAdapter) GetDepositHistory(startTime int64, endTime int64) ([]exchange.TransferHistory, error) {
	// Implement the logic to fetch deposit history from Okx API
	return []exchange.TransferHistory{}, nil
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
		"limit":    "100",
		"instType": "SPOT",
	}
	headers := o.headers.GetHeaders(endpoint, GET, queryParams, nil)
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
				FeeAmount:    util.ParseFloat32(order.Fee),
				FeeCurrency:  order.FeeCcy,
				Timestamp:    order.FillTime,
			})
		}
	}

	return orderHistory, nil
}

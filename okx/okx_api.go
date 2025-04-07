package okx

import (
	"crypto-skatt-go/exchange"
	"crypto-skatt-go/http"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type OkxApiAdapter struct {
	ApiKey     string
	SecretKey  string
	Passphrase string
	httpClient *http.Client
}

func NewOkxApiAdapter() *OkxApiAdapter {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file" + err.Error())
	}

	return &OkxApiAdapter{
		ApiKey:     os.Getenv("OKX_API_KEY"),
		SecretKey:  os.Getenv("OKX_API_SECRET"),
		Passphrase: os.Getenv("OKX_API_PASSPHRASE"),
		httpClient: http.NewClient("https://my.okx.com"),
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
	headers, err := o.getOkxHeaders(endpoint, "GET", queryParams, nil)
	if err != nil {
		return nil, err
	}

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
				boughtAmount = parseFloat64(order.AccFillSz)
				soldAmount = parseFloat64(order.AccFillSz) * parseFloat64(order.AvgPx)
			} else {
				boughtCoin = strings.Split(order.InstId, "-")[1]
				soldCoin = strings.Split(order.InstId, "-")[0]
				soldAmount = parseFloat64(order.AccFillSz)
				boughtAmount = parseFloat64(order.AccFillSz) * parseFloat64(order.AvgPx)
			}

			orderHistory = append(orderHistory, exchange.OrderHistory{
				BoughtCoin:   boughtCoin,
				BoughtAmount: boughtAmount,
				SoldCoin:     soldCoin,
				SoldAmount:   soldAmount,
				FeeAmount:    parseFloat32(order.Fee),
				FeeCurrency:  order.FeeCcy,
				Timestamp:    order.FillTime,
			})
		}
	}

	return orderHistory, nil
}

func (o *OkxApiAdapter) getOkxHeaders(url string, method string, queryParams map[string]string, bodyParams map[string]string) (map[string]string, error) {
	timestampIso := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	signature := Signature{
		Timestamp:   timestampIso,
		Method:      method,
		Endpoint:    url,
		SecretKey:   o.SecretKey,
		QueryParams: queryParams,
		BodyParams:  bodyParams,
	}.Encode()
	return map[string]string{
		"CONTENT-TYPE":         "application/json",
		"OK-ACCESS-KEY":        o.ApiKey,
		"OK-ACCESS-SIGN":       signature,
		"OK-ACCESS-TIMESTAMP":  timestampIso,
		"OK-ACCESS-PASSPHRASE": o.Passphrase,
	}, nil
}

func parseFloat64(s string) float64 {
	if s == "" {
		return 0
	}
	value, err := strconv.ParseFloat(s, 64)
	if err != nil {
		fmt.Printf("Error parsing '%s' as float: %v\n", s, err)
		return 0
	}
	return value
}

func parseFloat32(s string) float32 {
	if s == "" {
		return 0
	}
	value, err := strconv.ParseFloat(s, 32)
	if err != nil {
		fmt.Printf("Error parsing '%s' as float: %v\n", s, err)
		return 0
	}
	return float32(value)
}

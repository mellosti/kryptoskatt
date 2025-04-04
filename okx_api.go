package main

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

type OkxApi struct {
	ApiKey     string
	SecretKey  string
	Passphrase string
	BaseUrl    string
}

func (o *OkxApi) GetWithdrawHistory(startTime int64, endTime int64) ([]TransferHistory, error) {
	// Implement the logic to fetch withdraw history from Okx API
	return []TransferHistory{}, nil
}

func (o *OkxApi) GetDepositHistory(startTime int64, endTime int64) ([]TransferHistory, error) {
	// Implement the logic to fetch deposit history from Okx API
	return []TransferHistory{}, nil
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

func (o *OkxApi) GetOrderHistory(startTime int64, endTime int64) ([]OrderHistory, error) {
	endpoint := "/api/v5/trade/orders-history-archive"
	queryParams := map[string]string{
		"limit":    "100",
		"instType": "SPOT",
	}
	headers, err := o.getOkxHeaders(endpoint, "GET", queryParams, nil)
	if err != nil {
		return nil, err
	}

	response, err := Get(GetRequest{
		Url:         o.BaseUrl + endpoint,
		QueryParams: queryParams,
		Headers:     headers,
	})

	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var parsedResponse OkxOrderHistoryResponse
	if err := json.Unmarshal(body, &parsedResponse); err != nil {
		return nil, err
	}
	fmt.Println("Parsed response:", parsedResponse)

	return []OrderHistory{}, nil
}

func (o *OkxApi) getOkxHeaders(url string, method string, queryParams map[string]string, bodyParams map[string]string) (map[string]string, error) {
	timestampIso := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	signature, err := o.getOkxSignature(timestampIso, method, url, queryParams, bodyParams)
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"CONTENT-TYPE":         "application/json",
		"OK-ACCESS-KEY":        o.ApiKey,
		"OK-ACCESS-SIGN":       signature,
		"OK-ACCESS-TIMESTAMP":  timestampIso,
		"OK-ACCESS-PASSPHRASE": o.Passphrase,
	}, nil
}

func (o *OkxApi) getOkxSignature(timestampIso string, method string, endpoint string, queryParams map[string]string, bodyParams map[string]string) (string, error) {
	methodUpper := strings.ToUpper(method)

	var bodyString string = ""
	if bodyParams != nil {
		bytes, err := json.Marshal(bodyParams)
		if err != nil {
			fmt.Println("Error marshalling body params:", err)
			return "", err
		}
		bodyString = string(bytes)
	}

	var queryString string = ""
	if queryParams != nil {
		queryString = "?" + EncodeQueryParams(queryParams)
	}

	signatureString := fmt.Sprintf("%s%s%s%s%s", timestampIso, methodUpper, endpoint, queryString, bodyString)

	return GetHmac(signatureString, o.SecretKey), nil
}

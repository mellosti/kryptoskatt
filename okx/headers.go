package okx

import "time"

type Headers struct {
	ApiKey     string
	SecretKey  string
	Passphrase string
}

func (h Headers) GetHeaders(url string, method string, queryParams map[string]string, bodyParams map[string]string) map[string]string {
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	signature := Signature{
		Timestamp:   timestamp,
		Method:      method,
		Endpoint:    url,
		SecretKey:   h.SecretKey,
		QueryParams: queryParams,
		BodyParams:  bodyParams,
	}.Encode()
	return map[string]string{
		"CONTENT-TYPE":         "application/json",
		"OK-ACCESS-KEY":        h.ApiKey,
		"OK-ACCESS-SIGN":       signature,
		"OK-ACCESS-TIMESTAMP":  timestamp,
		"OK-ACCESS-PASSPHRASE": h.Passphrase,
	}
}

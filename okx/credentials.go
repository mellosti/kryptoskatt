package okx

import "time"

// Credentials stores API authentication information
type Credentials struct {
	ApiKey     string
	SecretKey  string
	Passphrase string
}

type Method string

const (
	GET    Method = "GET"
	POST   Method = "POST"
	PUT    Method = "PUT"
	DELETE Method = "DELETE"
	PATCH  Method = "PATCH"
)

// GenerateHeaders creates the required headers for OKX API authentication
func (c *Credentials) GenerateHeaders(url string, method Method, queryParams map[string]string, bodyParams map[string]string) map[string]string {
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	signature := signature{
		Timestamp:   timestamp,
		Method:      string(method),
		Endpoint:    url,
		SecretKey:   c.SecretKey,
		QueryParams: queryParams,
		BodyParams:  bodyParams,
	}.Encode()

	return map[string]string{
		"CONTENT-TYPE":         "application/json",
		"OK-ACCESS-KEY":        c.ApiKey,
		"OK-ACCESS-SIGN":       signature,
		"OK-ACCESS-TIMESTAMP":  timestamp,
		"OK-ACCESS-PASSPHRASE": c.Passphrase,
	}
}

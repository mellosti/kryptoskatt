package okx

import (
	"crypto-skatt-go/crypto"
	"crypto-skatt-go/http"
	"encoding/json"
	"fmt"
)

type signature struct {
	Timestamp   string
	Method      string
	Endpoint    string
	SecretKey   string
	QueryParams map[string]string
	BodyParams  map[string]string
}

func (s signature) Encode() string {
	var bodyString string
	if s.BodyParams != nil {
		bytes, err := json.Marshal(s.BodyParams)
		if err != nil {
			panic(("Error marshalling body params: " + err.Error()))
		}
		bodyString = string(bytes)
	}

	var queryString string
	if len(s.QueryParams) > 0 {
		queryString = "?" + http.EncodeQueryParams(s.QueryParams)
	}

	signatureString := fmt.Sprintf("%s%s%s%s%s", s.Timestamp, s.Method, s.Endpoint, queryString, bodyString)
	return crypto.GetHmac(signatureString, s.SecretKey)
}

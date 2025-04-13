package okx

import (
	"testing"
)

type test struct {
	name        string
	signature   signature
	expected    string
	description string
}

func TestSignatureEncode(t *testing.T) {
	tests := []test{
		{
			name: "Empty query params and body params",
			signature: signature{
				Timestamp:   "2023-04-01T12:00:00.000Z",
				Method:      "GET",
				Endpoint:    "/api/v5/account/balance",
				SecretKey:   "SEC123456789",
				QueryParams: map[string]string{},
				BodyParams:  nil,
			},
			expected:    "MIXqVEzXenuLCmYACL2eJ0TB31WvMOW5a9Yq5oHleUo=", // Replace with actual expected signature in real test
			description: "Should generate correct signature with no query params and no body",
		},
		{
			name: "With query params, no body params",
			signature: signature{
				Timestamp: "2023-04-01T12:00:00.000Z",
				Method:    "GET",
				Endpoint:  "/api/v5/account/balance",
				SecretKey: "SEC123456789",
				QueryParams: map[string]string{
					"ccy": "BTC",
				},
				BodyParams: nil,
			},
			expected:    "KtnE+TAZdFnf2NjgbeyYN45stQVGrn71bCtzLPRXzuE=", // Replace with actual expected signature in real test
			description: "Should generate correct signature with query params and no body",
		},
		{
			name: "With body params, no query params",
			signature: signature{
				Timestamp:   "2023-04-01T12:00:00.000Z",
				Method:      "POST",
				Endpoint:    "/api/v5/trade/order",
				SecretKey:   "SEC123456789",
				QueryParams: map[string]string{},
				BodyParams: map[string]string{
					"instId":  "BTC-USDT",
					"tdMode":  "cash",
					"side":    "buy",
					"ordType": "limit",
					"sz":      "1",
					"px":      "10000",
				},
			},
			expected:    "e3FqumMYKjAXrSKiLjMmoB1dy1mSRV5RLZNwHmiICd0=", // Replace with actual expected signature in real test
			description: "Should generate correct signature with body params and no query",
		},
		{
			name: "With both query params and body params",
			signature: signature{
				Timestamp: "2023-04-01T12:00:00.000Z",
				Method:    "POST",
				Endpoint:  "/api/v5/trade/order",
				SecretKey: "SEC123456789",
				QueryParams: map[string]string{
					"ccy": "BTC",
				},
				BodyParams: map[string]string{
					"instId": "BTC-USDT",
					"tdMode": "cash",
				},
			},
			expected:    "lX39yhEM0CaqlZKd2eIMqRIPhDHOjJMV4VNAp8bWOKg=", // Replace with actual expected signature in real test
			description: "Should generate correct signature with both query and body params",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.signature.Encode()
			// If you don't have expected signatures yet, you can print them and verify manually first
			t.Logf("Generated signature: %s", got)

			if got != tt.expected {
				t.Errorf("%s: expected signature %q, got %q", tt.description, tt.expected, got)
			}
		})
	}
}

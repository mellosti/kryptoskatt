package http

import "fmt"

func EncodeQueryParams(queryParams map[string]string) string {
	var queryString string
	for key, value := range queryParams {
		if queryString != "" {
			queryString += "&"
		}
		queryString += fmt.Sprintf("%s=%s", key, value)
	}
	return queryString
}

package main

import (
	"fmt"
	"io"
	"net/http"
)

type GetRequest struct {
	Url         string
	QueryParams map[string]string
	Headers     map[string]string
}

var client = &http.Client{}

func Get(request GetRequest) (*http.Response, error) {
	req, err := http.NewRequest("GET", request.Url, nil)
	if err != nil {
		return nil, err
	}

	addHeaders(req, request.Headers)
	addParamsToUrl(req, request.QueryParams)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		bodyString, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("error with code: %s, response body: %s", resp.Status, string(bodyString))
	}

	return resp, nil
}

func addHeaders(req *http.Request, headers map[string]string) {
	for headerKey, headerValue := range headers {
		req.Header.Add(headerKey, headerValue)
	}
}

func addParamsToUrl(req *http.Request, queryParams map[string]string) {
	if len(queryParams) == 0 {
		return
	}

	q := req.URL.Query()
	for paramKey, paramValue := range queryParams {
		q.Add(paramKey, paramValue)
	}
	req.URL.RawQuery = q.Encode()
}

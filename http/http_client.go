package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client wraps the standard http.Client with additional functionality
type Client struct {
	client         *http.Client
	baseURL        string
	defaultHeaders map[string]string
}

// RequestOptions contains all options for making HTTP requests
type RequestOptions struct {
	Method      string
	Path        string
	QueryParams map[string]string
	Headers     map[string]string
	Body        interface{}
	Context     context.Context
	Timeout     time.Duration
}

// NewClient creates a new HTTP client with sensible defaults
func NewClient(baseURL string) *Client {

	return &Client{
		client:  &http.Client{},
		baseURL: baseURL,
		defaultHeaders: map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
	}
}

// SetDefaultHeader sets a default header to be sent with every request
func (c *Client) SetDefaultHeader(key, value string) {
	c.defaultHeaders[key] = value
}

// Request performs an HTTP request and returns the response
func (c *Client) Request(opts RequestOptions) (*http.Response, []byte, error) {
	ctx := opts.Context
	if ctx == nil {
		ctx = context.Background()
	}

	// Apply timeout if provided
	var cancel context.CancelFunc
	if opts.Timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, opts.Timeout)
		defer cancel()
	}

	// Prepare URL
	url := c.baseURL
	if opts.Path != "" {
		url += opts.Path
	}

	// Prepare request body
	var bodyReader io.Reader
	if opts.Body != nil {
		jsonBody, err := json.Marshal(opts.Body)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewBuffer(jsonBody)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, opts.Method, url, bodyReader)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add default headers
	for key, value := range c.defaultHeaders {
		req.Header.Set(key, value)
	}

	// Add custom headers (overriding defaults if needed)
	for key, value := range opts.Headers {
		req.Header.Set(key, value)
	}

	// Add query parameters
	if len(opts.QueryParams) > 0 {
		q := req.URL.Query()
		for key, value := range opts.QueryParams {
			q.Add(key, value)
		}
		req.URL.RawQuery = q.Encode()
	}

	// Execute request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp, nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for non-2xx status codes
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return resp, body, fmt.Errorf("request failed with status %s: %s", resp.Status, string(body))
	}

	return resp, body, nil
}

// Get is a helper method for making GET requests
func (c *Client) Get(path string, queryParams map[string]string, headers map[string]string) (*http.Response, []byte, error) {
	return c.Request(RequestOptions{
		Method:      http.MethodGet,
		Path:        path,
		QueryParams: queryParams,
		Headers:     headers,
	})
}

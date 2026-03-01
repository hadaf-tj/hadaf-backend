package smsProvider

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultTimeoutSeconds = 20
)

type httpClient struct {
	baseURL string
	token   string
	login   string
	client  *http.Client
}

func newHTTPClient(baseURL, token, login string) *httpClient {
	return &httpClient{
		baseURL: baseURL,
		token:   token,
		login:   login,
		client: &http.Client{
			Timeout: defaultTimeoutSeconds * time.Second,
		},
	}
}

func (c *httpClient) doRequest(ctx context.Context, endpoint string, queryParams map[string]string) ([]byte, int, error) {
	// Build URL with query parameters
	reqURL := c.baseURL + endpoint
	if len(queryParams) > 0 {
		u, err := url.Parse(reqURL)
		if err != nil {
			return nil, 0, &NetworkError{Message: "failed to parse URL", Err: err}
		}

		q := u.Query()
		for k, v := range queryParams {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
		reqURL = u.String()
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, 0, &NetworkError{Message: "failed to create request", Err: err}
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, 0, &NetworkError{Message: "failed to execute request", Err: err}
	}
	defer resp.Body.Close()

	// Log response status (will be replaced with zerolog in step 7)

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, &NetworkError{Message: "failed to read response body", Err: err}
	}

	// Check for error status codes
	if resp.StatusCode >= 400 {
		return body, resp.StatusCode, parseAPIError(resp.StatusCode, body)
	}

	return body, resp.StatusCode, nil
}

package http

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	Client   *http.Client
	username *string
	password *string
}

func New(username, password *string, timeout time.Duration) *Client {
	return &Client{
		Client: &http.Client{
			Timeout: timeout,
		},
		username: username,
		password: password,
	}
}

func (c *Client) request(ctx context.Context, method, url string, headers map[string]string, body any) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if c.username != nil && c.password != nil {
		auth := base64.StdEncoding.EncodeToString([]byte(*c.username + ":" + *c.password))
		req.Header.Set("Authorization", "Basic "+auth)
	}

	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return c.Client.Do(req)
}

func (c *Client) Get(ctx context.Context, url string, headers map[string]string) (*http.Response, error) {
	return c.request(ctx, http.MethodGet, url, headers, nil)
}

func (c *Client) Post(ctx context.Context, url string, headers map[string]string, body any) (*http.Response, error) {
	return c.request(ctx, http.MethodPost, url, headers, body)
}

func (c *Client) Put(ctx context.Context, url string, headers map[string]string, body any) (*http.Response, error) {
	return c.request(ctx, http.MethodPut, url, headers, body)
}

func (c *Client) Delete(ctx context.Context, url string, headers map[string]string) (*http.Response, error) {
	return c.request(ctx, http.MethodDelete, url, headers, nil)
}

func IsOK(res *http.Response) bool {
	return res != nil &&
		res.StatusCode >= http.StatusOK && res.StatusCode < http.StatusMultipleChoices
}

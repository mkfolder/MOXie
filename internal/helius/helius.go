package helius

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/Makefolder/cynero/pkg/http"
)

type HeliusNet string

const (
	HeliusNetMainnet HeliusNet = "mainnet"
	HeliusNetDevnet  HeliusNet = "devnet"
)

type HeliusClient struct {
	http   *http.Client
	apiKey string
	domain string
}

func NewClient(http *http.Client, apiKey string, net HeliusNet) *HeliusClient {
	domain := "api-mainnet.helius-rpc.com"
	if net == HeliusNetDevnet {
		domain = "api-devnet.helius-rpc.com"
	}

	return &HeliusClient{
		http:   http,
		apiKey: apiKey,
		domain: domain,
	}
}

func (hc *HeliusClient) GetTransactions(ctx context.Context, address string) ([]Transaction, error) {
	url := fmt.Sprintf(
		"https://%s/v0/addresses/%s/transactions?api-key=%s",
		hc.domain, address, hc.apiKey,
	)

	res, err := hc.http.Get(
		ctx,
		url,
		nil,
	)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	if !http.IsOK(res) {
		return nil, err
	}

	var body []Transaction
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(b, &body); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	return body, nil
}

func (hc *HeliusClient) CreateWebhook(ctx context.Context, webhookURL *url.URL, address string) error {
	client := http.New(nil, nil, time.Second*10)

	url := fmt.Sprintf(
		"https://%s/v0/webhooks?api-key=%s",
		hc.domain, hc.apiKey,
	)

	res, err := client.Post(
		ctx,
		url,
		nil,
		map[string]any{
			"transactionTypes": []string{"TRANSFER"},
			"accountAddresses": []string{address},
			"webhookURL":       webhookURL.String(),
			"webhookType":      "enhanced",
		},
	)

	if err != nil {
		return err
	}

	defer res.Body.Close()
	if !http.IsOK(res) {
		return err
	}

	return nil
}

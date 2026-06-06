package helius

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"

	"github.com/mkfolder/moxie/pkg/http"
)

type HeliusNet string

const (
	HeliusNetMainnet HeliusNet = "mainnet"
	HeliusNetDevnet  HeliusNet = "devnet"
)

type HeliusClient struct {
	http       *http.Client
	webhookURL *url.URL
	apiKey     string
	domain     string
}

func NewClient(http *http.Client, webhookURL, apiKey string, net HeliusNet) *HeliusClient {
	domain := "api-mainnet.helius-rpc.com"
	if net == HeliusNetDevnet {
		domain = "api-devnet.helius-rpc.com"
	}

	webhook, err := url.Parse(webhookURL)
	if err != nil {
		panic(fmt.Sprintf("failed to parse webhook url: %v", err))
	}

	return &HeliusClient{
		http:       http,
		apiKey:     apiKey,
		domain:     domain,
		webhookURL: webhook,
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

func (hc *HeliusClient) CreateWebhook(ctx context.Context, addresses []string) error {
	url := fmt.Sprintf(
		"https://%s/v0/webhooks?api-key=%s",
		hc.domain, hc.apiKey,
	)

	res, err := hc.http.Post(
		ctx,
		url,
		nil,
		map[string]any{
			"transactionTypes": []string{"TRANSFER"},
			"accountAddresses": addresses,
			"webhookURL":       hc.webhookURL.String(),
			"webhookType":      HeliusWebhookTypeEnhanced,
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

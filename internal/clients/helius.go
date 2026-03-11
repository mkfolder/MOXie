package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/Makefolder/cynero/pkg/http"
)

type NativeTransfer struct {
	FromUserAccount string `json:"fromUserAccount"`
	ToUserAccount   string `json:"toUserAccount"`
	Amount          int    `json:"amount"`
}

type TokenTransfer struct {
	FromUserAccount  string  `json:"fromUserAccount"`
	ToUserAccount    string  `json:"toUserAccount"`
	FromTokenAccount string  `json:"fromTokenAccount"`
	ToTokenAccount   string  `json:"toTokenAccount"`
	TokenAmount      float64 `json:"tokenAmount"`
	Mint             string  `json:"mint"`
}

type RawTokenAmount struct {
	TokenAmount string `json:"tokenAmount"`
	Decimals    int    `json:"decimals"`
}

type TokenBalanceChange struct {
	UserAccount    string         `json:"userAccount"`
	TokenAccount   string         `json:"tokenAccount"`
	Mint           string         `json:"mint"`
	RawTokenAmount RawTokenAmount `json:"rawTokenAmount"`
}

type AccountData struct {
	Account             string               `json:"account"`
	NativeBalanceChange int                  `json:"nativeBalanceChange"`
	TokenBalanceChanges []TokenBalanceChange `json:"tokenBalanceChanges"`
}

type TransactionError struct {
	Error string `json:"error"`
}

type InnerInstruction struct {
	Accounts  []string `json:"accounts"`
	Data      string   `json:"data"`
	ProgramID string   `json:"programId"`
}

type Instruction struct {
	Accounts          []string           `json:"accounts"`
	Data              string             `json:"data"`
	ProgramID         string             `json:"programId"`
	InnerInstructions []InnerInstruction `json:"innerInstructions"`
}

type Transaction struct {
	Type             string            `json:"type"`
	Source           string            `json:"source"`
	Description      string            `json:"description"`
	Fee              int               `json:"fee"`
	FeePayer         string            `json:"fee_payer"`
	Signature        string            `json:"signature"`
	Slot             int               `json:"slot"`
	Timestamp        int               `json:"timestamp"`
	NativeTransfers  []NativeTransfer  `json:"nativeTransfers"`
	TokenTransfers   []TokenTransfer   `json:"tokenTransfers"`
	TransactionError *TransactionError `json:"transactionError"`
	Instructions     []Instruction     `json:"instructions"`
	Events           json.RawMessage   `json:"events"`
}

const (
	// domain = "api-devnet.helius-rpc.com"
	domain = "api-mainnet.helius-rpc.com"
	apiKey = "569dbc7b-d716-4219-aa73-7580f65b3011"
)

func GetTransactions(address string) ([]Transaction, error) {
	client := http.New(nil, nil, time.Second*10)

	url := fmt.Sprintf(
		"https://%s/v0/addresses/%s/transactions?api-key=%s",
		domain, address, apiKey,
	)

	res, err := client.Get(
		context.Background(),
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

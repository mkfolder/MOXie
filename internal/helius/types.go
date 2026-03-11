package helius

import "encoding/json"

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
	FeePayer         string            `json:"feePayer"`
	Signature        string            `json:"signature"`
	Slot             int               `json:"slot"`
	Timestamp        int               `json:"timestamp"`
	NativeTransfers  []NativeTransfer  `json:"nativeTransfers"`
	TokenTransfers   []TokenTransfer   `json:"tokenTransfers"`
	TransactionError *TransactionError `json:"transactionError"`
	Instructions     []Instruction     `json:"instructions"`
	Events           json.RawMessage   `json:"events"`
	AccountData      []AccountData     `json:"accountData"`
}

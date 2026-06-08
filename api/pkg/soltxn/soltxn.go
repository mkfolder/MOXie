package soltxn

import (
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

// Builder is the common interface for all transaction builders.
// Build assembles and returns a transaction using the provided blockhash.
type Builder interface {
	Build(blockhash solana.Hash) (*solana.Transaction, error)
}

// NativeTransferBuilder builds a native SOL transfer transaction,
// optionally with a memo instruction.
type NativeTransferBuilder struct {
	sender    solana.PublicKey
	recipient solana.PublicKey
	rawAmount uint64
	memo      *string
}

// TokenTransferBuilder builds an SPL token transfer transaction.
// It extends NativeTransferBuilder with token-specific fields.
type TokenTransferBuilder struct {
	NativeTransferBuilder

	programID solana.PublicKey
	tokenMint solana.PublicKey
	decimals  uint8
	rpc       *rpc.Client
}

func GetNativeTransferBuilder() *NativeTransferBuilder {
	return &NativeTransferBuilder{}
}

func GetTokenTransferBuilder(rpc *rpc.Client) *TokenTransferBuilder {
	return &TokenTransferBuilder{
		rpc: rpc,
	}
}

func buildMemoInstruction(memo string, sender solana.PublicKey) solana.Instruction {
	return solana.NewInstruction(
		solana.MemoProgramID,
		solana.AccountMetaSlice{
			solana.NewAccountMeta(sender, true, true),
		},
		[]byte(memo),
	)
}

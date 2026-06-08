package soltxn

import (
	"context"
	"time"

	"github.com/gagliardetto/solana-go"
	associatedtokenaccount "github.com/gagliardetto/solana-go/programs/associated-token-account"
	"github.com/gagliardetto/solana-go/programs/token"
)

func (b *TokenTransferBuilder) SetSender(sender solana.PublicKey) *TokenTransferBuilder {
	b.sender = sender
	return b
}

func (b *TokenTransferBuilder) SetRecipient(recipient solana.PublicKey) *TokenTransferBuilder {
	b.recipient = recipient
	return b
}

func (b *TokenTransferBuilder) SetAmount(rawAmount uint64) *TokenTransferBuilder {
	b.rawAmount = rawAmount
	return b
}

func (b *TokenTransferBuilder) SetMemo(memo string) *TokenTransferBuilder {
	b.memo = &memo
	return b
}

func (b *TokenTransferBuilder) SetProgramID(programID solana.PublicKey) *TokenTransferBuilder {
	b.programID = programID
	return b
}

func (b *TokenTransferBuilder) SetTokenMint(tokenMint solana.PublicKey) *TokenTransferBuilder {
	b.tokenMint = tokenMint
	return b
}

func (b *TokenTransferBuilder) SetDecimals(decimals uint8) *TokenTransferBuilder {
	b.decimals = decimals
	return b
}

func (b *TokenTransferBuilder) Build(blockhash solana.Hash) (*solana.Transaction, error) {
	if err := b.validate(); err != nil {
		return nil, err
	}

	instructions, err := b.buildTokenInstructions()
	if err != nil {
		return nil, err
	}

	return solana.NewTransaction(
		instructions,
		blockhash,
		solana.TransactionPayer(b.sender),
	)
}

func (b *TokenTransferBuilder) buildTokenInstructions() ([]solana.Instruction, error) {
	var instructions []solana.Instruction

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	senderATA, _, err := solana.FindAssociatedTokenAddress(b.sender, b.tokenMint)
	if err != nil {
		return nil, err
	}

	recipientATA, _, err := solana.FindAssociatedTokenAddress(b.recipient, b.tokenMint)
	if err != nil {
		return nil, err
	}

	account, err := b.rpc.GetAccountInfo(ctx, recipientATA)
	if err != nil || account == nil {
		createATAIx, err := associatedtokenaccount.NewCreateInstruction(
			b.sender, // payer
			b.recipient,
			b.tokenMint,
		).ValidateAndBuild()
		if err != nil {
			return nil, err
		}
		instructions = append(instructions, createATAIx)
	}

	transferIx, err := token.NewTransferCheckedInstruction(
		b.rawAmount,
		b.decimals,
		senderATA,
		b.tokenMint,
		recipientATA,
		b.sender,
		nil,
	).ValidateAndBuild()
	if err != nil {
		return nil, err
	}

	instructions = append(instructions, transferIx)

	if b.memo != nil {
		instructions = append(instructions, buildMemoInstruction(*b.memo, b.sender))
	}

	return instructions, nil
}

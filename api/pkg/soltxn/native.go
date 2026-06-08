package soltxn

import (
	"errors"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
)

func (b *NativeTransferBuilder) SetSender(sender solana.PublicKey) *NativeTransferBuilder {
	b.sender = sender
	return b
}

func (b *NativeTransferBuilder) SetRecipient(recipient solana.PublicKey) *NativeTransferBuilder {
	b.recipient = recipient
	return b
}

func (b *NativeTransferBuilder) SetAmount(rawAmount uint64) *NativeTransferBuilder {
	b.rawAmount = rawAmount
	return b
}

func (b *NativeTransferBuilder) SetMemo(memo string) *NativeTransferBuilder {
	b.memo = &memo
	return b
}

func (b *NativeTransferBuilder) Build(blockhash solana.Hash) (*solana.Transaction, error) {
	instructions, err := b.buildInstructions()
	if err != nil {
		return nil, err
	}

	return solana.NewTransaction(
		instructions,
		blockhash,
		solana.TransactionPayer(b.sender),
	)
}

func (b *NativeTransferBuilder) validate() error {
	if b.sender.IsZero() {
		return errors.New("soltxn: sender not set")
	}
	if b.recipient.IsZero() {
		return errors.New("soltxn: recipient not set")
	}
	if b.rawAmount == 0 {
		return errors.New("soltxn: amount not set")
	}
	return nil
}

func (b *NativeTransferBuilder) buildInstructions() ([]solana.Instruction, error) {
	if err := b.validate(); err != nil {
		return nil, err
	}

	transferIx, err := system.NewTransferInstruction(
		b.rawAmount,
		b.sender,
		b.recipient,
	).ValidateAndBuild()
	if err != nil {
		return nil, err
	}

	instructions := []solana.Instruction{transferIx}

	if b.memo != nil {
		instructions = append(instructions, buildMemoInstruction(*b.memo, b.sender))
	}

	return instructions, nil
}

package service

import (
	"context"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
)

func (s *Service) buildTransferTxn(
	ctx context.Context,
	sender, recipient solana.PublicKey,
	rawAmount uint64,
	memo *string,
) (*solana.Transaction, error) {
	var instructions []solana.Instruction

	transferIx, err := system.NewTransferInstruction(
		rawAmount,
		sender,
		recipient,
	).ValidateAndBuild()
	if err != nil {
		return nil, err
	}

	instructions = append(instructions, transferIx)
	if memo != nil {
		memoIx := buildMemoInstruction(*memo, sender)
		instructions = append(instructions, memoIx)
	}

	recent, err := s.rpc.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return nil, err
	}

	txn, err := solana.NewTransaction(
		instructions,
		recent.Value.Blockhash,
		solana.TransactionPayer(sender),
	)

	if err != nil {
		return nil, err
	}

	return txn, nil
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

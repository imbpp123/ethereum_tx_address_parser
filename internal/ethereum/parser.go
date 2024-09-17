package ethereum

import (
	"context"
	"fmt"

	"trust_walet/internal/ethereum/data"
)

type (
	AddressService interface {
		AddUnique(address string) bool
	}

	BlockService interface {
		GetCurrentNumber(defaultIfEmpty int) (int, error)
		ProcessNewBlocks(ctx context.Context) error
	}

	TransactionService interface {
		FetchAllByAddress(address string) []data.Transaction
	}

	Parser struct {
		address     AddressService
		block       BlockService
		transaction TransactionService
	}
)

func NewParser(
	address AddressService,
	block BlockService,
	transaction TransactionService,
) *Parser {
	return &Parser{
		transaction: transaction,
		block:       block,
		address:     address,
	}
}

func (p *Parser) GetCurrentBlock() int {
	value, err := p.block.GetCurrentNumber(0)
	if err != nil {
		return 0
	}

	return value
}

func (p *Parser) Subscribe(address string) bool {
	return p.address.AddUnique(address)
}

func (p *Parser) GetTransactions(address string) []data.Transaction {
	return p.transaction.FetchAllByAddress(address)
}

func (p *Parser) MonitorTransactions(ctx context.Context) error {
	err := p.block.ProcessNewBlocks(ctx)
	if err != nil {
		return fmt.Errorf("failed to process new blocks: %w", err)
	}

	return nil
}

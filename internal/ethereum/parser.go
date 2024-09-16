package ethereum

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"sync"
)

type (
	Transaction struct {
		Hash  string
		From  string
		To    string
		Value string
	}

	RpcClient interface {
		GetTransactionByHash(ctx context.Context, hash string) (*RPCTransaction, error)
		GetBlockByNumber(ctx context.Context, number string) (*RPCBlock, error)
	}

	TxStorage interface {
		Save(address string, transaction *Transaction)
		Exists(address, hash string) bool
		FetchAllByAddress(address string) []Transaction
	}

	Parser struct {
		client             RpcClient
		transactionStorage TxStorage

		addresses    []string
		currentBlock int

		muBlock   sync.RWMutex
		muAddress sync.RWMutex
	}
)

func NewParser(
	client RpcClient,
	transactionStorage TxStorage,
) *Parser {
	return &Parser{
		client:             client,
		transactionStorage: transactionStorage,
		currentBlock:       0,
	}
}

func (p *Parser) GetCurrentBlock() int {
	p.muBlock.RLock()
	defer p.muBlock.RUnlock()

	return p.currentBlock
}

func (p *Parser) Subscribe(address string) bool {
	p.muAddress.Lock()
	defer p.muAddress.Unlock()

	if slices.Contains(p.addresses, address) {
		return false
	}

	p.addresses = append(p.addresses, address)

	return true
}

func (p *Parser) GetTransactions(address string) []Transaction {
	transactions := p.transactionStorage.FetchAllByAddress(address)

	return transactions
}

func (p *Parser) MonitorTransactions(ctx context.Context) error {
	latestBlock, err := p.client.GetBlockByNumber(ctx, NumberLatest)
	if err != nil {
		return fmt.Errorf("error fetching latest block for monitoring: %v", err)
	}

	latestNumber, err := p.hexToInt(latestBlock.Number)
	if err != nil {
		return fmt.Errorf("error converting block number to int: %v", err)
	}

	p.muBlock.Lock()
	if p.currentBlock == 0 {
		p.currentBlock = latestNumber
	}
	currentBlock := p.currentBlock
	p.muBlock.Unlock()

	for i := currentBlock; i <= latestNumber; i++ {
		err := p.processBlock(ctx, i)
		if err != nil {
			return fmt.Errorf("error processing block %d: %v", i, err)
		}

		p.muBlock.Lock()
		p.currentBlock = i
		p.muBlock.Unlock()
	}

	return nil
}

// processBlock processes transactions from ethereum blocks
func (p *Parser) processBlock(ctx context.Context, blockNumber int) error {
	newBlock, err := p.client.GetBlockByNumber(ctx, fmt.Sprintf("0x%x", blockNumber))
	if err != nil {
		return fmt.Errorf("error getting block %d for processing: %v", blockNumber, err)
	}

	for _, tx := range newBlock.Transactions {
		for _, a := range []string{tx.From, tx.To} {
			if a == "" {
				// transaction.To field can be empty
				continue
			}

			p.muAddress.RLock()
			canSaveTransaction := slices.Contains(p.addresses, a)
			p.muAddress.RUnlock()

			if canSaveTransaction && !p.transactionStorage.Exists(a, tx.Hash) {
				p.transactionStorage.Save(a, &Transaction{
					Hash:  tx.Hash,
					From:  tx.From,
					To:    tx.To,
					Value: tx.Value,
				})
			}
		}
	}

	return nil
}

// hexToInt convert ethereum hex to int
func (p *Parser) hexToInt(hexStr string) (int, error) {
	number, err := strconv.ParseInt(hexStr[2:], 16, 0)
	if err != nil {
		return 0, fmt.Errorf("error pasring ethereum hex to int: %v", err)
	}

	return int(number), nil
}

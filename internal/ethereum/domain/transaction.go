package domain

import (
	"context"
	"fmt"

	"trust_walet/internal/ethereum/data"
	"trust_walet/internal/ethereum/rpc"

	"github.com/sirupsen/logrus"
)

type (
	AddressServiceInterface interface {
		IsSubscribed(address string) bool
	}

	TransactionStorage interface {
		SaveForAddress(address string, transaction *data.Transaction)
		Exists(address, hash string) bool
		FetchAllByAddress(address string) []data.Transaction
	}

	TransactionRpcClient interface {
		GetBlockByNumber(ctx context.Context, number string) (*rpc.Block, error)
	}

	TransactionService struct {
		client     TransactionRpcClient
		address    AddressServiceInterface
		transation TransactionStorage
	}
)

func NewTransactionService(
	client TransactionRpcClient,
	address AddressServiceInterface,
	transation TransactionStorage,
) *TransactionService {
	return &TransactionService{
		client:     client,
		address:    address,
		transation: transation,
	}
}

func (t *TransactionService) FetchAllByAddress(addr string) []data.Transaction {
	logrus.
		WithFields(logrus.Fields{
			"address": addr,
		}).
		Info("Fetch once transactions for address")

	return t.transation.FetchAllByAddress(addr)
}

func (t *TransactionService) ProcessBlockTransactionsByBlockNumber(ctx context.Context, number int) error {
	block, err := t.client.GetBlockByNumber(ctx, fmt.Sprintf("0x%x", number))
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"block_number": number,
			}).
			WithError(err).
			Error("failed to get block by number")

		return fmt.Errorf("error getting block %d for processing: %w", number, err)
	}

	for _, tx := range block.Transactions {
		txAddresses := []string{tx.From, tx.To}

		for _, a := range txAddresses {
			if t.address.IsSubscribed(a) && !t.transation.Exists(a, tx.Hash) {
				t.transation.SaveForAddress(a, &data.Transaction{
					Hash:  tx.Hash,
					From:  tx.From,
					To:    tx.To,
					Value: tx.Value,
				})

				logrus.
					WithFields(logrus.Fields{
						"block_number":     number,
						"address":          a,
						"transaction_hash": tx.Hash,
					}).
					Info("Block transaction was saved for address")
			}
		}
	}

	logrus.
		WithFields(logrus.Fields{
			"block_number": number,
		}).
		Info("Block transactions were processed")

	return nil
}

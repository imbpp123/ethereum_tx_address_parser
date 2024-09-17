package domain

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"

	"trust_walet/internal/ethereum/rpc"
	"trust_walet/internal/ethereum/storage"

	"github.com/sirupsen/logrus"
)

type (
	TransactionServiceInterface interface {
		ProcessBlockTransactionsByBlockNumber(ctx context.Context, number int) error
	}

	BlockRpcClient interface {
		GetBlockByNumber(ctx context.Context, number string) (*rpc.Block, error)
	}

	BlockStorage interface {
		SetCurrentBlockNumber(value int)
		GetCurrentBlockNumber() (int, error)
	}

	BlockService struct {
		client      BlockRpcClient
		storage     BlockStorage
		transaction TransactionServiceInterface

		mu sync.Mutex
	}
)

func NewBlockService(
	client BlockRpcClient,
	storage BlockStorage,
	transaction TransactionServiceInterface,
) *BlockService {
	return &BlockService{
		client:      client,
		storage:     storage,
		transaction: transaction,
	}
}

func (b *BlockService) GetCurrentNumber(defaultIfEmpty int) (int, error) {
	value, err := b.storage.GetCurrentBlockNumber()
	if err != nil {
		if errors.Is(err, storage.ErrBlockCurrentNotSet) {
			logrus.
				WithError(storage.ErrBlockCurrentNotSet).
				Info("current block number is not set in storage")

			return defaultIfEmpty, nil
		}

		logrus.
			WithError(err).
			Error("failed to get current block number from storage")

		return 0, fmt.Errorf("failed to get current block number: %w", err)
	}

	return value, nil
}

func (b *BlockService) ProcessNewBlocks(ctx context.Context) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	block, err := b.client.GetBlockByNumber(ctx, rpc.NumberLatest)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"block_number": rpc.NumberLatest,
			}).
			WithError(err).
			Error("failed to get block number")

		return fmt.Errorf("error fetching latest block for monitoring: %w", err)
	}

	lastNumber, err := strconv.ParseInt(block.Number[2:], 16, 0)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"block_number": block.Number,
			}).
			WithError(err).
			Error("failed to parse block number")

		return fmt.Errorf("error parsing ethereum hex to int: %w", err)
	}

	currentBlockNumber, err := b.GetCurrentNumber(int(lastNumber))
	if err != nil {
		logrus.
			WithError(err).
			Error("failed to get current block number from storage")

		return fmt.Errorf("failed to get current block number: %w", err)
	}

	for i := currentBlockNumber; i <= int(lastNumber); i++ {
		err := b.transaction.ProcessBlockTransactionsByBlockNumber(ctx, i)
		if err != nil {
			logrus.
				WithFields(logrus.Fields{
					"block_number": i,
				}).
				WithError(err).
				Error("failed to process transaction block")

			return fmt.Errorf("error processing block %d: %w", i, err)
		}

		b.storage.SetCurrentBlockNumber(i)
	}

	logrus.
		WithFields(logrus.Fields{
			"start_block": currentBlockNumber,
			"end_block":   lastNumber,
		}).
		Info("New blocks were processed")

	return nil
}

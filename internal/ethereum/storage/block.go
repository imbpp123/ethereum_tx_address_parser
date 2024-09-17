package storage

import (
	"errors"
	"sync"

	"github.com/sirupsen/logrus"
)

type BlockInMemory struct {
	data *int
	mu   sync.RWMutex
}

var ErrBlockCurrentNotSet = errors.New("current block is not set")

func NewBlockInMemory() *BlockInMemory {
	return &BlockInMemory{}
}

func (b *BlockInMemory) SetCurrentBlockNumber(value int) {
	b.mu.Lock()
	defer b.mu.Unlock()

	logrus.
		WithFields(logrus.Fields{
			"block_number": value,
		}).
		Debug("Current block number was set")

	b.data = &value
}

func (b *BlockInMemory) GetCurrentBlockNumber() (int, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.data == nil {
		return 0, ErrBlockCurrentNotSet
	}

	return *b.data, nil
}

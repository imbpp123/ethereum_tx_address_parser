package storage

import (
	"sync"

	"trust_walet/internal/ethereum/data"

	"github.com/sirupsen/logrus"
)

type TransactionInMemory struct {
	data map[string][]data.Transaction

	mu sync.RWMutex
}

func NewTransactionInMemory() *TransactionInMemory {
	return &TransactionInMemory{
		data: make(map[string][]data.Transaction),
	}
}

func (t *TransactionInMemory) SaveForAddress(address string, transaction *data.Transaction) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.data[address] = append(t.data[address], *transaction)
}

func (t *TransactionInMemory) Exists(address, hash string) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()

	for _, tx := range t.data[address] {
		if tx.Hash == hash {
			return true
		}
	}

	logrus.
		WithFields(logrus.Fields{
			"address": address,
			"hash":    hash,
		}).
		Debug("Address and hash were not found in storage")

	return false
}

func (t *TransactionInMemory) FetchAllByAddress(address string) []data.Transaction {
	t.mu.Lock()
	defer t.mu.Unlock()

	transactions := t.data[address]
	t.data[address] = nil

	return transactions
}

package ethereum

import "sync"

type InMemory struct {
	data map[string][]Transaction

	mu sync.RWMutex
}

func NewInMemory() *InMemory {
	return &InMemory{
		data: make(map[string][]Transaction),
	}
}

func (t *InMemory) Save(address string, transaction *Transaction) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.data[address] = append(t.data[address], *transaction)
}

func (t *InMemory) Exists(address, hash string) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()

	for _, tx := range t.data[address] {
		if tx.Hash == hash {
			return true
		}
	}

	return false
}

func (t *InMemory) FetchAllByAddress(address string) []Transaction {
	t.mu.Lock()
	defer t.mu.Unlock()

	transactions := t.data[address]
	t.data[address] = nil

	return transactions
}

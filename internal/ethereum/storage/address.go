package storage

import (
	"slices"
	"sync"
)

type AddressInMemory struct {
	data []string
	mu   sync.RWMutex
}

func NewAddressInMemory() *AddressInMemory {
	return &AddressInMemory{}
}

func (a *AddressInMemory) Exists(address string) bool {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return slices.Contains(a.data, address)
}

func (a *AddressInMemory) Add(address string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.data = append(a.data, address)
}

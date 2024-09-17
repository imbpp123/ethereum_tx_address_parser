package storage_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"trust_walet/internal/ethereum/data"
	"trust_walet/internal/ethereum/storage"
)

func TestInMemoryFetchAllByAddressSaveAndFetch(t *testing.T) {
	// arrange
	storage := storage.NewTransactionInMemory()
	storage.SaveForAddress("addr1", &data.Transaction{
		Hash:  "hash",
		From:  "addr1",
		To:    "addr2",
		Value: "value",
	})

	// act
	tx := storage.FetchAllByAddress("addr1")

	// assert
	assert.NotEmpty(t, tx)
	assert.Len(t, tx, 1)
}

func TestInMemoryFetchAllByAddressSaveAndFetchTwice(t *testing.T) {
	// arrange
	storage := storage.NewTransactionInMemory()
	storage.SaveForAddress("addr1", &data.Transaction{
		Hash:  "hash",
		From:  "addr1",
		To:    "addr2",
		Value: "value",
	})

	// act
	tx1 := storage.FetchAllByAddress("addr1")

	// assert
	assert.NotEmpty(t, tx1)
	assert.Len(t, tx1, 1)

	// act
	tx2 := storage.FetchAllByAddress("addr1")

	// assert
	assert.Empty(t, tx2)
}

func TestInMemoryFetchAllByAddressAddressNotExist(t *testing.T) {
	// arrange
	storage := storage.NewTransactionInMemory()

	// act
	tx := storage.FetchAllByAddress("addr1")

	// assert
	assert.Empty(t, tx)
}

func TestInMemoryExistsNoAddress(t *testing.T) {
	// arrange
	storage := storage.NewTransactionInMemory()

	// act
	result := storage.Exists("addr1", "hash")

	// assert
	assert.False(t, result)
}

func TestInMemoryExistsNoHash(t *testing.T) {
	// arrange
	storage := storage.NewTransactionInMemory()
	storage.SaveForAddress("addr1", &data.Transaction{
		Hash:  "another hash",
		From:  "addr1",
		To:    "addr2",
		Value: "value",
	})

	// act
	result := storage.Exists("addr1", "hash")

	// assert
	assert.False(t, result)
}

func TestInMemoryExists(t *testing.T) {
	// arrange
	storage := storage.NewTransactionInMemory()
	storage.SaveForAddress("addr1", &data.Transaction{
		Hash:  "hash",
		From:  "addr1",
		To:    "addr2",
		Value: "value",
	})

	// act
	result := storage.Exists("addr1", "hash")

	// assert
	assert.True(t, result)
}

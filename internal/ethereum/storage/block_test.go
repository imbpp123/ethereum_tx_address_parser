package storage_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"trust_walet/internal/ethereum/storage"
)

func TestBlockInMemoryGetCurrentBlockNumberCurrentEmpty(t *testing.T) {
	// arrange
	data := storage.NewBlockInMemory()

	// act
	result := data.GetCurrentBlockNumber()

	// assert
	assert.Equal(t, result, 0)
}

func TestBlockInMemorSetGetCurrentBlock(t *testing.T) {
	// arrange
	data := storage.NewBlockInMemory()
	data.SetCurrentBlockNumber(10)

	// act
	result := data.GetCurrentBlockNumber()

	// assert
	assert.Equal(t, result, 10)
}

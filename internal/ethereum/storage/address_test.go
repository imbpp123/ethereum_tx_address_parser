package storage_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"trust_walet/internal/ethereum/storage"
)

func TestAddressExists(t *testing.T) {
	// arrange
	data := storage.NewAddressInMemory()
	data.Add("any")

	// act
	result := data.Exists("any")

	// assert
	assert.True(t, result)
}

func TestAddressExistsNotExist(t *testing.T) {
	// arrange
	data := storage.NewAddressInMemory()

	// act
	result := data.Exists("any")

	// assert
	assert.False(t, result)
}

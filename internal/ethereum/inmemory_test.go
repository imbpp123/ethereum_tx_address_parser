package ethereum_test

import (
	"testing"

	"trust_walet/internal/ethereum"
)

func TestInMemoryFetchAllByAddress(t *testing.T) {
	storage := ethereum.NewInMemory()

	storage.Save("addr1", &ethereum.Transaction{
		Hash:  "hash",
		From:  "addr1",
		To:    "addr2",
		Value: "value",
	})

	tx1 := storage.FetchAllByAddress("addr1")
	if len(tx1) != 1 {
		t.Errorf("\nwrong transaction amount for %s, should be 1, but %d instead", "addr1", len(tx1))
	}

	tx2 := storage.FetchAllByAddress("addr1")
	if len(tx2) != 0 {
		t.Errorf("\nwrong transaction amount for %s, should be 0, but %d instead", "addr1", len(tx2))
	}
}

func TestInMemoryExists(t *testing.T) {
	storage := ethereum.NewInMemory()

	storage.Save("addr1", &ethereum.Transaction{
		Hash:  "hash",
		From:  "addr1",
		To:    "addr2",
		Value: "value",
	})

	if !storage.Exists("addr1", "hash") {
		t.Errorf("\naddress %s does not have any transaction", "addr1")
	}

	if storage.Exists("addr2", "hash") {
		t.Errorf("\naddress %s does not have any transaction", "addr2")
	}
}

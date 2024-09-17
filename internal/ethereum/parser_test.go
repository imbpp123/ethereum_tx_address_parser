package ethereum_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"trust_walet/internal/ethereum"
	"trust_walet/internal/ethereum/domain"
	"trust_walet/internal/ethereum/rpc"
	"trust_walet/internal/ethereum/storage"
)

func TestParserMonitorTransactions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := createGetBlockResponse(
			createBlockResponse(1,
				[]map[string]interface{}{
					createTransactionResponse("0x1", "addr1", "addr2"),
					createTransactionResponse("0x2", "addr2", "addr3"),
				},
			),
		)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		}
	}))
	defer server.Close()

	// arrange
	ctx := context.Background()

	client := rpc.NewHttp(&http.Client{}, server.URL)

	addressService := domain.NewAddressService(
		storage.NewAddressInMemory(),
	)
	transactionService := domain.NewTransactionService(
		client,
		addressService,
		storage.NewTransactionInMemory(),
	)
	blockService := domain.NewBlockService(
		client,
		storage.NewBlockInMemory(),
		transactionService,
	)

	parser := ethereum.NewParser(
		addressService,
		blockService,
		transactionService,
	)

	// act
	parser.Subscribe("addr2")
	parser.MonitorTransactions(ctx)

	// assert
	assert.Equal(t, parser.GetCurrentBlock(), 1)

	tx := parser.GetTransactions("addr2")
	assert.Len(t, tx, 2)
}

func createTransactionResponse(hash, from, to string) map[string]interface{} {
	transaction := make(map[string]interface{})
	transaction["hash"] = hash
	transaction["from"] = from
	transaction["to"] = to
	transaction["value"] = "0xf3"

	return transaction
}

func createBlockResponse(number int, transactions []map[string]interface{}) map[string]interface{} {
	block := make(map[string]interface{})
	block["number"] = fmt.Sprintf("0x%x", number)
	block["transactions"] = transactions

	return block
}

func createGetBlockResponse(block map[string]interface{}) map[string]interface{} {
	request := make(map[string]interface{})
	request["jsonrpc"] = "2.0"
	request["result"] = block

	return request
}

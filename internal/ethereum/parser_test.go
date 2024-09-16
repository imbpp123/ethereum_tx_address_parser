package ethereum_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"trust_walet/internal/ethereum"
)

func TestParserMonitorTransactions(t *testing.T) {
	responses := map[float64]map[string]interface{}{
		1: createGetBlockByNumberResponse(
			1,
			createBlockResponse(1,
				[]map[string]interface{}{
					createTransactionResponse("0x1", "addr1", "addr2"),
					createTransactionResponse("0x2", "addr2", "addr3"),
				},
			),
		),
		2: createGetBlockByNumberResponse(
			1,
			createBlockResponse(1,
				[]map[string]interface{}{
					createTransactionResponse("0x1", "addr1", "addr2"),
					createTransactionResponse("0x2", "addr2", "addr3"),
				},
			),
		),
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if http.MethodPost != r.Method {
			t.Errorf("\nrequest method shoud be POST, %s is used instead", r.Method)
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("\nfailed to read request body %v", err)
		}
		fmt.Printf("\nreceived request: %s", string(body))

		rpcRequest := make(map[string]interface{})
		if err = json.Unmarshal(body, &rpcRequest); err != nil {
			t.Errorf("\nfailed to decode request: %v", err)
		}

		if rpcRequest["method"] != "eth_getBlockByNumber" {
			t.Errorf("\nRPC method is not eth_getBlockByNumber, %s is set instead", rpcRequest["method"])
		}

		response, ok := responses[rpcRequest["id"].(float64)]
		if !ok {
			t.Errorf("\nresponse not found for request id %f", rpcRequest["id"].(float64))
		}

		responseJSON, _ := json.Marshal(response)
		fmt.Printf("\nsending response: %s\n", string(responseJSON))

		w.Header().Set("Content-Type", "application/json")
		w.Write(responseJSON)
	}))
	defer server.Close()

	ctx := context.Background()

	client := ethereum.NewRpc(&http.Client{}, server.URL)
	storage := ethereum.NewInMemory()
	parser := ethereum.NewParser(client, storage)

	parser.Subscribe("addr1")
	parser.Subscribe("addr2")
	parser.MonitorTransactions(ctx)

	if parser.GetCurrentBlock() != 1 {
		t.Errorf("\nfailed to get current block number, it should be %d, but returned %d", 1, parser.GetCurrentBlock())
	}

	addr1Tx := parser.GetTransactions("addr1")
	if len(addr1Tx) != 1 {
		t.Errorf("\n%s tx array should be %d, %d elements in array instead", "addr1", 1, len(addr1Tx))
	}

	addr2Tx := parser.GetTransactions("addr2")
	if len(addr2Tx) != 2 {
		t.Errorf("\n%s tx array should be %d, %d elements in array instead", "addr2", 2, len(addr2Tx))
	}

	addr3Tx := parser.GetTransactions("addr3")
	if len(addr3Tx) != 0 {
		t.Errorf("\n%s tx array should be %d, %d elements in array instead", "addr1", 0, len(addr3Tx))
	}
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

func createGetBlockByNumberResponse(id int, block map[string]interface{}) map[string]interface{} {
	request := make(map[string]interface{})
	request["jsonrpc"] = "2.0"
	request["result"] = block
	request["id"] = id

	return request
}

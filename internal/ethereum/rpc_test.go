package ethereum_test

import (
	"context"
	"net/http"
	"testing"

	"trust_walet/internal/ethereum"
)

func TestRpcGetBlockByNumber(t *testing.T) {
	blockNumber := "0x13cb921"

	ctx := context.Background()
	client := ethereum.NewRpc(&http.Client{}, ethereum.RPCUrl)

	block, err := client.GetBlockByNumber(ctx, blockNumber)

	if err != nil {
		t.Errorf("\nerror while fetching block: %v", err)
	}

	if block == nil {
		t.Errorf("\nblock is nil")
	}

	if block != nil && block.Number != blockNumber {
		t.Errorf("\block number is not equal to %s", blockNumber)
	}

	if len(block.Transactions) == 0 {
		t.Errorf("\block %s should contain transactions", blockNumber)
	}
}

package ethereum

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
)

const (
	RPCUrl     = "https://ethereum-rpc.publicnode.com"
	rpcVersion = "2.0"

	methodGetBlockByNumber     = "eth_getBlockByNumber"
	methodGetTransactionByHash = "eth_getTransactionByHash"

	NumberLatest = "latest"
)

type (
	RPCBlock struct {
		Number       string        `json:"number"`
		Transactions []Transaction `json:"transactions"`
	}

	RPCTransaction struct {
		Hash  string `json:"hash"`
		From  string `json:"from"`
		To    string `json:"to,omitempty"`
		Value string `json:"value"`
	}

	rpcRequest struct {
		JSONRPC string        `json:"jsonrpc"`
		Method  string        `json:"method"`
		Params  []interface{} `json:"params"`
		ID      uint64        `json:"id"`
	}

	rpcResponse struct {
		JSONRPC string          `json:"jsonrpc"`
		ID      int             `json:"id"`
		Result  json.RawMessage `json:"result"`
		Error   *rpcError       `json:"error,omitempty"`
	}

	rpcError struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	Rpc struct {
		client *http.Client

		idCounter uint64
		url       string
	}
)

func NewRpc(
	client *http.Client,
	url string,
) *Rpc {
	return &Rpc{
		client:    client,
		idCounter: 0,
		url:       url,
	}
}

func (r *Rpc) GetTransactionByHash(ctx context.Context, hash string) (*RPCTransaction, error) {
	reqBody := rpcRequest{
		JSONRPC: rpcVersion,
		Method:  methodGetTransactionByHash,
		Params:  []interface{}{hash},
		ID:      atomic.AddUint64(&r.idCounter, 1),
	}

	resp, err := r.sendRequest(ctx, &reqBody)
	if err != nil {
		return nil, fmt.Errorf("error during %s request: %v", methodGetTransactionByHash, err)
	}

	var tx RPCTransaction
	if err := json.Unmarshal(resp.Result, &tx); err != nil {
		return nil, fmt.Errorf("error unmarshaling block: %v", err)
	}

	return &tx, nil
}

func (r *Rpc) GetBlockByNumber(ctx context.Context, number string) (*RPCBlock, error) {
	reqBody := rpcRequest{
		JSONRPC: rpcVersion,
		Method:  methodGetBlockByNumber,
		Params:  []interface{}{number, true},
		ID:      atomic.AddUint64(&r.idCounter, 1),
	}

	resp, err := r.sendRequest(ctx, &reqBody)
	if err != nil {
		return nil, fmt.Errorf("error during %s request: %v", methodGetTransactionByHash, err)
	}

	var block RPCBlock
	if err := json.Unmarshal(resp.Result, &block); err != nil {
		return nil, fmt.Errorf("error unmarshaling block: %v", err)
	}

	return &block, nil
}

func (r *Rpc) sendRequest(ctx context.Context, reqBody *rpcRequest) (*rpcResponse, error) {
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error encoding request: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", r.url, bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	var rpcResp rpcResponse
	if err := json.Unmarshal(body, &rpcResp); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	if rpcResp.Error != nil {
		return nil, fmt.Errorf("RPC error: code %d, message %s", rpcResp.Error.Code, rpcResp.Error.Message)
	}

	return &rpcResp, nil
}

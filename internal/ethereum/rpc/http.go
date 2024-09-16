package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	EthereumUrl = "https://ethereum-rpc.publicnode.com"

	NumberLatest = "latest"

	rpcVersion = "2.0"

	methodGetBlockByNumber = "eth_getBlockByNumber"
)

type (
	Block struct {
		Number       string        `json:"number"`
		Transactions []Transaction `json:"transactions"`
	}

	Transaction struct {
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

	Http struct {
		client *http.Client

		idCounter uint64
		url       string
	}
)

var (
	ErrEthereumServerUnavailable = errors.New("ethereum server is unavailable")
	ErrRPCResponseError          = errors.New("rpc error is returned")
)

func NewHttp(
	client *http.Client,
	url string,
) *Http {
	return &Http{
		client:    client,
		idCounter: 0,
		url:       url,
	}
}

func (r *Http) GetBlockByNumber(ctx context.Context, number string) (*Block, error) {
	reqBody := rpcRequest{
		JSONRPC: rpcVersion,
		Method:  methodGetBlockByNumber,
		Params:  []interface{}{number, true},
		ID:      atomic.AddUint64(&r.idCounter, 1),
	}

	resp, err := r.sendRequest(ctx, &reqBody)
	if err != nil {
		return nil, fmt.Errorf("error during %s request: %w", methodGetBlockByNumber, err)
	}

	var block Block
	if err := json.Unmarshal(resp.Result, &block); err != nil {
		return nil, fmt.Errorf("error unmarshaling block: %w", err)
	}

	return &block, nil
}

func (r *Http) sendRequest(ctx context.Context, reqBody *rpcRequest) (*rpcResponse, error) {
	log := logrus.WithFields(logrus.Fields{
		"method":     reqBody.Method,
		"params":     reqBody.Params,
		"request_id": reqBody.ID,
	})

	start := time.Now()

	log.Debug("Sending request to Ethereum node")

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error encoding request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", r.url, bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, ErrEthereumServerUnavailable
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var rpcResp rpcResponse
	if err := json.Unmarshal(body, &rpcResp); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	if rpcResp.Error != nil {
		log.
			WithFields(logrus.Fields{
				"error_code":    rpcResp.Error.Code,
				"error_message": rpcResp.Error.Message,
			}).
			Error("RPC error happened")

		return nil, ErrRPCResponseError
	}

	log.
		WithFields(logrus.Fields{
			"duration": time.Since(start),
		}).
		Debug("Request is done")

	return &rpcResp, nil
}

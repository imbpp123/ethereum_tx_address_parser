package rpc_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"trust_walet/internal/ethereum/rpc"
)

var (
	getBlockByNumberResponseBody string = `{
		"jsonrpc": "2.0",
		"result": {
			"baseFeePerGas": "0x18d08d45f",
			"blobGasUsed": "0x0",
			"difficulty": "0x0",
			"excessBlobGas": "0xa0000",
			"extraData": "0x6265617665726275696c642e6f7267",
			"gasLimit": "0x1c9c380",
			"gasUsed": "0xeae28d",
			"hash": "0x27bfd8fa4bde92f1d0ef2461d3d3835fd6b8e2949c191d31fb6282e6740ce4bc",
			"logsBloom": "0x332502a2c10c9304f0544a7e8070aba11307cca1cdc16187060d7d503fd5b90f57d75708c0bbe08e6a593f6975c679340a90f1fa9fe3ea72d18caee5c1fe38207d02fe5b17e13e39bc0776a9e0eff224e520f11953dd5d3304789ec58cee648910df25d0a7720478110ed99c81526f11efaf2a65223c1413925da79500dd3a4169b0c61c108a948a38c8b1422a645352f04e06c1a3c81079a421a1570732399ebb9d1184181175aa26c1f1c6de9e25582097a99775681d0a4f74669b02ca81d503f111cff7968d538121083d4693981454ffed3e070f66b041054716e3ce685e2af1f619aa11a4744b7417c8bd290ed6093993c3598830f3f5bd6931b4d5368b",
			"miner": "0x95222290dd7278aa3ddd389cc1e1d165cc4bafe5",
			"mixHash": "0xaa7d3b83601e38a24047d77276803432a6b89684f15d5ab66253541a667c10fa",
			"nonce": "0x0000000000000000",
			"number": "0x13cdb01",
			"parentBeaconBlockRoot": "0x808f741101d63b49701ff98ea4dce03989172f62f51c268613c92e083e80cd13",
			"parentHash": "0x116d163f9bc5608d77acd947390c1d60238cbdf8c36504289dcb2b68767172d2",
			"receiptsRoot": "0x267c2b2d89e10c5049d6648846835057c1981528df32f46c68f6fdc310113604",
			"sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
			"size": "0x292a1",
			"stateRoot": "0x61a3ac9868a70b518cdbf25d492d23ff8e1b8a974360b45a2a89b5d3b84a1c4a",
			"timestamp": "0x66e88f3b",
			"totalDifficulty": "0xc70d815d562d3cfa955",
			"transactions":[
				{
					"accessList":[],
					"blockHash":"0xd31a13092d2c95917da5547d0df4ac81327dacba9adbe4b24a32d4f1f5d336c6",
					"blockNumber":"0x13cdb47",
					"chainId":"0x1",
					"from":"0x3b7a1379b3aeb36420db4533035aa04e9a3ce615",
					"gas":"0x6a021",
					"gasPrice":"0x18b0ee9f5",
					"hash":"0x69307e8052d183489c1197687a586b42572a88fed9c76bbdf784f3f068d14284",
					"input":"0x0162e2d000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000001600000000000000000000000007a250d5630b4cf539739df2c5dacb4c659f2488d000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000066e89283000000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000051b660cdd580000000000000000000000000000000000000000000000000000000000000000002000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc200000000000000000000000026e2abc5edf473a2d5cedef24f47e6f2e29c3c3e",
					"maxFeePerGas":"0x2b1a06a7e",
					"maxPriorityFeePerGas":"0x3747177",
					"nonce":"0x71c",
					"r":"0x976270a2a69c16b6a1cf8b8aaf4c47023fe817259cf5284ee718955bf2acc40",
					"s":"0x5a7bb8611c0612fdbbe5522c78b799d8ca01edf45de8654ed3673e310fbf13a0",
					"to":"0x3328f7f4a1d1c57c35df56bbf0c9dcafca309c49",
					"transactionIndex":"0x0",
					"type":"0x2",
					"v":"0x1",
					"value":"0x1b4fbd92b5f8000",
					"yParity":"0x1"
				},
				{
					"accessList":[],
					"blockHash":"0xd31a13092d2c95917da5547d0df4ac81327dacba9adbe4b24a32d4f1f5d336c6",
					"blockNumber":
					"0x13cdb47",
					"chainId":"0x1",
					"from":"0xc6d065e48fc3848328d3b8eefbe38f0068542d9c",
					"gas":"0x6a021",
					"gasPrice":"0x1885c0398",
					"hash":"0x0fed43938113e0151d728d3ad9fa039d5216222d289451e011503b37c4b46142",
					"input":"0x0162e2d000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000001600000000000000000000000007a250d5630b4cf539739df2c5dacb4c659f2488d000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000066e89283000000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000051b660cdd580000000000000000000000000000000000000000000000000000000000000000002000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc200000000000000000000000026e2abc5edf473a2d5cedef24f47e6f2e29c3c3e",
					"maxFeePerGas":"0x2b1a06a7e",
					"maxPriorityFeePerGas":"0xc18b1a",
					"nonce":"0x3b5",
					"r":"0x47e51efd85f917c315a34d63e8945abe6e922b1ece2bcde8c4cd81926e25a68f",
					"s":"0x7f11ba3cb7322bd7a4e8c5d759eb77c66671aa74515e8c9c00a932907e7a44f1",
					"to":"0x3328f7f4a1d1c57c35df56bbf0c9dcafca309c49",
					"transactionIndex":"0x1",
					"type":"0x2",
					"v":"0x1",
					"value":"0x1b4fbd92b5f8000",
					"yParity":"0x1"
				}
			],
			"transactionsRoot": "0x71280fc79c51ae750d48693e3b27cf23586c88cb2e705d98b0ab892974d73564",
			"uncles": [],
			"withdrawals": [
				{
					"index": "0x38ec2fb",
					"validatorIndex": "0x16c2ba",
					"address": "0x2f777e9f26aa138ed21c404079e80656b448c711",
					"amount": "0x1247eaf"
				}
			],
			"withdrawalsRoot": "0xdd6143733d321c4e08c6e58d7e27247bc0b9ad4454ba956c2408fd27bbf61971"
		}
	}`
	rpcErrorMessage string = `{
		"jsonrpc": "2.0",
		"error": {
			"code": 1234,
			"message": "error happened"
		}
	}`
)

func TestRpcGetBlockByNumber(t *testing.T) {
	logrus.SetOutput(io.Discard)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		body, err := io.ReadAll(r.Body)
		if !assert.NoError(t, err) {
			return
		}

		rpcRequest := make(map[string]interface{})
		if !assert.NoError(t, json.Unmarshal(body, &rpcRequest)) {
			return
		}
		assert.Equal(t, "eth_getBlockByNumber", rpcRequest["method"])

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(getBlockByNumberResponseBody))
	}))
	defer server.Close()

	// arrange
	blockNumber := "0x13cdb01"
	ctx := context.Background()
	client := rpc.NewHttp(&http.Client{}, server.URL)

	// act
	block, err := client.GetBlockByNumber(ctx, blockNumber)

	// assert
	if assert.NoError(t, err) {
		if assert.NotNil(t, block) {
			assert.Equal(t, blockNumber, block.Number)
			assert.Len(t, block.Transactions, 2)
		}
	}
}

func TestRpcGetBlockByNumberServerUnavailable(t *testing.T) {
	logrus.SetOutput(io.Discard)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	// arrange
	ctx := context.Background()
	client := rpc.NewHttp(&http.Client{}, server.URL)

	// act
	block, err := client.GetBlockByNumber(ctx, "any")

	// assert
	assert.ErrorIs(t, err, rpc.ErrEthereumServerUnavailable)
	assert.Nil(t, block)
}

func TestRpcGetBlockByNumberMailformedJson(t *testing.T) {
	logrus.SetOutput(io.Discard)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{..."))
	}))
	defer server.Close()

	// arrange
	ctx := context.Background()
	client := rpc.NewHttp(&http.Client{}, server.URL)

	// act
	block, err := client.GetBlockByNumber(ctx, "any")

	// assert
	assert.Error(t, err)
	assert.Nil(t, block)
}

func TestRpcGetBlockByNumberRpcError(t *testing.T) {
	logrus.SetOutput(io.Discard)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(rpcErrorMessage))
	}))
	defer server.Close()

	// arrange
	ctx := context.Background()
	client := rpc.NewHttp(&http.Client{}, server.URL)

	// act
	block, err := client.GetBlockByNumber(ctx, "any")

	// assert
	assert.ErrorIs(t, err, rpc.ErrRPCResponseError)
	assert.Nil(t, block)
}

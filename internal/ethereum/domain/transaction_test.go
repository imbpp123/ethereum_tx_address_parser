package domain_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"trust_walet/internal/ethereum/domain"
	mockDomain "trust_walet/internal/ethereum/domain/mock"
	"trust_walet/internal/ethereum/rpc"
)

type unitTransactionService struct {
	mockClient             *mockDomain.MockTransactionRpcClient
	mockAddressService     *mockDomain.MockAddressServiceInterface
	mockTransactionStorage *mockDomain.MockTransactionStorage
	transactionService     *domain.TransactionService
}

func newUnitTransactionService(ctrl *gomock.Controller) *unitTransactionService {
	unit := unitTransactionService{
		mockClient:             mockDomain.NewMockTransactionRpcClient(ctrl),
		mockAddressService:     mockDomain.NewMockAddressServiceInterface(ctrl),
		mockTransactionStorage: mockDomain.NewMockTransactionStorage(ctrl),
	}

	unit.transactionService = domain.NewTransactionService(
		unit.mockClient,
		unit.mockAddressService,
		unit.mockTransactionStorage,
	)

	return &unit
}

func TestTransactionServiceProcessBlockTransactionsByBlockNumber(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// arrange
	tc := newUnitTransactionService(ctrl)
	block := rpc.Block{
		Number: "0x1",
		Transactions: []rpc.Transaction{
			{
				Hash: "hash",
				From: "addr1",
				To:   "addr2",
			},
		},
	}

	// assert
	tc.mockClient.EXPECT().GetBlockByNumber(gomock.Any(), gomock.Eq(block.Number)).Return(&block, nil)

	tc.mockAddressService.EXPECT().IsSubscribed(gomock.Eq("addr1")).Return(true)
	tc.mockTransactionStorage.EXPECT().Exists(gomock.Eq("addr1"), gomock.Eq("hash")).Return(false)
	tc.mockTransactionStorage.EXPECT().SaveForAddress(gomock.Eq("addr1"), gomock.Any())

	tc.mockAddressService.EXPECT().IsSubscribed(gomock.Eq("addr2")).Return(true)
	tc.mockTransactionStorage.EXPECT().Exists(gomock.Eq("addr2"), gomock.Eq("hash")).Return(false)
	tc.mockTransactionStorage.EXPECT().SaveForAddress(gomock.Eq("addr2"), gomock.Any())

	// act
	err := tc.transactionService.ProcessBlockTransactionsByBlockNumber(context.Background(), 1)

	// assert
	assert.NoError(t, err)
}

func TestTransactionServiceProcessBlockTransactionsByBlockNumberAddrNotSubscribed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// arrange
	tc := newUnitTransactionService(ctrl)
	block := rpc.Block{
		Number: "0x1",
		Transactions: []rpc.Transaction{
			{
				Hash: "hash",
				From: "addr1",
				To:   "addr2",
			},
		},
	}

	// assert
	tc.mockClient.EXPECT().GetBlockByNumber(gomock.Any(), gomock.Eq(block.Number)).Return(&block, nil)

	tc.mockAddressService.EXPECT().IsSubscribed(gomock.Eq("addr1")).Return(true)
	tc.mockTransactionStorage.EXPECT().Exists(gomock.Eq("addr1"), gomock.Eq("hash")).Return(false)
	tc.mockTransactionStorage.EXPECT().SaveForAddress(gomock.Eq("addr1"), gomock.Any())

	tc.mockAddressService.EXPECT().IsSubscribed(gomock.Eq("addr2")).Return(false)
	tc.mockTransactionStorage.EXPECT().Exists(gomock.Eq("addr2"), gomock.Eq("hash")).Times(0)
	tc.mockTransactionStorage.EXPECT().SaveForAddress(gomock.Eq("addr2"), gomock.Any()).Times(0)

	// act
	err := tc.transactionService.ProcessBlockTransactionsByBlockNumber(context.Background(), 1)

	// assert
	assert.NoError(t, err)
}

func TestTransactionServiceProcessBlockTransactionsByBlockNumberTransactionInStorage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// arrange
	tc := newUnitTransactionService(ctrl)
	block := rpc.Block{
		Number: "0x1",
		Transactions: []rpc.Transaction{
			{
				Hash: "hash",
				From: "addr1",
				To:   "addr2",
			},
		},
	}

	// assert
	tc.mockClient.EXPECT().GetBlockByNumber(gomock.Any(), gomock.Eq(block.Number)).Return(&block, nil)

	tc.mockAddressService.EXPECT().IsSubscribed(gomock.Eq("addr1")).Return(true)
	tc.mockTransactionStorage.EXPECT().Exists(gomock.Eq("addr1"), gomock.Eq("hash")).Return(false)
	tc.mockTransactionStorage.EXPECT().SaveForAddress(gomock.Eq("addr1"), gomock.Any())

	tc.mockAddressService.EXPECT().IsSubscribed(gomock.Eq("addr2")).Return(true)
	tc.mockTransactionStorage.EXPECT().Exists(gomock.Eq("addr2"), gomock.Eq("hash")).Return(true)
	tc.mockTransactionStorage.EXPECT().SaveForAddress(gomock.Eq("addr2"), gomock.Any()).Times(0)

	// act
	err := tc.transactionService.ProcessBlockTransactionsByBlockNumber(context.Background(), 1)

	// assert
	assert.NoError(t, err)
}

func TestTransactionServiceProcessBlockTransactionsByBlockNumberClientFailed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// arrange
	tc := newUnitTransactionService(ctrl)

	// assert
	tc.mockClient.EXPECT().GetBlockByNumber(gomock.Any(), gomock.Any()).Return(nil, errors.New("failed to make request"))

	tc.mockTransactionStorage.EXPECT().SaveForAddress(gomock.Any(), gomock.Any()).Times(0)

	// act
	err := tc.transactionService.ProcessBlockTransactionsByBlockNumber(context.Background(), 1)

	// assert
	assert.Error(t, err)
}

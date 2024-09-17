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
	"trust_walet/internal/ethereum/storage"
)

type unitBlockService struct {
	mockClient             *mockDomain.MockBlockRpcClient
	mockBlockStorage       *mockDomain.MockBlockStorage
	mockTransactionService *mockDomain.MockTransactionServiceInterface
	blockService           *domain.BlockService
}

func newUnitBlockService(ctrl *gomock.Controller) *unitBlockService {
	unit := unitBlockService{
		mockClient:             mockDomain.NewMockBlockRpcClient(ctrl),
		mockBlockStorage:       mockDomain.NewMockBlockStorage(ctrl),
		mockTransactionService: mockDomain.NewMockTransactionServiceInterface(ctrl),
	}
	unit.blockService = domain.NewBlockService(
		unit.mockClient,
		unit.mockBlockStorage,
		unit.mockTransactionService,
	)

	return &unit
}

func TestBlockServiceGetCurrentNumber(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// arrange
	tc := newUnitBlockService(ctrl)

	tc.mockBlockStorage.EXPECT().GetCurrentBlockNumber().Return(10, nil)

	// act
	result, err := tc.blockService.GetCurrentNumber(20)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, result, 10)
}

func TestBlockServiceGetCurrentNumberNotSet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// arrange
	tc := newUnitBlockService(ctrl)

	tc.mockBlockStorage.EXPECT().GetCurrentBlockNumber().Return(0, storage.ErrBlockCurrentNotSet)

	// act
	result, err := tc.blockService.GetCurrentNumber(20)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, result, 20)
}

func TestBlockServiceGetCurrentNumberStorageError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// arrange
	tc := newUnitBlockService(ctrl)

	tc.mockBlockStorage.EXPECT().GetCurrentBlockNumber().Return(0, errors.New("any error"))

	// act
	result, err := tc.blockService.GetCurrentNumber(20)

	// assert
	assert.Error(t, err)
	assert.Equal(t, result, 0)
}

func TestBlockServiceProcessNewBlocks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// arrange
	tc := newUnitBlockService(ctrl)

	block := rpc.Block{
		Number: "0x2",
	}

	// assert
	tc.mockClient.EXPECT().GetBlockByNumber(gomock.Any(), gomock.Eq("latest")).Return(&block, nil)
	tc.mockBlockStorage.EXPECT().GetCurrentBlockNumber().Return(1, nil)

	processed1 := tc.mockTransactionService.EXPECT().ProcessBlockTransactionsByBlockNumber(gomock.Any(), gomock.Eq(1)).Return(nil)
	setCurrent1 := tc.mockBlockStorage.EXPECT().SetCurrentBlockNumber(gomock.Eq(1)).After(processed1)

	processed2 := tc.mockTransactionService.EXPECT().ProcessBlockTransactionsByBlockNumber(gomock.Any(), gomock.Eq(2)).Return(nil).After(setCurrent1)
	tc.mockBlockStorage.EXPECT().SetCurrentBlockNumber(gomock.Eq(2)).After(processed2)

	// act
	err := tc.blockService.ProcessNewBlocks(context.Background())

	// assert
	assert.NoError(t, err)
}

func TestBlockServiceProcessNewBlocksFailedProcess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// arrange
	tc := newUnitBlockService(ctrl)

	block := rpc.Block{
		Number: "0x2",
	}

	// assert
	tc.mockClient.EXPECT().GetBlockByNumber(gomock.Any(), gomock.Eq("latest")).Return(&block, nil)
	tc.mockBlockStorage.EXPECT().GetCurrentBlockNumber().Return(1, nil)

	processed1 := tc.mockTransactionService.EXPECT().ProcessBlockTransactionsByBlockNumber(gomock.Any(), gomock.Eq(1)).Return(nil)
	setCurrent1 := tc.mockBlockStorage.EXPECT().SetCurrentBlockNumber(gomock.Eq(1)).After(processed1)

	tc.mockTransactionService.EXPECT().ProcessBlockTransactionsByBlockNumber(gomock.Any(), gomock.Eq(2)).
		Return(errors.New("failed to process")).
		After(setCurrent1)
	tc.mockBlockStorage.EXPECT().SetCurrentBlockNumber(gomock.Eq(2)).Times(0)

	// act
	err := tc.blockService.ProcessNewBlocks(context.Background())

	// assert
	assert.Error(t, err)
}

func TestBlockServiceProcessNewBlocksFailedGetLatest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// arrange
	tc := newUnitBlockService(ctrl)

	// assert
	tc.mockClient.EXPECT().GetBlockByNumber(gomock.Any(), gomock.Eq("latest")).Return(nil, errors.New("failed to get block"))
	tc.mockTransactionService.EXPECT().ProcessBlockTransactionsByBlockNumber(gomock.Any(), gomock.Any()).Times(0)
	tc.mockBlockStorage.EXPECT().SetCurrentBlockNumber(gomock.Any()).Times(0)

	// act
	err := tc.blockService.ProcessNewBlocks(context.Background())

	// assert
	assert.Error(t, err)
}

func TestBlockServiceProcessNewBlocksCurrentNotSet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// arrange
	tc := newUnitBlockService(ctrl)

	block := rpc.Block{
		Number: "0x2",
	}

	// assert
	tc.mockClient.EXPECT().GetBlockByNumber(gomock.Any(), gomock.Eq("latest")).Return(&block, nil)
	tc.mockBlockStorage.EXPECT().GetCurrentBlockNumber().Return(0, storage.ErrBlockCurrentNotSet)

	processed := tc.mockTransactionService.EXPECT().ProcessBlockTransactionsByBlockNumber(gomock.Any(), gomock.Eq(2)).Return(nil)
	tc.mockBlockStorage.EXPECT().SetCurrentBlockNumber(gomock.Eq(2)).Times(1).After(processed)

	// act
	err := tc.blockService.ProcessNewBlocks(context.Background())

	// assert
	assert.NoError(t, err)
}

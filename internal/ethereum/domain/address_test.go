package domain_test

import (
	"testing"
	"trust_walet/internal/ethereum/domain"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	mockDomain "trust_walet/internal/ethereum/domain/mock"
)

func TestAddressServiceAddUnique(t *testing.T) {
	testCases := map[string]struct {
		exist           bool
		expectedAddCall int
	}{
		"addr unique": {
			exist:           false,
			expectedAddCall: 1,
		},
		"addr not unique": {
			exist:           true,
			expectedAddCall: 0,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// arrange
			mockStorage := mockDomain.NewMockAddressStorage(ctrl)
			service := domain.NewAddressService(mockStorage)

			// assert
			mockStorage.EXPECT().Exists(gomock.Eq("addr")).
				Return(tc.exist).
				Times(1)

			mockStorage.EXPECT().Add(gomock.Eq("addr")).Times(tc.expectedAddCall)

			// act
			service.AddUnique("addr")
		})
	}
}

func TestAddressServiceIsSubscribed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// arrange
	mockStorage := mockDomain.NewMockAddressStorage(ctrl)
	service := domain.NewAddressService(mockStorage)

	// assert
	mockStorage.EXPECT().Exists(gomock.Eq("addr")).Return(true)

	// act
	assert.True(t, service.IsSubscribed("addr"))
}

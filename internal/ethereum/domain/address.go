package domain

import "github.com/sirupsen/logrus"

type (
	AddressStorage interface {
		Exists(address string) bool
		Add(address string)
	}

	AddressService struct {
		storage AddressStorage
	}
)

func NewAddressService(addressStorage AddressStorage) *AddressService {
	return &AddressService{
		storage: addressStorage,
	}
}

func (a *AddressService) AddUnique(address string) bool {
	if a.storage.Exists(address) {
		logrus.WithFields(logrus.Fields{
			"address": address,
		}).Warn("address is not unique")

		return false
	}

	logrus.WithFields(logrus.Fields{
		"address": address,
	}).Info("address is added to subscribe list")

	a.storage.Add(address)

	return false
}

func (a *AddressService) IsSubscribed(address string) bool {
	subscribed := a.storage.Exists(address)

	message := "address is in subscribe list"
	if !subscribed {
		message = "address is not in subscribe list"
	}
	logrus.WithFields(logrus.Fields{
		"address": address,
	}).Debug(message)

	return subscribed
}

.PHONY: run
run:
	go run ./cmd/main.go

tests:
	go test ./... -v

mocks:
	mockgen -source internal/ethereum/domain/address.go -destination internal/ethereum/domain/mock/address.go -package=mockDomain
	mockgen -source internal/ethereum/domain/block.go -destination internal/ethereum/domain/mock/block.go -package=mockDomain
	mockgen -source internal/ethereum/domain/transaction.go -destination internal/ethereum/domain/mock/transaction.go -package=mockDomain
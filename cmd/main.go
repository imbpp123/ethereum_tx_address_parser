package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"trust_walet/internal/ethereum"
	"trust_walet/internal/ethereum/domain"
	"trust_walet/internal/ethereum/rpc"
	"trust_walet/internal/ethereum/storage"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.WarnLevel)
	//logrus.SetLevel(logrus.InfoLevel)
	//logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true)

	address := "0xe7d36d7f5832349f7a9f04c898a1e47992f02bd5"
	ctx, cancel := context.WithCancel(context.Background())

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	parser := createParser(rpc.EthereumUrl)

	parser.Subscribe(address)

	fmt.Printf("Collecting transactions for %s address...\n", address)

	go func() {
		<-signalChan
		fmt.Println("Received shutdown signal, exiting...")
		cancel()
	}()

	go func(ctx context.Context) {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				fmt.Println("Stopping transaction logger...")
				return
			case <-ticker.C:
				tx := parser.GetTransactions(address)
				for _, t := range tx {
					fmt.Printf("New transaction: hash=%s value=%s from=%s to=%s\n", t.Hash, t.Value, t.From, t.To)
				}
			}
		}
	}(ctx)

	go func(ctx context.Context) {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				fmt.Println("Stopping block logger...")
				return
			case <-ticker.C:
				fmt.Printf("Current block: 0x%x\n", parser.GetCurrentBlock())
			}
		}
	}(ctx)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Stopping transaction monitor...")
			return
		default:
			if err := parser.MonitorTransactions(ctx); err != nil {
				fmt.Printf("Error while monitoring new transactions: %v\n", err)
				cancel()
				return
			}

			time.Sleep(5 * time.Second)
		}
	}
}

func createParser(url string) *ethereum.Parser {
	client := rpc.NewHttp(&http.Client{}, url)

	addressService := domain.NewAddressService(
		storage.NewAddressInMemory(),
	)
	transactionService := domain.NewTransactionService(
		client,
		addressService,
		storage.NewTransactionInMemory(),
	)
	blockService := domain.NewBlockService(
		client,
		storage.NewBlockInMemory(),
		transactionService,
	)

	return ethereum.NewParser(
		addressService,
		blockService,
		transactionService,
	)
}

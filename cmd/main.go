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
)

func main() {
	address := "0xa8dfb8cc7f9843c3e7bec636bd08c3487b72dc40"
	ctx, cancel := context.WithCancel(context.Background())

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	client := ethereum.NewRpc(&http.Client{}, ethereum.RPCUrl)
	storage := ethereum.NewInMemory()
	parser := ethereum.NewParser(client, storage)

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

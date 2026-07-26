package main

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
)

func BlockHeight() (uint64, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return 0, fmt.Errorf("dial %s: %w", rpcURL, err)
	}
	defer client.Close()

	height, err := client.BlockNumber(context.Background())
	if err != nil {
		return 0, fmt.Errorf("block number: %w", err)
	}

	return height, nil
}

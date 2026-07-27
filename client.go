package main

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
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

func Balance(hexAddr string) (*big.Int, error) {
	isAddr := common.IsHexAddress(hexAddr)
	if !isAddr {
		return nil, fmt.Errorf("hexAddr %s: invalid address", hexAddr)
	}

	addr := common.HexToAddress(hexAddr)

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("dial %s: %w", rpcURL, err)
	}
	defer client.Close()

	balance, err := client.BalanceAt(context.Background(), addr, nil)
	if err != nil {
		return nil, fmt.Errorf("balance: %w", err)
	}

	return balance, nil
}

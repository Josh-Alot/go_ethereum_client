package main

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Client struct {
	eth *ethclient.Client
}

type BlockSummary struct {
	Height       uint64
	Hash         string
	PreviousHash string
	StateRoot    string
	Timestamp    uint64
	Transactions int
}

func NewClient(rpcURL string) (*Client, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("dial %s: %w", rpcURL, err)
	}

	return &Client{eth: client}, nil
}

func (client *Client) Close() {
	client.eth.Close()
}

func (client *Client) BlockHeight() (uint64, error) {
	height, err := client.eth.BlockNumber(context.Background())
	if err != nil {
		return 0, fmt.Errorf("block number: %w", err)
	}

	return height, nil
}

func (client *Client) Balance(hexAddr string) (*big.Int, error) {
	isAddr := common.IsHexAddress(hexAddr)
	if !isAddr {
		return nil, fmt.Errorf("hexAddr %s: invalid address", hexAddr)
	}

	addr := common.HexToAddress(hexAddr)

	balance, err := client.eth.BalanceAt(context.Background(), addr, nil)
	if err != nil {
		return nil, fmt.Errorf("balance: %w", err)
	}

	return balance, nil
}

func (client *Client) BlockInfo(height *big.Int) (*BlockSummary, error) {
	block, err := client.eth.BlockByNumber(context.Background(), height)
	if err != nil {
		return nil, fmt.Errorf("blockinfo: %w", err)
	}

	summary := BlockSummary{
		Height:       block.NumberU64(),
		Hash:         block.Hash().Hex(),
		PreviousHash: block.ParentHash().Hex(),
		StateRoot:    block.Root().Hex(),
		Timestamp:    block.Time(),
		Transactions: len(block.Transactions()),
	}

	return &summary, nil
}

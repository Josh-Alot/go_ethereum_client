package main

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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

type TxSummary struct {
	Pending bool
	Hash    string
	Nonce   uint64
	From    string
	To      string
	Value   *big.Int
	Type    uint8
	ChainID *big.Int
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

func (client *Client) TransactionByHash(txHash string) (*TxSummary, error) {
	isTxHash := common.IsHexHash(txHash)
	if !isTxHash {
		return nil, errors.New("tx info valid hash error: invalid tx hash")
	}

	tx, isPending, err := client.eth.TransactionByHash(context.Background(), common.HexToHash(txHash))
	if err != nil {
		return nil, fmt.Errorf("tx info pending error: %w", err)
	}

	txSigner := types.LatestSignerForChainID(tx.ChainId())
	from, err := types.Sender(txSigner, tx)
	if err != nil {
		return nil, fmt.Errorf("tx info signer error: %w", err)
	}

	to := ""
	if tx.To() != nil {
		to = tx.To().Hex()
	}

	newTx := TxSummary{
		Pending: isPending,
		Hash:    tx.Hash().Hex(),
		Nonce:   tx.Nonce(),
		From:    from.Hex(),
		To:      to,
		Value:   tx.Value(),
		Type:    tx.Type(),
		ChainID: tx.ChainId(),
	}

	return &newTx, nil
}

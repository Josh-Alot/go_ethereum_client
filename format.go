package main

import (
	"fmt"
	"math/big"
	"time"
)

var divisor = big.NewInt(1e18)

func FormatWeiToEther(balance *big.Int) string {
	remain := new(big.Int)
	quo := new(big.Int)

	quo.QuoRem(balance, divisor, remain)
	return fmt.Sprintf("%s.%018s", quo, remain)
}

func FormatTimestampDate(timestamp uint64) string {
	return time.Unix(int64(timestamp), 0).UTC().Format(time.RFC3339)
}

func FormatBlockSummary(summary BlockSummary) {
	fmt.Printf("block height: %d\n", summary.Height)
	fmt.Printf("block hash: %s\n", summary.Hash)
	fmt.Printf("previous hash: %s\n", summary.PreviousHash)
	fmt.Printf("state root hash: %s\n", summary.StateRoot)
	fmt.Printf("block timestamp: %d (%s)\n", summary.Timestamp, FormatTimestampDate(summary.Timestamp))
	fmt.Printf("total transactions: %d\n", summary.Transactions)
}

func FormatTxSummary(tx TxSummary) {
	fmt.Printf("Transaction hash: %s\n", tx.Hash)
	fmt.Printf("Nonce: %d\n", tx.Nonce)

	if !tx.Pending {
		fmt.Println("Status: validated")
	} else {
		fmt.Println("Status: validating")
	}

	fmt.Printf("From: %s\n", tx.From)

	if tx.To != "" {
		fmt.Printf("To: %s\n", tx.To)
	} else {
		fmt.Println("To: Transaction is a contract creation and it does not have a destiny")
	}

	fmt.Printf("Value: %s ETH \n", FormatWeiToEther(tx.Value))
	fmt.Printf("Type: %d\n", tx.Type)
	fmt.Printf("Chain ID: %d\n", tx.ChainID)
}

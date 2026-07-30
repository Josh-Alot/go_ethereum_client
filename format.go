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

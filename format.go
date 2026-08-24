package main

import (
	"fmt"
	"math/big"
	"strings"
	"time"
)

var divisor = big.NewInt(1e18)

func FormatWeiToEther(balance *big.Int) string {
	remain := new(big.Int)
	quo := new(big.Int)

	quo.QuoRem(balance, divisor, remain)
	return fmt.Sprintf("%s.%018s", quo, remain)
}

func ParseEtherToWei(amount string) (*big.Int, error) {
	before, after, _ := strings.Cut(amount, ".")
	if len(after) > 18 {
		return nil, fmt.Errorf("parse ether to wei error: invalid amount, the decimals must be at max 18")
	}

	after = after + strings.Repeat("0", (18-len(after)))
	parsedWei := before + after

	wei, ok := new(big.Int).SetString(parsedWei, 10)
	if !ok {
		return nil, fmt.Errorf("parse ether to wei error: failed to parse %s", amount)
	}

	if wei.Sign() <= 0 {
		return nil, fmt.Errorf("parse ether to wei error: invalid amount")
	}

	return wei, nil
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
	fmt.Printf("From: %s\n", tx.From)

	if tx.To != "" {
		fmt.Printf("To: %s\n", tx.To)
	} else {
		fmt.Println("To: contract creation (address assigned when mined)")
	}

	if tx.Receipt == nil {
		fmt.Println("Status: validating")
	} else {
		fmt.Printf("Block number: %d\n", tx.Receipt.BlockNumber)

		if tx.Receipt.ContractAddress != "" {
			fmt.Printf("Contract deployment address: %s\n", tx.Receipt.ContractAddress)
		}

		if tx.Receipt.Status == 0 {
			fmt.Println("Status: reverted")
		} else {
			fmt.Println("Status: completed")
		}

		fmt.Printf("Gas Used: %d Gas\n", tx.Receipt.GasUsed)
		fmt.Printf("Effective Gas Price: %d Wei\n", tx.Receipt.EffectiveGasPrice)

		bigGasUsed := new(big.Int).SetUint64(tx.Receipt.GasUsed)
		fmt.Printf("Transaction cost: %d Wei\n", new(big.Int).Mul(bigGasUsed, tx.Receipt.EffectiveGasPrice))
	}

	fmt.Printf("Value: %s ETH\n", FormatWeiToEther(tx.Value))
	fmt.Printf("Type: %d\n", tx.Type)
	fmt.Printf("Chain ID: %d\n", tx.ChainID)
}

func FormatAccountCreated(account AccountSummary) {
	fmt.Printf("Account created: %s\n", account.Address)
}

func FormatAccountsList(accounts []AccountSummary) {
	if len(accounts) == 0 {
		fmt.Println("No accounts found")
	} else {
		fmt.Println("Accounts available:")
		for _, account := range accounts {
			fmt.Printf("%s\n", account.Address)
		}
	}
}

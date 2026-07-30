package main

import (
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: blocknumber, getbalance, blockinfo")
		os.Exit(1)
	}

	client, err := NewClient(rpcURL)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	defer client.Close()

	switch os.Args[1] {
	case "blocknumber":
		blocknumber := flag.NewFlagSet("blocknumber", flag.ExitOnError)
		blocknumber.Parse(os.Args[2:])

		height, err := client.BlockHeight()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%d\n", height)

	case "getbalance":
		getbalance := flag.NewFlagSet("getbalance", flag.ExitOnError)

		address := getbalance.String("address", "", "the required address to check balance")
		getbalance.Parse(os.Args[2:])

		if *address == "" {
			fmt.Printf("a wallet address is required\n")
			os.Exit(1)
		}

		balance, err := client.Balance(*address)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s\n", FormatWeiToEther(balance))

	case "blockinfo":
		blockinfo := flag.NewFlagSet("blockinfo", flag.ExitOnError)
		height := blockinfo.Int("height", -1, "the block height on chain")

		blockinfo.Parse(os.Args[2:])

		if *height < 0 {
			fmt.Printf("must give a block height\n")
			os.Exit(1)
		}

		summary, err := client.BlockInfo(big.NewInt(int64(*height)))
		if err != nil {
			log.Fatal(err)
		}

		FormatBlockSummary(*summary)

	default:
		fmt.Printf("command not found\n")
		os.Exit(1)
	}
}

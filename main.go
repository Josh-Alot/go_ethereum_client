package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: blocknumber, getbalance") // name if after
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

	default:
		fmt.Printf("command not found\n")
		os.Exit(1)
	}
}

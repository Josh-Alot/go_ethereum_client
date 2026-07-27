package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: [not defined yet]") // name if after
		os.Exit(1)
	}

	switch os.Args[1] {
	case "blocknumber":
		blocknumber := flag.NewFlagSet("blocknumber", flag.ExitOnError)
		blocknumber.Parse(os.Args[2:])

		height, err := BlockHeight()
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

		balance, err := Balance(*address)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s\n", FormatWeiToEther(balance))

	default:
		fmt.Printf("command not found\n")
		os.Exit(1)
	}
}

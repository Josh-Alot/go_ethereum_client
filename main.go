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
		fmt.Fprintf(os.Stderr, "usage: blocknumber, getbalance, blockinfo, txinfo, createaccount, listaccounts, sendeth\n")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "blocknumber":
		blocknumber := flag.NewFlagSet("blocknumber", flag.ExitOnError)
		blocknumber.Parse(os.Args[2:])

		client, err := createClient()
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close()

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
			fmt.Fprintf(os.Stderr, "a wallet address is required\n")
			os.Exit(1)
		}

		client, err := createClient()
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close()

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
			fmt.Fprintf(os.Stderr, "must give a block height\n")
			os.Exit(1)
		}

		client, err := createClient()
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close()

		summary, err := client.BlockInfo(big.NewInt(int64(*height)))
		if err != nil {
			log.Fatal(err)
		}

		FormatBlockSummary(*summary)

	case "txinfo":
		txinfo := flag.NewFlagSet("txinfo", flag.ExitOnError)
		txHash := txinfo.String("hash", "", "the transaction hash")

		txinfo.Parse(os.Args[2:])

		if *txHash == "" {
			fmt.Fprintf(os.Stderr, "must give a transaction hash\n")
			os.Exit(1)
		}

		client, err := createClient()
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close()

		tx, err := client.TransactionByHash(*txHash)
		if err != nil {
			log.Fatal(err)
		}

		FormatTxSummary(*tx)

	case "createaccount":
		createaccount := flag.NewFlagSet("createaccount", flag.ExitOnError)
		password := createaccount.String("password", "", "the encryption password")

		createaccount.Parse(os.Args[2:])

		var pw []byte
		var err error

		if *password != "" {
			fmt.Fprintln(os.Stderr, "CAUTION: you provided the password on the terminal, it can be leaked through the bash history, use createaccount without the password flag")
			pw = []byte(*password)
		} else {
			pw, err = ReadPassword("Enter the keystore password: ")
			if err != nil {
				log.Fatal(err)
			}
		}

		account, err := CreateAccount(string(pw))
		if err != nil {
			log.Fatal(err)
		}

		FormatAccountCreated(*account)

	case "listaccounts":
		listaccounts := flag.NewFlagSet("listaccounts", flag.ExitOnError)
		listaccounts.Parse(os.Args[2:])

		accounts := ListAccounts()
		FormatAccountsList(accounts)

	case "sendeth":
		sendeth := flag.NewFlagSet("sendeth", flag.ExitOnError)
		from := sendeth.String("from", "", "the required origin address")
		to := sendeth.String("to", "", "the required destination address")
		amount := sendeth.String("amount", "", "the required value to send on ether")

		sendeth.Parse(os.Args[2:])

		if *from == "" {
			fmt.Fprintf(os.Stderr, "the origin address is required\n")
			os.Exit(1)
		}

		if *to == "" {
			fmt.Fprintf(os.Stderr, "the destiny address is required\n")
			os.Exit(1)
		}

		if *amount == "" {
			fmt.Fprintf(os.Stderr, "the amount value to send is required\n")
			os.Exit(1)
		}

		password, err := ReadPassword("Enter the keystore password: ")
		if err != nil {
			log.Fatal(err)
		}

		key, err := LoadPrivateKey(*from, password)
		if err != nil {
			log.Fatal(err)
		}

		client, err := createClient()
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close()

		tx, err := client.PrepareTx(*from, *to, *amount)
		if err != nil {
			log.Fatal(err)
		}

		signedTx, err := SignTx(tx, key)
		if err != nil {
			log.Fatal(err)
		}

		err = client.SendTx(signedTx)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Transaction %v sent to the network\n", signedTx.Hash())
		fmt.Printf("If you want to know the tx status, type txinfo -hash %v\n", signedTx.Hash())

	default:
		fmt.Fprintf(os.Stderr, "command not found\n")
		os.Exit(1)
	}
}

func createClient() (*Client, error) {
	defaultRpcUrl := "http://127.0.0.1:8545"

	rpc := os.Getenv("ETH_RPC_URL")
	if rpc == "" {
		fmt.Fprint(os.Stderr, "CAUTION: You are using anvil, not the sepolia network\n")
		rpc = defaultRpcUrl
	}

	client, err := NewClient(rpc)
	if err != nil {
		return nil, fmt.Errorf("client connection error: %w", err)
	}

	return client, nil
}

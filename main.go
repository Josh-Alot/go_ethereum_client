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

		height, err := BlockNumber()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%d\n", height)
	default:
		fmt.Printf("command not found\n")
		os.Exit(1)
	}
}

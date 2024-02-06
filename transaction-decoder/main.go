package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"transaction-decoder/bitcoin"
)

// main is the entry point of the program.
func main() {
	version := flag.Bool("v", false, "Show Version")
	inputs := flag.Bool("in", false, "Show Inputs")
	outputs := flag.Bool("out", false, "Show Outputs")
	locktime := flag.Bool("lock-time", false, "Show Locktime")
	trxHex := flag.String("hex", "", "Transaction hex")

	bitcoinClient := bitcoin.BitcoinClient()

	flag.Parse()

	if *trxHex == "" {
		fmt.Println("Please provide transaction Hex")
		fmt.Println("-hex is required!")
		return
	}
	byteString, err := hex.DecodeString(*trxHex)
	if err != nil {
		fmt.Println("Error decoding hex", err)
		return
	}

	decodedTrx, err := bitcoinClient.DecodeRawTransaction(byteString)
	if err != nil {
		fmt.Println("Error decoding transaction", err)
		return
	}

	displayTrxDetails(*decodedTrx, *version, *inputs, *outputs, *locktime)
}

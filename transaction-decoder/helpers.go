package main

import (
	"fmt"

	"github.com/btcsuite/btcd/btcjson"
)

func displayTrxDetails(decodedTrx btcjson.TxRawResult, version, inputs, outputs, locktime bool) {
	fmt.Println("================================= START =================================")
	if version {
		fmt.Printf("Version: %v\n", decodedTrx.Version)
	}
	if inputs {
		fmt.Printf("Inputs: %v\n", decodedTrx.Vin)
	}
	if outputs {
		fmt.Printf("Outputs: %v\n", decodedTrx.Vout)
	}
	if locktime {
		fmt.Printf("Locktime: %v\n", decodedTrx.LockTime)
	}
	fmt.Println("================================= END =================================")
}

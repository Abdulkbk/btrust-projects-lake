package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

const MAX_WEIGHT = 4000000

var currentWeight = 0

var fees = 0

var trxTxids = []string{}

func main() {
	records := readMempool()

	// for _, record := range records {
	// 	fmt.Println(strings.Split(record[0], ",")[3])
	// }

	records = sortMempool(records, 1)

	fmt.Println(records[0])

	constructBlock(records)

	fmt.Println("Total fees:", fees)
	fmt.Println("Total weight:", currentWeight)
	fmt.Println("Total transactions:", len(trxTxids))
	fmt.Println("Total mempool transactions:", len(records))

	fmt.Println("=================================  Start  =================================")
	joinedTrx := strings.Join(trxTxids, "\n")
	fmt.Println(joinedTrx[:1039])
	fmt.Println("=================================   End   =================================")

}

func readMempool() [][]string {
	file, err := os.Open("mempool.csv")
	if err != nil {
		fmt.Println("Error when opening file", err)
		return nil
	}

	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error when reading file", err)
		return nil
	}

	return records

}

func constructBlock(mempool [][]string) {
	for _, record := range mempool {
		txid, fee, weight, parents := record[0], record[1], record[2], record[3]

		if parents != "" {
			parentsAppeared := checkParentsAppeared(parents)

			if !parentsAppeared {
				continue
			}
		}

		weightToInt, err := strconv.Atoi(weight)
		if err != nil {
			fmt.Println("Error when converting weight", err)
			return
		}

		feeToInt, err := strconv.Atoi(fee)
		if err != nil {
			fmt.Println("Error when converting fee", err)
			return
		}

		if currentWeight+weightToInt > MAX_WEIGHT {
			continue
		}

		trxTxids = append(trxTxids, txid)
		currentWeight += weightToInt
		fees += feeToInt
	}
}

// func destructure(trx []string) (string, string, string, string) {

// 	data := strings.Split(trx[0], ",")

// 	txid, fee, weight, parents := data[0], data[1], data[2], data[3]

// 	return txid, fee, weight, parents

// }

func checkParentsAppeared(parentsTxid string) bool {
	pTxids := strings.Split(parentsTxid, ";")

	for _, pTxid := range pTxids {
		exists := contains(trxTxids, pTxid)
		if !exists {
			return false
		}
	}
	return true
}

func sortMempool(mempool [][]string, index int) [][]string {
	records := [][]string{}

	for _, record := range mempool {
		des := strings.Split(record[0], ",")
		records = append(records, des)
	}

	sort.Slice(records, func(i, j int) bool {
		a, _ := strconv.Atoi(records[i][index])
		b, _ := strconv.Atoi(records[j][index])

		return a > b
	})

	return records

}

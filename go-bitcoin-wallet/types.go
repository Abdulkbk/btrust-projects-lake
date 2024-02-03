package main

type UTXO struct {
	Status struct {
		BlockHash   string `json:"block_hash"`
		BlockHeight int    `json:"block_height"`
		BlockTime   int    `json:"block_time"`
		Confirmed   bool   `json:"confirmed"`
	} `json:"status"`
	Txid  string `json:"txid"`
	Value int64  `json:"value"` // Assuming integer value
	Vout  int    `json:"vout"`
}

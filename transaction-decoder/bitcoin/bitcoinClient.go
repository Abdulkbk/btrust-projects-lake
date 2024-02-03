package bitcoin

import (
	"log"
	"transaction-decoder/constants"

	"github.com/btcsuite/btcd/rpcclient"
)

func BitcoinClient() *rpcclient.Client {

	connConfig := &rpcclient.ConnConfig{
		Host:         "127.0.0.1:" + constants.REGTEST_PORT,
		User:         "user",
		Pass:         "pass",
		HTTPPostMode: true,
		DisableTLS:   true,
	}

	client, err := rpcclient.New(connConfig, nil)

	if err != nil {
		log.Fatal("Error connecting to bitcoind: ", err)
	}

	// blockChainInfo, err2 := client.GetBlockChainInfo()

	// if err2 != nil {
	// 	log.Fatal("Error getting blockchaininfo: ", err2)
	// }

	// log.Println(blockChainInfo)
	// fmt.Println(blockChainInfo)
	return client

}

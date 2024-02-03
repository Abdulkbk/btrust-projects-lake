package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// const TESTNET_RPC_URL = "http://127.0.0.1:18332"
const SIGNET_RPC_URL = "http://127.0.0.1:38332"

const USERNAME = "user"
const PASSWORD = "pass"

var auth = USERNAME + ":" + PASSWORD
var authEnc = base64.StdEncoding.EncodeToString([]byte(auth))

func main() {
	fmt.Println("Welcome to connecting to bitcoin")

	r := gin.Default()

	r.GET("/get_utxos", handleGetUtxos)
	r.GET("/get_balance", handleGetBalance)
	r.GET("/get_new_address", handleGetNewAddress)
	r.GET("/get_address_info/:address", handleGetAddressInfo)

	r.Run("localhost:4000")
}

func handleGetUtxos(c *gin.Context) {
	jsonStr := []byte(`{"jsonrpc": "1.0", "id": "curltest", "method": "listunspent", "params": []}`)
	auth := USERNAME + ":" + PASSWORD
	authEnc := base64.StdEncoding.EncodeToString([]byte(auth))

	req, err := http.NewRequest("POST", SIGNET_RPC_URL, bytes.NewBuffer(jsonStr))

	if err != nil {
		fmt.Println("Error creating request")
		return
	}

	req.Header.Set("Authorization", "Basic "+authEnc)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request")
		return
	}

	defer resp.Body.Close()

	var utxos any

	err = json.NewDecoder(resp.Body).Decode(&utxos)

	if err != nil {
		fmt.Println("Error decoding response")
		return
	}

	c.IndentedJSON(200, utxos)

}
func handleGetBalance(c *gin.Context) {
	jsonStr := []byte(`{"jsonrpc": "1.0", "id": "curltest", "method": "getbalance", "params": []}`)

	req, err := http.NewRequest("POST", SIGNET_RPC_URL, bytes.NewBuffer(jsonStr))

	if err != nil {
		fmt.Println("Error creating request")
		return
	}

	req.Header.Set("Authorization", "Basic "+authEnc)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request")
		return
	}

	defer resp.Body.Close()

	var walletBalance any

	err = json.NewDecoder(resp.Body).Decode(&walletBalance)

	if err != nil {
		fmt.Println("Error decoding response")
		return
	}

	c.IndentedJSON(200, walletBalance)

}
func handleGetAddressInfo(c *gin.Context) {
	address := c.Param("address")
	jsonStr := []byte(`{"jsonrpc": "1.0", "id": "curltest", "method": "getaddressinfo", "params": ["` + address + `"]}`)

	req, err := http.NewRequest("POST", SIGNET_RPC_URL, bytes.NewBuffer(jsonStr))

	if err != nil {
		fmt.Println("Error creating request")
		return
	}

	req.Header.Set("Authorization", "Basic "+authEnc)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request")
		return
	}

	defer resp.Body.Close()

	var addressInfo any

	err = json.NewDecoder(resp.Body).Decode(&addressInfo)

	if err != nil {
		fmt.Println("Error decoding response")
		return
	}

	c.IndentedJSON(200, addressInfo)

}
func handleGetNewAddress(c *gin.Context) {
	addrType := c.Query("addressType")

	jsonStr := []byte(`{"jsonrpc": "1.0", "id": "curltest", "method": "getnewaddress", "params": ["` + addrType + `"]}`)

	req, err := http.NewRequest("POST", SIGNET_RPC_URL, bytes.NewBuffer(jsonStr))

	if err != nil {
		fmt.Println("Error creating request")
		return
	}

	req.Header.Set("Authorization", "Basic "+authEnc)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request")
		return
	}

	defer resp.Body.Close()

	var newAddress any

	err = json.NewDecoder(resp.Body).Decode(&newAddress)

	if err != nil {
		fmt.Println("Error decoding response")
		return
	}

	c.IndentedJSON(200, newAddress)

}

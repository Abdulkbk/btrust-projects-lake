package main

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"

// 	"github.com/btcsuite/btcd/btcutil"
// 	"github.com/btcsuite/btcd/btcutil/hdkeychain"
// 	"github.com/btcsuite/btcd/chaincfg"
// 	"github.com/btcsuite/btcd/chaincfg/chainhash"
// 	"github.com/btcsuite/btcd/txscript"
// 	"github.com/btcsuite/btcd/wire"
// 	"github.com/gin-gonic/gin"
// 	"github.com/tyler-smith/go-bip32"
// 	"github.com/tyler-smith/go-bip39"
// )

// var RPC_URL = "http://127.0.0.1:18332"
// var RPC_BITSTREAM = "https://blockstream.info/testnet/api/"

// type UTXO struct {
// 	Status struct {
// 		BlockHash   string `json:"block_hash"`
// 		BlockHeight int    `json:"block_height"`
// 		BlockTime   int    `json:"block_time"`
// 		Confirmed   bool   `json:"confirmed"`
// 	} `json:"status"`
// 	Txid  string `json:"txid"`
// 	Value int64  `json:"value"` // Assuming integer value
// 	Vout  int    `json:"vout"`
// }

// type trxRequest struct {
// 	PrivKey      string `json:"privKey"`
// 	To           string `json:"to"`
// 	Amount       int    `json:"amount"`
// 	MasterPubKey string `json:"masterPubkey"`
// }

// func depmain() {
// 	r := gin.Default()

// 	r.GET("/", func(c *gin.Context) {
// 		c.JSON(200, gin.H{
// 			"message": "Hello World",
// 		})
// 	})

// 	r.GET("/generate_keys", handleGenerateAddress)
// 	r.GET("/generate_p2pkh", handleGenerateP2PKHAddress)
// 	r.GET("/generate_p2wpkh", handleGenerateP2WPKHAddress)
// 	r.GET("/test", handleTest)
// 	r.POST("/create_trx", handleCreateTransaction)

// 	r.Run("localhost:4000")
// }

// func handleGenerateAddress(c *gin.Context) {
// 	entropy, _ := bip39.NewEntropy(128)

// 	mnemonics, _ := bip39.NewMnemonic(entropy)

// 	seed := bip39.NewSeed(mnemonics, "")

// 	masterKey, _ := bip32.NewMasterKey(seed)

// 	pubKey := masterKey.PublicKey()

// 	testnetKey, _ := masterKey.NewChildKey(bip32.FirstHardenedChild + 44)

// 	derivedPrivKey := testnetKey.Key

// 	derivedPubKey := testnetKey.PublicKey()

// 	c.IndentedJSON(200, gin.H{
// 		"mnemonics":      mnemonics,
// 		"masterKey":      masterKey,
// 		"pubKey":         pubKey,
// 		"derivedPrivKey": string(derivedPrivKey),
// 		"derivedPubKey":  derivedPubKey.String(),
// 	})

// }

// func handleGenerateP2PKHAddress(c *gin.Context) {
// 	masterPubKey := c.Query("masterPubkey")

// 	masterKey, _ := hdkeychain.NewKeyFromString(masterPubKey)

// 	childKey, _ := masterKey.Derive(0)

// 	pubKey, _ := childKey.ECPubKey()

// 	pubKeyHash := btcutil.Hash160(pubKey.SerializeCompressed())

// 	testnetPParams := &chaincfg.TestNet3Params

// 	p2pkhAddr, _ := btcutil.NewAddressPubKeyHash(pubKeyHash, testnetPParams)

// 	c.IndentedJSON(200, gin.H{
// 		"p2pkhAddr": p2pkhAddr.String(),
// 	})
// }

// func handleGenerateP2WPKHAddress(c *gin.Context) {
// 	masterPubKey := c.Query("masterPubkey")

// 	masterKey, _ := hdkeychain.NewKeyFromString(masterPubKey)

// 	childKey, _ := masterKey.Derive(0)

// 	pubKey, _ := childKey.ECPubKey()

// 	pubKeyHash := btcutil.Hash160(pubKey.SerializeCompressed())

// 	testnetPParams := &chaincfg.TestNet3Params

// 	p2wpkhAddr, _ := btcutil.NewAddressWitnessPubKeyHash(pubKeyHash, testnetPParams)

// 	c.IndentedJSON(200, gin.H{
// 		"p2wpkhAddr": p2wpkhAddr.String(),
// 	})
// }

// func handleTest(c *gin.Context) {

// 	body := map[string]interface{}{
// 		"jsonrpc": "1.0",
// 		"id":      "curltext",
// 		"method":  "listunspent",
// 		"params":  []string{},
// 	}

// 	jsonBody, _ := json.Marshal(body)

// 	req, _ := http.NewRequest("POST", RPC_URL, bytes.NewBuffer(jsonBody))

// 	req.SetBasicAuth("user", "pass")

// 	req.Header.Set("Content-Type", "application/json")

// 	client := &http.Client{}

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		panic(err)
// 	}

// 	defer resp.Body.Close()

// 	var responseData interface{}
// 	err = json.NewDecoder(resp.Body).Decode(&responseData)

// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println(responseData)

// 	c.IndentedJSON(200, responseData)
// }

// func handleCreateTransaction(c *gin.Context) {
// 	var req trxRequest

// 	if err := c.BindJSON(&req); err != nil {
// 		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	masterKey, _ := hdkeychain.NewKeyFromString(req.MasterPubKey)

// 	childKey, _ := masterKey.Derive(0)

// 	pubKey, _ := childKey.ECPubKey()

// 	pubKeyHash := btcutil.Hash160(pubKey.SerializeCompressed())

// 	testnetPParams := &chaincfg.TestNet3Params

// 	p2pkhAddr, _ := btcutil.NewAddressPubKeyHash(pubKeyHash, testnetPParams)

// 	fmt.Println(p2pkhAddr.String())

// 	request, _ := http.NewRequest("GET", RPC_BITSTREAM+"address/"+p2pkhAddr.String()+"/utxo", nil)

// 	client := &http.Client{}

// 	resp, err := client.Do(request)
// 	if err != nil {
// 		panic(err)
// 	}

// 	defer resp.Body.Close()

// 	var utxos []UTXO

// 	err = json.NewDecoder(resp.Body).Decode(&utxos)

// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println("")
// 	fmt.Println(utxos[0].Txid)
// 	fmt.Println("")

// 	txid, _ := chainhash.NewHashFromStr(utxos[0].Txid)
// 	prevOut := wire.OutPoint{Hash: *txid, Index: 1}
// 	amountToSend := btcutil.Amount(req.Amount)
// 	fee := btcutil.Amount(1000)

// 	toAddr, _ := btcutil.DecodeAddress(req.To, &chaincfg.TestNet3Params)
// 	changeAddr, _ := btcutil.DecodeAddress(p2pkhAddr.String(), &chaincfg.TestNet3Params)
// 	recipientPkScript, _ := txscript.PayToAddrScript(toAddr)
// 	changePkScript, _ := txscript.PayToAddrScript(changeAddr)

// 	// Create trx
// 	tx := wire.NewMsgTx(wire.TxVersion)
// 	tx.AddTxIn(&wire.TxIn{
// 		PreviousOutPoint: prevOut,
// 		Sequence:         wire.MaxTxInSequenceNum,
// 	})

// 	tx.AddTxOut(&wire.TxOut{
// 		Value:    int64(amountToSend),
// 		PkScript: recipientPkScript,
// 	})

// 	tx.AddTxOut(&wire.TxOut{
// 		Value:    utxos[0].Value - int64(amountToSend) - int64(fee),
// 		PkScript: changePkScript,
// 	})

// 	// Sign trx
// 	privKey, _ := btcutil.DecodeWIF(req.PrivKey)
// 	sigScript, _ := txscript.SignatureScript(tx, 0, prevOut.Hash[:], txscript.SigHashAll, privKey.PrivKey, true)

// 	tx.TxIn[0].SignatureScript = sigScript

// 	// Serialized Trx
// 	c.IndentedJSON(200, tx)

// }

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

func main() {
	redeemScript := generateRedeemScript()

	fmt.Println(redeemScript)

	p2shAddress := generateAddressFromRedeemScript(redeemScript)

	fmt.Println(p2shAddress)

	prevTxHashStr := "5dc27ef4b998e1b2bcb62bf53d99a83f703ff110977fc63f72194d65e280a7cb"
	p2shAddr := "2NEWkaQQEMGyMyFyFXdgmpo1GHN6qiPhP5t"
	amount := btcutil.Amount(1000000)

	// Calling the transaction constructor to send bitcoin to the p2sh address created
	sendTrx := transactionContructor(prevTxHashStr, p2shAddr, amount)

	// Building a transaction to spend bitcoin from the p2sh address created
	spendTrx := transactionSpenderConstructor(sendTrx.TxHash().String(), btcutil.Amount(100000), "2NEWkaQQEMGyMyFyFXdgmpo1GHN6qiPhP5t", "2NEWkaQQEMGyMyFyFXdgmpo1GHN6qiPhP5t")

	fmt.Println(spendTrx)

}

func generateRedeemScript() string {
	// preImage := "Btrust Builders"

	// Byte encoding of "Btrust BUlders"
	preImage := "427472757374204275696c64657273"

	hash := sha256.Sum256([]byte(preImage))

	lockHex := hex.EncodeToString(hash[:])

	// The redeem script for the p2sh
	redeemScript := fmt.Sprintf("OP_SHA256 %s OP_EQUAL", lockHex)

	return redeemScript
}

func generateAddressFromRedeemScript(redeemScript string) btcutil.Address {
	// Convert the redeem script to bytes
	redeemScriptBytes := []byte(redeemScript)

	// Hash the redeem script with Sha256 to give a 32 byte hash
	hash := sha256.Sum256(redeemScriptBytes)

	// Hash the result with RIPEMD160  to give a 20 byte hash
	ripemd160Hash := btcutil.Hash160(hash[:])

	// Append the 0x05 version prefix
	p2shHash := append([]byte{0x05}, ripemd160Hash...)

	// Double hash the result with Sha256 to give a 32 byte hash
	checksum := sha256.Sum256(p2shHash)
	checksum = sha256.Sum256(checksum[:])

	// Append the checksum for the p2sh address
	p2shHash = append(p2shHash, checksum[:4]...)

	// Convert the p2sh hash to a p2sh address
	p2shAddress, _ := btcutil.NewAddressScriptHash(p2shHash, &chaincfg.SigNetParams)
	return p2shAddress
}

func transactionContructor(prevTxHashStr, p2ScriptAddrStr string, amount btcutil.Amount) *wire.MsgTx {

	tx := wire.NewMsgTx(wire.TxVersion)

	prevTxHash, err := chainhash.NewHashFromStr(prevTxHashStr)

	if err != nil {
		fmt.Println("Error while decoding prev tx hash", err)
		return tx
	}

	// Creating the input for the transaction (referencing a utxo)
	prevTxOutPoint := wire.NewOutPoint(prevTxHash, 0)
	txIn := wire.NewTxIn(prevTxOutPoint, nil, nil)
	tx.AddTxIn(txIn)

	// p2ScriptAddrStr := "2NEWkaQQEMGyMyFyFXdgmpo1GHN6qiPhP5t"
	p2shAddress, err := btcutil.DecodeAddress(p2ScriptAddrStr, &chaincfg.SigNetParams)

	if err != nil {
		fmt.Println("Error while decoding p2sh address", err)
		return tx
	}

	//
	script, err := txscript.PayToAddrScript(p2shAddress)

	if err != nil {
		fmt.Println("Error while creating p2s script", err)
		return tx
	}

	txOut := wire.NewTxOut(int64(amount), script)
	tx.AddTxOut(txOut)

	fmt.Println(tx)
	return tx

}

func transactionSpenderConstructor(prevTxHashStr string, mainAmount btcutil.Amount, to, changeAddressStr string) *wire.MsgTx {
	prevTxHash, err := chainhash.NewHashFromStr(prevTxHashStr)
	tx := wire.NewMsgTx(wire.TxVersion)

	if err != nil {
		fmt.Println("Error while decoding prev tx hash", err)
		return tx
	}

	prevTxOutPoint := wire.NewOutPoint(prevTxHash, 0)

	txIn := wire.NewTxIn(prevTxOutPoint, nil, nil)
	tx.AddTxIn(txIn)

	mainAddress, err := btcutil.DecodeAddress(to, &chaincfg.SigNetParams)

	if err != nil {
		fmt.Println("Error while decoding main address", err)
		return tx
	}

	mainScript, err := txscript.PayToAddrScript(mainAddress)

	if err != nil {
		fmt.Println("Error while creating main script", err)
		return tx
	}

	mainTxOut := wire.NewTxOut(int64(mainAmount), mainScript)
	tx.AddTxOut(mainTxOut)

	inputAmount := int64(10000)
	changeAmount := inputAmount - int64(mainAmount) // Calculate change amount
	changeAddress, err := btcutil.DecodeAddress(changeAddressStr, &chaincfg.SigNetParams)

	if err != nil {
		fmt.Println("Error while decoding change address", err)
		return tx
	}

	changeScript, err := txscript.PayToAddrScript(changeAddress)

	if err != nil {
		fmt.Println("Error while creating change script", err)
		return tx
	}

	changeTxOut := wire.NewTxOut(changeAmount, changeScript)
	tx.AddTxOut(changeTxOut)

	unlockingScript, err := txscript.PayToAddrScript(mainAddress)

	if err != nil {
		fmt.Println("Error while creating unlocking script", err)
		return tx
	}

	tx.TxIn[0].SignatureScript = unlockingScript

	// fmt.Println(tx)
	return tx
}

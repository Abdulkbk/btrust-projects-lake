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

	// ==================== 2
	p2shAddress := generateAddressFromRedeemScript(redeemScript)

	fmt.Println(p2shAddress)

	// ==================== 3
	prevTxHashStr := ""
	p2shAddr := ""
	ammount := btcutil.Amount(100000)
	sendTrx := transactionContructor(prevTxHashStr, p2shAddr, ammount)
	TransactionSpenderConstructor(sendTrx.TxHash().String())

}

func generateRedeemScript() string {
	preImage := "Btrust Builders"

	hash := sha256.Sum256([]byte(preImage))

	redeemScript := hex.EncodeToString(hash[:])
	return redeemScript
}

func generateAddressFromRedeemScript(redeemScript string) btcutil.Address {
	redeemScriptBytes, err := hex.DecodeString(redeemScript)

	if err != nil {
		fmt.Println("Error while decoding redeem script", err)
	}

	hash2 := sha256.Sum256(redeemScriptBytes)

	ripemd160Hash := btcutil.Hash160(hash2[:])

	p2shHash := append([]byte{0x05}, ripemd160Hash...)

	checksum := sha256.Sum256(p2shHash)
	checksum = sha256.Sum256(checksum[:])

	p2shHash = append(p2shHash, checksum[:4]...)

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

	prevTxOutPoint := wire.NewOutPoint(prevTxHash, 0)
	txIn := wire.NewTxIn(prevTxOutPoint, nil, nil)
	tx.AddTxIn(txIn)

	// p2ScriptAddrStr := "2NEWkaQQEMGyMyFyFXdgmpo1GHN6qiPhP5t"
	p2sAddress, err := btcutil.DecodeAddress(p2ScriptAddrStr, &chaincfg.SigNetParams)

	if err != nil {
		fmt.Println("Error while decoding p2s address", err)
		return tx
	}

	script, err := txscript.PayToAddrScript(p2sAddress)

	if err != nil {
		fmt.Println("Error while creating p2s script", err)
		return tx
	}

	txOut := wire.NewTxOut(int64(amount), script)
	tx.AddTxOut(txOut)

	fmt.Println(tx)
	return tx

}

func TransactionSpenderConstructor(prevTxHashStr string, mainAmount btcutil.Amount, to, changeAddressStr string) *wire.MsgTx {
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

	inputAmount := int64(100000000)                 // Total amount from previous transaction
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

	fmt.Println(tx)
	return tx
}

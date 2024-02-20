package main

import (
	"bytes"
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
	// Task 1: Generate the redeem script
	preImage := "427472757374204275696c64657273"
	redeemScript := generateRedeemScript(preImage)
	fmt.Println("Redeem script: ", redeemScript)

	// Task 2: Derive an address from the redeem script
	p2shAddress := generateAddressFromRedeemScript(redeemScript)
	fmt.Println("Address:", p2shAddress)

	prevTxHashStr := "f54b974828be38649b8c056bca4b83d94d50417ff8660f1f27d98567f9c7da1a" // Hash of the previous transaction from local signet wallet
	p2shAddr := "2N74S7ccUsw1hn4w9NASZAgcsrcEzmf6cox"                                   // Address gotten from the p2sh
	amount := btcutil.Amount(10000)

	// Task 3: Create the transaction that sends some bitcoin to an address

	sendTrx := transactionContructor(prevTxHashStr, p2shAddr, amount)

	serializedTx := new(bytes.Buffer) // Create a new buffer to append the serialized transaction
	sendTrx.Serialize(serializedTx)

	hexEnc := hex.EncodeToString(serializedTx.Bytes()) // Convert the buffer to hex
	fmt.Println("Transaction :", hexEnc)

	// Trx hex: 01000000011adac7f96785d9271f0f66f87f41504dd9834bca6b058c9b6438be2848974bf500000000000000000001102700000000000017a91497874ed8e883a4d2a5442a95315035605c7501378700000000
	// Trx signed hex: 01000000011adac7f96785d9271f0f66f87f41504dd9834bca6b058c9b6438be2848974bf5000000006a47304402202a2156e92252b65c8c27fc211b40c524b821833d62566120579366f8e293214a022071d599e35e9028d678b29e7e5e631baeb40bb42c71992d1c359f3c934b111f32012102be8c25ee7213c3d3a7622534e91cf1ccf0c6cc11442bfa340fbe75cf57ae7f4c0000000001102700000000000017a91497874ed8e883a4d2a5442a95315035605c7501378700000000
	// Txid: 7a4ecf5ff85cf58df4197eec8fc7b2eedccd72dea5e09aace8ba3f2bafdac4b2
	// Building a transaction to spend bitcoin from the p2sh address created

	// Task 4: Create another transaction that spends the bitcoin
	prevTxHashStr2 := "7a4ecf5ff85cf58df4197eec8fc7b2eedccd72dea5e09aace8ba3f2bafdac4b2"
	myP2shAddr := "2N74S7ccUsw1hn4w9NASZAgcsrcEzmf6cox"
	myChangeAddr := "tb1qsqhkt8035ljqj53vfd3d56tdz66p7rnsgrguvh"
	spendTrx := transactionSpenderConstructor(prevTxHashStr2, btcutil.Amount(9000), myP2shAddr, myChangeAddr)

	serializedTx2 := new(bytes.Buffer)
	spendTrx.Serialize(serializedTx2)

	hexEnc2 := hex.EncodeToString(serializedTx2.Bytes())

	fmt.Println("Transaction :", hexEnc2)

}

func generateRedeemScript(preImage string) string {
	// preImage := "Btrust Builders"

	// Byte encoding of "Btrust BUlders"

	hash := sha256.Sum256([]byte(preImage))

	lockHex := hex.EncodeToString(hash[:])

	// The redeem script for the p2sh
	redeemScript := fmt.Sprintf("OP_SHA256 %s OP_EQUAL", lockHex)

	return redeemScript
}

func generateAddressFromRedeemScript(redeemScript string) btcutil.Address {
	// Convert the redeem script to bytes array
	redeemScriptBytes := []byte(redeemScript)

	// Hash the result with RIPEMD160  to give a 20 byte hash
	ripemd160Hash := btcutil.Hash160(redeemScriptBytes)

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

	// Add signature
	// sigScript, err := txscript.SignatureScript(tx, 0, script, txscript.SigHashAll, nil, true)

	if err != nil {
		fmt.Println("Error while creating signature script", err)
		return tx
	}

	// txIn.SignatureScript = sigScript
	txIn.Sequence = 0

	return tx

}

func transactionSpenderConstructor(prevTxHashStr string, mainAmount btcutil.Amount, to, changeAddressStr string) *wire.MsgTx {
	actualPreimage := "Btrust Builders"
	hashedPreimage := sha256.Sum256([]byte(actualPreimage))
	unlockingScript := hex.EncodeToString(hashedPreimage[:])

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

	if err != nil {
		fmt.Println("Error while creating unlocking script", err)
		return tx
	}

	// tx.TxIn[0].SignatureScript =

	// fmt.Println(tx)
	return tx
}

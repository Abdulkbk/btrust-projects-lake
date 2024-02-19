package main

import (
	"fmt"
)

// const HEX = "02000000000101d9014dbbdb08a4b93b53d2b62194139d0f85ba20e522d1b4afd92fa39fec562b1f00000000fffffffd014a01000000000000225120245091249f4f29d30820e5f36e1e5d477dc3386144220bd6f35839e94de4b9ca03403373581c4772771f039fd66572ae7524416d5336633bf0c625f405de6ab1d05fb1a018a2448eb858ca6fec63f01949dd127ee39ed2a506645815a6be2366d69edb2055bba80746653a70cc871710c6671a88e4b0035e070e98bf7340d1ffd9b3afc9ac0063036f7264010118746578742f706c61696e3b636861727365743d7574662d38004c947b224269746d616e204944223a3530303831382c22426972746864617465223a22323031372d31322d32342031313a3038222c2253706563696573223a2230783230303030303030222c2253697a65223a313038333831372c22576569676874223a333939333133332c225765616c7468223a3736383438343833392c22576973646f6d223a313837333130353437353232317d6821c055bba80746653a70cc871710c6671a88e4b0035e070e98bf7340d1ffd9b3afc900000000"

const HEX = "0200000000010179aaafbe7c9d3b0812a489facaf77508c08c190ec7dfd82f129aeb995aca23ab0000000000fdffffff020bd2190000000000160014d2caa7b08db89cd62c9af34da53332d30e53bb1598151b00000000001600143d4427468cbe7ae396427a1aa9128fa05b18c7db024730440220573fd27574cfdde484347621e1f48f85ae975cb8c2265a04496ded038896822302204a5e04a3a2d160c3158caa39b58bfc91ac64c484078ec0225a7d4d2d4454661f012103d96e3819b52245e42c76f869c9a875f6ea5344cf1aee2e6b3ab03adcfef0d80ede3b0b00"

func main() {

	// convert hex to decodedBytes
	// byte1, _ := hex.DecodeString(HEX)
	decodedBytes := convertHexToBytes(HEX)

	version := decodedBytes[0]

	witnessFlag := decodedBytes[5]

	noOfInputs := decodedBytes[6]

	endInputIndex := noOfInputs * 43
	inputs := decodedBytes[7:endInputIndex]

	// lenSig := decodedBytes[endInputIndex]

	startSeqIndex := endInputIndex + 1
	stopSeqIndex := startSeqIndex + 4
	sequence := decodedBytes[startSeqIndex:stopSeqIndex]

	noOfOutputs := decodedBytes[stopSeqIndex]

	startOutputIndex := stopSeqIndex + 1
	stopOutputIndex := startOutputIndex + noOfOutputs*31
	outputs := decodedBytes[startOutputIndex:stopOutputIndex]

	witnessCount := decodedBytes[stopOutputIndex]

	var witnessesData [][]byte

	witnessStart := stopOutputIndex + 2

	for i := 0; i < int(witnessCount); i += 1 {
		witnessStop := witnessStart + decodedBytes[witnessStart-1]
		// fmt.Println("witnessStop: ", witnessStop, witnessStart, decodedBytes[witnessStart-1])
		witnessesData = append(witnessesData, decodedBytes[witnessStart:witnessStop])
		witnessStart = witnessStop + 1
	}

	// fmt.Println("Witness stop: ", witnessStart)
	locktime := decodedBytes[witnessStart-1:]

	// fmt.Println(decodedBytes)
	fmt.Printf("version: %x\n", version)
	fmt.Println("witnessFlag: ", witnessFlag)

	inputsStringRep := fmt.Sprintf("%x", inputs)
	fmt.Println("inputs: ", inputsStringRep)
	// fmt.Println("lenSig: ", lenSig)

	seqStrRep := fmt.Sprintf("%x", sequence)
	fmt.Println("sequence: ", seqStrRep)

	outStrRep := fmt.Sprintf("%x", outputs)
	fmt.Println("outputs: ", outStrRep)

	witness1StrRep := fmt.Sprintf("%x", witnessesData[0])
	witness2StrRep := fmt.Sprintf("%x", witnessesData[1])
	fmt.Println("witnessesData1: ", witness1StrRep)
	fmt.Println("witnessesData2: ", witness2StrRep)

	stringRep := fmt.Sprintf("%x", locktime)
	fmt.Println("locktime: ", stringRep)
}

func convertHexToBytes(hex string) []byte {

	// create a byte array
	var result []byte

	for i := 0; i < len(hex); i += 2 {

		byteValue := (byteFromHexChar(hex[i]) << 4) | byteFromHexChar(hex[i+1])

		result = append(result, byteValue)
	}

	return result
}

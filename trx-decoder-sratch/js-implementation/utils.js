const convertEndian = (hexString) => {
  // Add padding if the length of hexString is odd
  if (hexString.length % 2 !== 0) {
    hexString = '0' + hexString;
  }
  const convertedHex = hexString.match(/../g).reverse().join('');
  return convertedHex;
}


const decodeFixedLength = (hexValue) => {
  const bigEndianHex = convertEndian(hexValue);
  return [parseInt(bigEndianHex, 16), hexValue.length];
}

const decodeCompactSize = (hexValue) => {

  let skip = 2
  const firstByte = parseInt(hexValue.substring(0, 2), 16);

  if (firstByte < 253) {
    return [firstByte, skip];
  } else if (firstByte === 253) {
    skip = 4
    const [decValue, hexSkip] = decodeFixedLength(hexValue.substring(2, 6));
    return [decValue, skip];
  } else if (firstByte === 254) {
    skip = 8
    const [decValue, hexSkip] = decodeFixedLength(hexValue.substring(2, 10));
    return [decValue, skip];
  } else if (firstByte === 255) {
    skip = 16
    const [decValue, hexSkip] = decodeFixedLength(hexValue.substring(2, 18));
    return [decValue, skip];

  }
}

const decodeField = (startIndex, length, decodeFunction, rawTransaction) => {
  const hexValue = rawTransaction.substring(startIndex, startIndex + length);
  const decodedValue = decodeFunction(hexValue);
  return decodedValue;
}

const getVersion = (rawTransaction) => {
  const version = decodeField(0, 8, decodeFixedLength, rawTransaction);
  return version; // we get 2 as the version 
}



const isSegWitTransaction = (rawTransaction) => {
  if (rawTransaction.length < 10) {
    throw new Error("Raw transaction is too short");
  }

  const version = parseInt(rawTransaction.substring(0, 8), 16);
  if (version < 2) {
    return false;
  }

  // Check if the byte after version is 0 and the next byte is 1
  const marker = parseInt(rawTransaction.substring(8, 10), 16);
  const flag = parseInt(rawTransaction.substring(10, 12), 16);
  return marker === 0 && flag === 1;
}


const getInputCount = (rawTransaction, isSegWit) => {
  let [, nextOffset] = getVersion(rawTransaction);

  // skip marker and flag (2 bytes or 4 hex) if it is segwit transaction
  nextOffset = isSegWit ? nextOffset + 4 : nextOffset;
  const [inputCount, hexToSkip] = decodeField(nextOffset, 2, decodeCompactSize, rawTransaction);

  return [inputCount, nextOffset + hexToSkip];
}

const getInputs = (rawTransaction, isSegWit) => {
  let inputs = [];
  let [inputCount, nextOffset] = getInputCount(rawTransaction, isSegWit);

  currentOffset = nextOffset
  for (let i = 0; i < inputCount; i++) {
    const [txid, hexToSkip] = decodeField(currentOffset, 64, convertEndian, rawTransaction);
    currentOffset += hexToSkip;

    const [vout, hexToSkip1] = decodeField(currentOffset, 8, decodeFixedLength, rawTransaction);
    currentOffset += hexToSkip1;

    const scriptSigLenHex = rawTransaction.substring(currentOffset, currentOffset + 2);
    const [scriptSigLen, hexToSkip2] = decodeCompactSize(scriptSigLenHex);
    currentOffset += hexToSkip2;

    const scriptSig = rawTransaction.substring(currentOffset, currentOffset + scriptSigLen * 2);
    currentOffset += scriptSigLen * 2;

    const [sequence, hexToSkip3] = decodeField(currentOffset, 8, decodeFixedLength, rawTransaction);
    currentOffset += hexToSkip3;

    const input = {
      txid: txid,
      vout: vout,
      scriptSig: scriptSig,
      txinwitness: [], // will be field when outputs are extracted
      sequence: sequence,
    };

    inputs.push(input);
  }

  return [inputs, currentOffset];
}

const getOutputCount = (rawTransaction, isSegWit) => {
  const [, nextOffset] = getVersion(rawTransaction, isSegWit);

  const [inputCount, hexToSkip] = decodeField(nextOffset, 2, decodeCompactSize, rawTransaction);

  return [inputCount, nextOffset + hexToSkip];
}

const getOutputs = (rawTransaction, isSegWit) => {
  let outputs = [];
  const [outputCount, nextOffset] = getOutputCount(rawTransaction, isSegWit);
  currentOffset = nextOffset;

  for (let i = 0; i < outputCount; i++) {
    const outputAmountHex = rawTransaction.substring(currentOffset, currentOffset + 16);
    const outputAmountBigEndian = convertEndian(outputAmountHex);

    const outputAmountSatoshi = parseInt(outputAmountBigEndian, 16);
    const outputAmountBitcoin = satoshisToBitcoin(outputAmountSatoshi);
    currentOffset += 16;

    const outputScriptLenHex = rawTransaction.substring(currentOffset, currentOffset + 2);
    const [outputScriptLen, hexToSkip] = decodeCompactSize(outputScriptLenHex);
    currentOffset += hexToSkip;

    const scriptPubKeyHex = rawTransaction.substring(currentOffset, currentOffset + outputScriptLen * 2);
    currentOffset += outputScriptLen * 2;

    const output = {
      value: outputAmountBitcoin,
      n: i,
      scriptPubKey: {
        hex: scriptPubKeyHex
      }
    };

    outputs.push(output);
  }

  return [outputs, currentOffset];
}

const getWitnessItems = (rawTransaction) => {
  const isSegWit = true;
  const [inputCount, _] = getInputCount(rawTransaction, isSegWit);
  const [, nextOffset] = getOutputs(rawTransaction, isSegWit);
  let currentOffset = nextOffset;

  const witnessItemsForAllInputs = [];
  for (let i = 0; i < inputCount; i++) {
    const witnessItemsCountHex = rawTransaction.substring(currentOffset, currentOffset + 2);
    const witnessItemsCount = parseInt(witnessItemsCountHex, 16);
    currentOffset += 2;

    const witnessItemsForOneInput = [];
    for (let j = 0; j < witnessItemsCount; j++) {
      const WitnessItemSizeHex = rawTransaction.substring(currentOffset, currentOffset + 2);
      const [WitnessItemSize, hexToSkip] = decodeCompactSize(WitnessItemSizeHex);
      currentOffset += hexToSkip;

      const WitnessItemHex = rawTransaction.substring(currentOffset, currentOffset + WitnessItemSize * 2);
      currentOffset += WitnessItemSize * 2;

      witnessItemsForOneInput[j] = WitnessItemHex
    }
    witnessItemsForAllInputs.push(witnessItemsForOneInput);
  }

  return [witnessItemsForAllInputs, currentOffset];
}

const getLockTime = (rawTransaction) => {
  const lockTimeHex = rawTransaction.substring(rawTransaction.length - 8);
  const [lockTime,] = decodeFixedLength(lockTimeHex);
  return lockTime;
}







module.exports = {
  convertEndian,
  decodeFixedLength,
  decodeCompactSize,
  decodeField,
  isSegWitTransaction,
  getInputCount,
  getInputs,
  getOutputCount,
  getOutputs,
  getWitnessItems,
  getLockTime,
  getVersion
}

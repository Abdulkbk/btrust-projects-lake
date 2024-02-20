## Bitcoin Transaction Decoder from Scratch
This repository contains a Go program that decodes Bitcoin raw transactions without using any library. It provides a way to extract information about transactions including inputs, outputs, values, version, and more.

### Features
- Transaction Decoding: Decode raw Bitcoin transactions into a human-readable format.
- The program decodes Segwit transactions

### Getting Started
To use the Bitcoin CLI Transaction Decoder, follow these steps:

- Clone the repository to your local machine.
package.
- Run the go run main.go command, providing the Bitcoin raw transaction hex as an argument in the func main.
- The program will decode the transaction and display the results in the terminal.

```bash
    go run main.go
```
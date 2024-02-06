## Welcome to Week 2 deliverable: Developing a transaction `HEX` Decoding program

### ABOUT
This is a CLI application that decode Bitcoin transaction `hex`. It is implemented using Golang. The main package used is the `btcd`. 

### How to use the Program
After clonning the repo and downloading the necessary dependencies, proceed and run the program like so:

```bash
go run . -v -in -out -lock-time -hex <replace with actual transaction hex>
```

#### Parameters

| Params | Description | Required |
|--------|-------------|----------|
| -v     | Print transaction version | false |
| -in    | Print transaction inputs | false |
| -out   | Print transaction outputs | false |
| -lock-time | Print transaction locktime | false |
| -hex   | Specify the transaction Hex | **true** |
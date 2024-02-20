## Transaction Builder - week 2 deliverable of the Btrust Bitcoin Fellowship

### ABOUT
The transaction builder is implemented using golang. It takes a byte encoding of the string `Btrust Builders` and generate a redeem script which can be used with a `p2sh` transaction.

### Functions
- `generateRedeemScript`: Takes in the preimage, hash it to get a lockhex, and then combine the lockhex with the required `opcodes` for spending the transaction.

- `generateAddressFromRedeemScript`: This function generates a `p2sh` address from the redeem script generated earelier, which can be used to accept bitcoin payments.

- `transactionContructor`: This takes reference to utxo, amount, and other necessary parameters to construct a transaction.

- `transactionSpenderConstructor`: This function takes some parameters and construct a transaction based on the previous transaction gotten from `transactionContructor` function.

### Installation
- Install Go on your machine.
- Clone this repository
- Navigate to the repository: cd transaction-builder
- Run go get to install the required dependencies.
- Run the application: go run main.go
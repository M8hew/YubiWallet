# YubiWallet

YubiWallet is a command-line tool that allows you to interact with a YubiKey for Bitcoin transaction signing. The tool can be used for initializing the YubiKey and constructing, signing, and marshaling Bitcoin transactions.

## Features
* Initialize YubiKey: Generate keys and store them on your YubiKey.
* Construct Bitcoin Transactions: Create and construct a Bitcoin transaction.
* Sign Bitcoin Transactions: Sign transactions using the private key stored on the YubiKey.
* Broadcast Bitcoin Transactions: Serialize the transaction to hexadecimal format, ready to be broadcasted to the Bitcoin network.

## Installation
Clone the repository:

```bash
git clone https://github.com/yourusername/yubiwallet.git
cd yubiwallet
```

Install dependencies:
```bash
go mod tidy
```

## Usage
### Command-Line Flags
The tool uses flags for different operations. Here's an overview of the available flags:

| Flag                  |	Description 
| ---                   |   --- 
| `-init`               |	Initialize the YubiKey by generating and storing keys. 
| `-tx`	                |   Create, sign, and broadcast a Bitcoin transaction. 
| `-utxo_hash`          |	The hash of the UTXO to spend. Required when constructing a transaction.
| `-utxo_index`         |	The index of the UTXO to spend. Required when constructing a transaction.
| `-recipient`          | addr	The recipient Bitcoin address. Required when constructing a transaction.
| `-amount`	            | The amount of Bitcoin (in satoshis) to send. Required when constructing a transaction.
| `-transaction_fee`    |	The fee for the transaction in satoshis. Required when constructing a transaction.

### Example Commands
Initialize the YubiKey (Generates and stores the key on the YubiKey):

```bash
go run main.go -init
Create and Sign a Bitcoin Transaction:
```

```bash
go run main.go -tx -utxo_hash <UTXO_HASH> -utxo_index <UTXO_INDEX> -recipient <RECIPIENT_ADDRESS> -amount <AMOUNT_IN_SATOSHIS> -transaction_fee <FEE_IN_SATOSHIS>
```

Example:

```bash
go run main.go -tx -utxo_hash "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef" -utxo_index 0 -recipient "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa" -amount 100000 -transaction_fee 5000
```

This will:

* Construct a Bitcoin transaction.
* Sign it using the private key stored on the YubiKey.
* Attach the signature to the transaction.
* Prepare the transaction for broadcasting.
* Serialize the Transaction: After signing, the transaction will be serialized into a hex format, which you can use to broadcast it via Bitcoin Core or other broadcasting services.

## Notes
Make sure you have the YubiKey properly set up and connected to your system before using the `-init` or `-tx` commands.

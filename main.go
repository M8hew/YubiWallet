package main

import (
	"flag"
	"fmt"

	"yubiwallet/internal/transaction"
	"yubiwallet/internal/yubi"
)

type Config struct {
	init bool
	tx   bool

	utxoHash  string
	utxoIndex uint64
	recipient string
	amount    int64
	fee       int64
}

var cfg Config

func init() {
	flag.BoolVar(&cfg.init, "init", false, "init yubikey by generating keys")
	flag.BoolVar(&cfg.tx, "tx", false, "init yubikey by generating keys")
	flag.StringVar(&cfg.utxoHash, "utxo_hash", "", "utxo hash for transaction")
	flag.Uint64Var(&cfg.utxoIndex, "utxo_index", 0, "utxo index for transaction")
	flag.StringVar(&cfg.recipient, "recipient addr", "", "recipient addr for transaction")
	flag.Int64Var(&cfg.amount, "amount", 0, "satoshi amount for transaction")
	flag.Int64Var(&cfg.fee, "transaction fee", 0, "fee amount for transaction")
}

func main() {
	flag.Parse()

	switch {
	case cfg.init:
		yubi.InitKey()
		return

	case cfg.tx:
		tx, txHash, err := transaction.Construct(
			cfg.utxoHash,
			uint32(cfg.utxoIndex),
			cfg.recipient,
			cfg.amount,
			cfg.fee,
		)
		if err != nil {
			fmt.Println("Failed to construct transaction: %w", err)
			return
		}

		sig, err := transaction.Sign(txHash)
		if err != nil {
			fmt.Println("Failed to sign transaction: %w", err)
			return
		}

		pubKey, err := yubi.GetPubKey()
		if err != nil {
			fmt.Println("Failed to get public key: %w", err)
			return
		}

		if err := transaction.AttachSignature(tx, sig, pubKey); err != nil {
			fmt.Println("Failed to attach signature: %w", err)
			return
		}

		strTx, err := transaction.PrepareBroadcast(tx)
		if err != nil {
			fmt.Println("Failed to marshal transaction: %w", err)
			return
		}

		fmt.Println("Transaction ready to broadcast: %w", strTx)
		return

	default:
		fmt.Println("No option choosen")
	}
}

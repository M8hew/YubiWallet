package transaction

import (
	"bytes"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"

	"yubiwallet/internal/yubi"
)

// Construct creates an unsigned Bitcoin transaction.
func Construct(utxoHash string, utxoIndex uint32, recipientAddr string, amount, fee int64) (*wire.MsgTx, []byte, error) {
	tx := wire.NewMsgTx(wire.TxVersion)

	utxoHashBytes, err := hex.DecodeString(utxoHash)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid UTXO hash: %w", err)
	}

	hashedUTXO, err := chainhash.NewHash(utxoHashBytes)
	if err != nil {
		return nil, nil, fmt.Errorf("error hashing utxo: %w", err)
	}

	utxo := wire.NewOutPoint(hashedUTXO, utxoIndex)
	txIn := wire.NewTxIn(utxo, nil, nil)
	tx.AddTxIn(txIn)

	if amount <= fee {
		return nil, nil, fmt.Errorf("fee exceeds amount")
	}
	outputAmount := amount - fee

	// Create the output for the recipient
	addr, err := btcutil.DecodeAddress(recipientAddr, &chaincfg.MainNetParams)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid recipient address: %w", err)
	}

	script, err := txscript.PayToAddrScript(addr)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create script: %w", err)
	}

	txOut := wire.NewTxOut(outputAmount, script)
	tx.AddTxOut(txOut)

	// Serialize the transaction and calculate the hash for signing
	var buf bytes.Buffer
	if err := tx.Serialize(&buf); err != nil {
		return nil, nil, fmt.Errorf("failed to serialize transaction: %w", err)
	}
	txHash := sha256.Sum256(buf.Bytes())

	return tx, txHash[:], nil
}

// Sign transaction with YubiKey.
func Sign(hash []byte) ([]byte, error) {
	sig, err := yubi.SignWithKey(hash)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}
	return sig, nil
}

// Attach signature to transaction
func AttachSignature(tx *wire.MsgTx, sig []byte, pubKey *rsa.PublicKey) error {
	pubKeyDer, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return fmt.Errorf("failed to marshal RSA public key: %w", err)
	}

	pubKeyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubKeyDer,
	})

	sigScipt, err := txscript.NewScriptBuilder().
		AddData(sig).
		AddData(pubKeyPem).
		Script()
	if err != nil {
		return fmt.Errorf("failed to build script: %w", err)
	}

	if len(tx.TxIn) == 0 {
		return errors.New("transaction has no inputs to attach")
	}
	tx.TxIn[0].SignatureScript = sigScipt
	return nil
}

func PrepareBroadcast(tx *wire.MsgTx) (string, error) {
	var buf bytes.Buffer
	if err := tx.Serialize(&buf); err != nil {
		return "", fmt.Errorf("failed to serialize transaction %w", err)
	}
	return hex.EncodeToString(buf.Bytes()), nil
}

package yubi

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/asn1"
	"errors"
	"fmt"
	"strings"

	"github.com/go-piv/piv-go/v2/piv"
)

func getKeyObj() (*piv.YubiKey, error) {
	cards, err := piv.Cards()
	if err != nil {
		return nil, fmt.Errorf("couldn't find any connected yubikeys: %w", err)
	}

	var yk *piv.YubiKey
	for _, card := range cards {
		if strings.Contains(strings.ToLower(card), "yubikey") {
			if yk, err = piv.Open(card); err != nil {
				return yk, nil
			}
			break
		}
	}

	return nil, errors.New("couldn't open any yubikeys")
}

func InitKey() error {
	yk, err := getKeyObj()
	if err != nil {
		return fmt.Errorf("error getting key obj: %w", err)
	}
	defer yk.Close()

	key := piv.Key{
		Algorithm:   piv.AlgorithmRSA2048,
		PINPolicy:   piv.PINPolicyAlways,
		TouchPolicy: piv.TouchPolicyAlways,
	}

	pub, err := yk.GenerateKey(piv.DefaultManagementKey, piv.SlotAuthentication, key)
	if err != nil {
		return fmt.Errorf("failed to generate key: %w", err)
	}

	fmt.Printf("Key generated successfully. Public Key: %x\n", pub)
	return nil
}

func GetPubKey() (*rsa.PublicKey, error) {
	yk, err := getKeyObj()
	if err != nil {
		return nil, fmt.Errorf("error getting key obj: %w", err)
	}
	defer yk.Close()

	info, err := yk.KeyInfo(piv.SlotAuthentication)
	if err != nil {
		return nil, fmt.Errorf("error extracting public key: %w", err)
	}

	pubKey, ok := info.PublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("error processing public key")
	}
	return pubKey, nil
}

func SignWithKey(message []byte) ([]byte, error) {
	yk, err := getKeyObj()
	if err != nil {
		return nil, fmt.Errorf("error getting key obj: %w", err)
	}
	defer yk.Close()

	info, err := yk.KeyInfo(piv.SlotAuthentication)
	if err != nil {
		return nil, fmt.Errorf("error extracting public key: %w", err)
	}

	privKey, err := yk.PrivateKey(piv.SlotAuthentication, info.PublicKey, piv.KeyAuth{PIN: piv.DefaultPIN})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve private key: %w", err)
	}

	rsaPrivKey, ok := privKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("retrieved key is not an RSA private key")
	}

	hashed := sha256.Sum256(message)

	// Sign the hashed message
	signedBytes, err := rsa.SignPSS(rand.Reader, rsaPrivKey, crypto.SHA256, hashed[:], nil)
	if err != nil {
		return nil, fmt.Errorf("failed to sign message: %w", err)
	}

	// Serialize the signature as ASN.1(Bitcoin format)
	sig, err := asn1.Marshal(signedBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize signature: %w", err)
	}
	return sig, nil
}

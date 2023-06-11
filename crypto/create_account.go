package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/crypto"
)

func CreateAddress() (string, string, error) {
	// Generate a new private key
	privateKey, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	if err != nil {
		return "", "", err
	}

	// Get the public key from the private key
	publicKey := privateKey.Public()

	// Convert the public key to an Ethereum address
	address := crypto.PubkeyToAddress(*publicKey.(*ecdsa.PublicKey)).Hex()

	// Convert the private key to a string
	privateKeyString := hex.EncodeToString(privateKey.D.Bytes())
	return address, privateKeyString, nil
}

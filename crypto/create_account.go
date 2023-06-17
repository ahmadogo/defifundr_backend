package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
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

func GenerateAccountKeyStone(password string) (string, string, error) {
	// Generate a new random private key
	key := keystore.NewKeyStore("./../tmp", keystore.StandardScryptN, keystore.StandardScryptP)

	// Create a new account with the specified encryption passphrase
	passwordKey := password
	account, err := key.NewAccount(passwordKey)

	if err != nil {
		log.Error().Err(err).Msg("cannot create account")
	}

	filename := account.URL.Path[strings.LastIndex(account.URL.Path, "/")+1:]

	accountName := account.Address.Hex()

	return filename, accountName, nil
}

func DecryptPrivateKey(path string, passphrase string) (*ecdsa.PrivateKey, string, error) {
	filePath := fmt.Sprintf("./../tmp/%s", path)

	b, err := os.ReadFile(filePath)
	if err != nil {
		return nil, "", err
	}

	key, err := keystore.DecryptKey(b, passphrase)
	if err != nil {
		return nil, "", err
	}

	pData := crypto.FromECDSA(key.PrivateKey)
	privateKey := hexutil.Encode(pData)
	publicKey, err := GeneratePublicKeyFromPrivateKey(privateKey)
	if err != nil {
		return nil, "", err
	}

	return key.PrivateKey, publicKey, nil
}

func GeneratePublicKeyFromPrivateKey(privateKey string) (string, error) {
	pData, err := hexutil.Decode(privateKey)
	if err != nil {
		return "", err
	}

	privateKeyECDSA, err := crypto.ToECDSA(pData)
	if err != nil {
		return "", err
	}

	publicKey := crypto.PubkeyToAddress(privateKeyECDSA.PublicKey).Hex()

	return publicKey, nil
}

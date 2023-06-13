package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncode(t *testing.T) {
	privateKey, publicKey := GenerateKeyPair()
	pemEncoded, pemEncodedPub, err := Encode(privateKey, publicKey)
	require.NoError(t, err)
	require.NotEmpty(t, pemEncoded)
	require.NotEmpty(t, pemEncodedPub)

	privateKey2, publicKey2, err := Decode(pemEncoded, pemEncodedPub)
	require.Equal(t, privateKey, privateKey2)
	require.Equal(t, publicKey, publicKey2)
	require.NoError(t, err)
}

func GenerateKeyPair() (*ecdsa.PrivateKey, *ecdsa.PublicKey) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	publicKey := &privateKey.PublicKey

	return privateKey, publicKey
}

func TestInvalidEncode(t *testing.T) {
	_, _, err := Encode(nil, nil)
	require.Error(t, err)
}

func TestInvalidDecode(t *testing.T) {
	_, _, err := Decode("", "")
	require.Error(t, err)
}

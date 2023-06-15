package utils

import (
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

func TestInvalidEncode(t *testing.T) {
	_, _, err := Encode(nil, nil)
	require.Error(t, err)
}

func TestInvalidDecode(t *testing.T) {
	_, _, err := Decode("", "")
	require.Error(t, err)
}

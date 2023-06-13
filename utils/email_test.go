package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSendEmail(t *testing.T) {
	emailAddr := RandomEmail()
	details := "This is a test email"

	result, err := SendEmail(emailAddr, RandomOwner(), details)
	require.NoError(t, err)
	require.NotEmpty(t, result)
}

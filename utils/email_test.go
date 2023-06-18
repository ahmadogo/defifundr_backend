package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSendEmail(t *testing.T) {
	emailAddr := RandomEmail()
	details := "This is a test email"

	info := EmailInfo{
		Name:    "Test",
		Otp:     RandomOtp(),
		Details: details,
		Subject: "Test Email",
	}

	result, err := SendEmail(emailAddr, RandomOwner(), info, ".")
	require.NoError(t, err)
	require.NotEmpty(t, result)
}

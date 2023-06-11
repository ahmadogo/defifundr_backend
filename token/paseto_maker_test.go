package token

import (
	"testing"
	"time"

	"github.com/demola234/defiraise/utils"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewTokenMaker(utils.RandomString(32))
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	username := utils.RandomString(6)
	duration := time.Minute
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiresAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewTokenMaker(utils.RandomString(32))
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	username := utils.RandomString(6)
	duration := -time.Minute

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Empty(t, payload)
}

func TestInvalidToken(t *testing.T) {
	maker, err := NewTokenMaker(utils.RandomString(32))
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	payload, err := maker.VerifyToken("invalid_token")
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Empty(t, payload)
}

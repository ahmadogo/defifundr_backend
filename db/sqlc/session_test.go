package db

import (
	"context"
	"testing"
	"time"

	"github.com/demola234/defiraise/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func createRandomSession(t *testing.T) UserSession {

	expired_at := time.Now().Add(time.Hour * 24 * 7)

	arg := CreateSessionParams{
		Username:     "test",
		ID:           uuid.New(),
		RefreshToken: utils.RandomString(6),
		UserAgent:    utils.RandomString(6),
		ClientIp:     utils.RandomString(6),
		IsBlocked:    false,
		ExpiresAt:    expired_at,
	}
	session, err := testQueries.CreateSession(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, session)

	require.Equal(t, arg.Username, session.Username)

	require.NotZero(t, session.Username)
	require.NotZero(t, session.CreatedAt)

	return session
}
func TestCreateSession(t *testing.T) {
	createRandomSession(t)
}

func TestGetSession(t *testing.T) {
	session := createRandomSession(t)

	session2, err := testQueries.GetSession(context.Background(), session.ID)
	require.NoError(t, err)
	require.NotEmpty(t, session2)

	require.Equal(t, session.Username, session2.Username)

	require.NotZero(t, session2.Username)
	require.NotZero(t, session2.CreatedAt)
}

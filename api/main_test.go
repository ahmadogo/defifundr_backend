package api

import (
	"os"
	"testing"
	"time"

	db "github.com/demola234/defiraise/db/sqlc"
	"github.com/demola234/defiraise/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := utils.Config{
		TokenSymmetricKey:   utils.RandomString(32),
		AccessTokenDuration: time.Minute,
	}
	newServer, err := NewServer(config, store)
	require.NoError(t, err)
	require.NotNil(t, newServer)

	return newServer
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

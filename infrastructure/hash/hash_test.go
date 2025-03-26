package commons

import (
	"strings"
	"testing"

	"github.com/demola234/defifundr/pkg/random"
	"github.com/stretchr/testify/require"
)

func TestPasswordHashing(t *testing.T) {
	t.Parallel()

	password := random.RandomString(12)

	t.Run("successful hash and verify", func(t *testing.T) {
		t.Parallel()
		hashPassword, err := HashPassword(password)
		require.NoError(t, err)
		require.NotEmpty(t, hashPassword)
		require.True(t, strings.HasPrefix(hashPassword, "$argon2id$"))

		match, err := CheckPassword(password, hashPassword)
		require.NoError(t, err)
		require.True(t, match)
	})

	t.Run("wrong password fails", func(t *testing.T) {
		t.Parallel()
		hashPassword, err := HashPassword(password)
		require.NoError(t, err)

		wrongPassword := random.RandomString(12)
		match, err := CheckPassword(wrongPassword, hashPassword)
		require.NoError(t, err)
		require.False(t, match)
	})

	t.Run("different salts produce different hashes", func(t *testing.T) {
		t.Parallel()
		hash1, err := HashPassword(password)
		require.NoError(t, err)

		hash2, err := HashPassword(password)
		require.NoError(t, err)

		require.NotEqual(t, hash1, hash2)
	})

	t.Run("empty password fails", func(t *testing.T) {
		t.Parallel()
		_, err := HashPassword("")
		require.Error(t, err)
	})

	t.Run("malformed hash fails verification", func(t *testing.T) {
		t.Parallel()
		_, err := CheckPassword(password, "$invalid$hash$format")
		require.Error(t, err)
	})

	t.Run("password too long fails", func(t *testing.T) {
		t.Parallel()
		longPassword := random.RandomString(100)
		_, err := HashPassword(longPassword)
		require.Error(t, err)
	})
}

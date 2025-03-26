package commons

import (
	"os"
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

func TestArgon2Parameters(t *testing.T) {
	t.Parallel()

	t.Run("custom parameters via env vars", func(t *testing.T) {
		t.Parallel()
		os.Setenv("ARGON2_MEMORY", "128000") // 128MB
		os.Setenv("ARGON2_ITERATIONS", "5")
		os.Setenv("ARGON2_PARALLELISM", "4")
		defer os.Unsetenv("ARGON2_MEMORY")
		defer os.Unsetenv("ARGON2_ITERATIONS")
		defer os.Unsetenv("ARGON2_PARALLELISM")

		password := random.RandomString(12)
		hashStr, err := HashPassword(password)
		require.NoError(t, err)

		// Verify parameters were applied
		parts := strings.Split(hashStr, "$")
		require.Equal(t, "m=128000,t=5,p=4", parts[3])
	})
}

func TestHashFormat(t *testing.T) {
	t.Parallel()

	password := random.RandomString(12)
	hashStr, err := HashPassword(password)
	require.NoError(t, err)

	parts := strings.Split(hashStr, "$")
	require.Len(t, parts, 6)
	require.Equal(t, "argon2id", parts[1])
	require.True(t, strings.HasPrefix(parts[2], "v=19"))
	require.Regexp(t, `^m=\d+,t=\d+,p=\d+$`, parts[3])
	require.NotEmpty(t, parts[4]) // salt
	require.NotEmpty(t, parts[5]) // hash
}

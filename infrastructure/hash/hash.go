package commons

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"unsafe"

	"golang.org/x/crypto/argon2"
)

const (
	maxPasswordLength = 72 
)

type Argon2Params struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

var DefaultParams = Argon2Params{
	Memory:      64 * 1024, // 64MB
	Iterations:  3,
	Parallelism: 2,
	SaltLength:  16,
	KeyLength:   32,
}

func init() {
	// Securely wipe memory when params are garbage collected
	runtime.SetFinalizer(&DefaultParams, secureWipeParams)
}

func secureWipeParams(params *Argon2Params) {
	// Wipe memory that contained sensitive data
	mem := unsafe.Pointer(params)
	size := unsafe.Sizeof(*params)
	for i := uintptr(0); i < size; i++ {
		*(*byte)(unsafe.Pointer(uintptr(mem) + i)) = 0
	}
}

func LoadParams() Argon2Params {
	params := DefaultParams

	if memStr := os.Getenv("ARGON2_MEMORY"); memStr != "" {
		if mem, err := strconv.ParseUint(memStr, 10, 32); err == nil {
			params.Memory = uint32(mem)
		}
	}

	if iterationsStr := os.Getenv("ARGON2_ITERATIONS"); iterationsStr != "" {
		if iterations, err := strconv.ParseUint(iterationsStr, 10, 32); err == nil {
			params.Iterations = uint32(iterations)
		}
	}

	if parallelismStr := os.Getenv("ARGON2_PARALLELISM"); parallelismStr != "" {
		if parallelism, err := strconv.ParseUint(parallelismStr, 10, 8); err == nil {
			params.Parallelism = uint8(parallelism)
		}
	}

	return params
}

func HashPassword(password string) (string, error) {

	//check if password is empty
	if password == "" {
		return "", errors.New("password cannot be empty")
	}

	// Security check: prevent DoS through extremely long passwords
	if len(password) > maxPasswordLength {
		return "", errors.New("password length exceeds maximum allowed")
	}

	params := LoadParams()

	// Generate cryptographically secure random salt
	salt := make([]byte, params.SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// Generate hash using Argon2id
	hash := argon2.IDKey(
		[]byte(password),
		salt,
		params.Iterations,
		params.Memory,
		params.Parallelism,
		params.KeyLength,
	)

	// Encode components
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Format: $argon2id$v=19$m=65536,t=3,p=2$salt$hash
	encoded := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		params.Memory,
		params.Iterations,
		params.Parallelism,
		b64Salt,
		b64Hash,
	)

	// Securely wipe temporary buffers
	defer func() {
		for i := range hash {
			hash[i] = 0
		}
		for i := range salt {
			salt[i] = 0
		}
	}()

	return encoded, nil
}

func CheckPassword(password, encodedHash string) (bool, error) {
	// Security check: prevent DoS through extremely long passwords
	if len(password) > maxPasswordLength {
		return false, errors.New("password length exceeds maximum allowed")
	}

	// Parse the encoded hash
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false, errors.New("invalid hash format")
	}

	// Check algorithm identifier
	if parts[1] != "argon2id" {
		return false, errors.New("unsupported algorithm")
	}

	// Parse version
	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return false, fmt.Errorf("invalid version format: %w", err)
	}
	if version != argon2.Version {
		return false, errors.New("incompatible argon2 version")
	}

	// Parse parameters
	params := Argon2Params{}
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d",
		&params.Memory,
		&params.Iterations,
		&params.Parallelism); err != nil {
		return false, fmt.Errorf("invalid parameter format: %w", err)
	}

	// Decode salt and hash
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, fmt.Errorf("invalid salt encoding: %w", err)
	}
	params.SaltLength = uint32(len(salt))

	storedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, fmt.Errorf("invalid hash encoding: %w", err)
	}
	params.KeyLength = uint32(len(storedHash))

	// Recompute the hash
	computedHash := argon2.IDKey(
		[]byte(password),
		salt,
		params.Iterations,
		params.Memory,
		params.Parallelism,
		params.KeyLength,
	)

	// Securely wipe temporary buffers
	defer func() {
		for i := range computedHash {
			computedHash[i] = 0
		}
		for i := range salt {
			salt[i] = 0
		}
	}()

	// Constant-time comparison
	if subtle.ConstantTimeCompare(computedHash, storedHash) == 1 {
		return true, nil
	}
	return false, nil
}

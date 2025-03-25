package commons

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"unsafe"

	"golang.org/x/crypto/argon2"
)

const (
	maxPasswordLength = 72 // Argon2's maximum input length
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

package passwd

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Argon2 parameters definition
const (
	time    = 1         // Computational cost (number of iterations)
	memory  = 64 * 1024 // Memory usage (in KB)
	threads = 4         // Number of parallel threads
	keyLen  = 32        // Length of the derived key (in bytes)
)

// New generates a hashed password using Argon2 and a randomly generated salt.
// Input: password `p` and size of the salt `saltSize`
// Output: hashed password (salt and hash encoded as a string), or an error
func New(p string, saltSize int) (string, error) {
	// Generate a random salt
	salt, err := Salt(saltSize)
	if err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// Hash the password with the generated salt
	hashed, err := Hash(p, salt)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return hashed, nil
}

func MustNew(p string, saltSize int) string {
	hashed, err := New(p, saltSize)
	if err != nil {
		panic(err)
	}
	return hashed
}

// Salt generates a random salt of the specified size.
// Input: size of the salt `size`
// Output: generated random salt, or an error
func Salt(size int) ([]byte, error) {
	salt := make([]byte, size)
	// Generate random bytes
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}
	return salt, nil
}

// Hash creates a hashed password using Argon2 with the given password and salt.
// Input: password `p` and salt `salt`
// Output: a string containing the base64-encoded salt and hash, or an error
func Hash(p string, salt []byte) (string, error) {
	// Generate the Argon2 key (hash)
	hash := argon2.IDKey([]byte(p), salt, time, memory, threads, keyLen)

	// Encode the salt and hash in base64
	saltEncoded := base64.RawStdEncoding.EncodeToString(salt)
	hashEncoded := base64.RawStdEncoding.EncodeToString(hash)

	// Return the concatenated salt and hash
	return fmt.Sprintf("%s:%s", saltEncoded, hashEncoded), nil
}

// Verify checks if the provided password matches the encoded hash.
// Input: password `p` and encoded hash `encodedHash` (in "salt:hash" format)
// Output: true if the password matches, or false and an error if it doesn't
func Verify(p, encodedHash string) (bool, error) {
	// Split the encoded hash into salt and hash parts
	parts := strings.Split(encodedHash, ":")
	if len(parts) != 2 {
		return false, fmt.Errorf("invalid hash format")
	}

	// Decode the salt from base64
	salt, err := base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		return false, fmt.Errorf("failed to decode salt: %w", err)
	}

	// Decode the hash from base64
	hash, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return false, fmt.Errorf("failed to decode hash: %w", err)
	}

	// Recompute the hash with the provided password and salt
	expectedHash := argon2.IDKey([]byte(p), salt, time, memory, threads, keyLen)

	// Compare the hashes and return the result
	return string(expectedHash) == string(hash), nil
}

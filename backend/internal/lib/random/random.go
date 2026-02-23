package random

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
)

// NewString generates a new random string with the specified byte length.
// The returned string is the hexadecimal representation of the random bytes,
// so its length will be twice the input length.
func NewString(length int) (string, error) {
	bytes, err := NewBytes(length)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// NewBytes generates a new random byte slice with the specified length.
func NewBytes(length int) ([]byte, error) {
	if length <= 0 {
		return nil, errors.New("length must be a positive integer")
	}
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return nil, fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return b, nil
}

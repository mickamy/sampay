package random

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"
)

// NewString generates a new random string with the specified byte length.
// The returned string is the hexadecimal representation of the random bytes,
// so its length will be twice the input length.
func NewString(length int) (string, error) {
	bytes, err := NewBytes(length)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", bytes), nil
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

// NewPinCode generates a random numeric pin code of the specified length.
func NewPinCode(length int) (string, error) {
	if length <= 0 {
		return "", errors.New("length must be a positive integer")
	}

	var builder strings.Builder
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", fmt.Errorf("failed to generate random pin code: %w", err)
		}
		builder.WriteByte('0' + byte(n.Int64()))
	}
	return builder.String(), nil
}

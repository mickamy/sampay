package random_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mickamy.com/sampay/internal/lib/random"
)

func TestNewString(t *testing.T) {
	t.Parallel()

	// arrange
	length := 32

	// act
	got, err := random.NewString(length)

	// assert
	require.NoError(t, err)
	assert.Len(t, got, length*2)
}

func TestNewBytes(t *testing.T) {
	t.Parallel()

	// arrange
	length := 32

	// act
	got, err := random.NewBytes(length)

	// assert
	require.NoError(t, err)
	assert.Len(t, got, length)
}

func TestNewPinCode(t *testing.T) {
	t.Parallel()

	// act
	got, err := random.NewPinCode(6)

	// assert
	require.NoError(t, err)
	assert.Len(t, got, 6)
}

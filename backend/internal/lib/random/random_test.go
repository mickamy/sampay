package random_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"mickamy.com/sampay/internal/lib/random"
)

func Test_NewString(t *testing.T) {
	t.Parallel()

	// arrange
	length := 32

	// act
	got, err := random.NewString(length)

	// assert
	assert.NoError(t, err)
	assert.Len(t, got, length*2)
}

func Test_NewBytes(t *testing.T) {
	t.Parallel()

	// arrange
	length := 32

	// act
	got, err := random.NewBytes(length)

	// assert
	assert.NoError(t, err)
	assert.Len(t, got, length)
}

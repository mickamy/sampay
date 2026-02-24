package ptr_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mickamy/sampay/internal/lib/ptr"
)

func TestMap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		val    *int
		mapper func(*int) string
		want   string
	}{
		{
			name: "val is nil",
			val:  nil,
			mapper: func(*int) string {
				return "nil"
			},
			want: "",
		},
		{
			name: "val is not nil",
			val:  new(int),
			mapper: func(*int) string {
				return "not nil"
			},
			want: "not nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := ptr.Map(tt.val, tt.mapper)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNullIfZero(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		val  int
		want *int
	}{
		{
			name: "val is empty",
			val:  0,
			want: nil,
		},
		{
			name: "val is not nil",
			val:  1,
			want: ptr.Of(1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := ptr.NullIfZero(tt.val)

			assert.Equal(t, tt.want, got)
		})
	}
}

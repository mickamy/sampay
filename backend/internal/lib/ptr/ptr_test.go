package ptr_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mickamy/sampay/internal/lib/ptr"
)

func TestMap(t *testing.T) {
	t.Parallel()

	tcs := []struct {
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

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := ptr.Map(tc.val, tc.mapper)

			assert.Equal(t, tc.want, got)
		})
	}
}

func TestNullIfZero(t *testing.T) {
	t.Parallel()

	tcs := []struct {
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

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := ptr.NullIfZero(tc.val)

			assert.Equal(t, tc.want, got)
		})
	}
}

package ptr_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"mickamy.com/sampay/internal/lib/ptr"
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

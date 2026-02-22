package database_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mickamy/sampay/internal/infra/storage/database"
)

func TestStringArray_Scan(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   any
		want    database.StringArray
		wantErr bool
	}{
		{
			name:  "nil",
			input: nil,
			want:  nil,
		},
		{
			name:  "empty array",
			input: "{}",
			want:  []string{},
		},
		{
			name:  "single element",
			input: "{hello}",
			want:  []string{"hello"},
		},
		{
			name:  "multiple elements",
			input: "{a,b,c}",
			want:  []string{"a", "b", "c"},
		},
		{
			name:  "quoted elements",
			input: `{"hello world","foo bar"}`,
			want:  []string{"hello world", "foo bar"},
		},
		{
			name:  "escaped quotes",
			input: `{"say \"hi\""}`,
			want:  []string{`say "hi"`},
		},
		{
			name:    "invalid type",
			input:   123,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var a database.StringArray
			err := a.Scan(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, a)
		})
	}
}

func TestStringArray_Value(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		arr  database.StringArray
		want any
	}{
		{
			name: "nil",
			arr:  nil,
			want: nil,
		},
		{
			name: "empty",
			arr:  database.StringArray{},
			want: `{}`,
		},
		{
			name: "single element",
			arr:  database.StringArray{"hello"},
			want: `{"hello"}`,
		},
		{
			name: "multiple elements",
			arr:  database.StringArray{"a", "b", "c"},
			want: `{"a","b","c"}`,
		},
		{
			name: "with special characters",
			arr:  database.StringArray{`say "hi"`, `back\slash`},
			want: `{"say \"hi\"","back\\slash"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.arr.Value()
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

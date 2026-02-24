package model_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mickamy/sampay/internal/domain/user/model"
)

func TestValidateSlug(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		slug    string
		wantErr bool
	}{
		{name: "valid simple", slug: "alice", wantErr: false},
		{name: "valid with digits", slug: "user123", wantErr: false},
		{name: "valid with hyphens", slug: "my-page", wantErr: false},
		{name: "valid min length", slug: "abc", wantErr: false},
		{name: "valid 30 chars", slug: strings.Repeat("a", 30), wantErr: false},
		{name: "too short", slug: "ab", wantErr: true},
		{name: "too long", slug: strings.Repeat("a", 31), wantErr: true},
		{name: "uppercase", slug: "Alice", wantErr: true},
		{name: "underscore", slug: "my_page", wantErr: true},
		{name: "starts with hyphen", slug: "-abc", wantErr: true},
		{name: "ends with hyphen", slug: "abc-", wantErr: true},
		{name: "uuid", slug: "550e8400-e29b-41d4-a716-446655440000", wantErr: true},
		{name: "reserved admin", slug: "admin", wantErr: true},
		{name: "reserved setup", slug: "setup", wantErr: true},
		{name: "reserved api", slug: "api", wantErr: true},
		{name: "empty", slug: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := model.ValidateSlug(tt.slug)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

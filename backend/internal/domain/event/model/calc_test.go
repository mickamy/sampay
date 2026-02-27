package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mickamy/sampay/internal/domain/event/model"
)

func TestCalcAmounts(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		totalAmount  int
		participants []model.EventParticipant
		want         map[string]int
	}{
		{
			name:         "no participants",
			totalAmount:  30000,
			participants: nil,
			want:         map[string]int{},
		},
		{
			name:        "single tier (uniform split)",
			totalAmount: 30000,
			participants: []model.EventParticipant{
				{ID: "a", Tier: 1},
				{ID: "b", Tier: 1},
				{ID: "c", Tier: 1},
				{ID: "d", Tier: 1},
				{ID: "e", Tier: 1},
				{ID: "f", Tier: 1},
			},
			want: map[string]int{
				"a": 5000, "b": 5000, "c": 5000,
				"d": 5000, "e": 5000, "f": 5000,
			},
		},
		{
			name:        "5-tier split",
			totalAmount: 30000,
			participants: []model.EventParticipant{
				{ID: "a", Tier: 5},
				{ID: "b", Tier: 3},
				{ID: "c", Tier: 3},
				{ID: "d", Tier: 1},
			},
			// total weight = 12
			// a: 30000*5/12=12500, b: 30000*3/12=7500, c: 7500, d: 30000*1/12=2500
			want: map[string]int{
				"a": 12500, "b": 7500, "c": 7500, "d": 2500,
			},
		},
		{
			name:        "3-tier split",
			totalAmount: 10000,
			participants: []model.EventParticipant{
				{ID: "a", Tier: 3},
				{ID: "b", Tier: 2},
				{ID: "c", Tier: 1},
			},
			// total weight = 6
			// a: 10000*3/6=5000, b: 10000*2/6=3333, c: 10000*1/6=1666
			want: map[string]int{
				"a": 5000, "b": 3333, "c": 1666,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := model.CalcAmounts(tt.totalAmount, tt.participants)
			assert.Equal(t, tt.want, got)
		})
	}
}

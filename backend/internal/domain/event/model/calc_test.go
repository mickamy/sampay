package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mickamy/sampay/internal/domain/event/model"
)

func TestCalcTierAmounts(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		totalAmount   int
		tiers         []model.EventTier
		wantAmounts   map[int]int
		wantRemainder int
	}{
		{
			name:          "no tiers",
			totalAmount:   30000,
			tiers:         nil,
			wantAmounts:   map[int]int{},
			wantRemainder: 0,
		},
		{
			name:        "single tier (uniform)",
			totalAmount: 30000,
			tiers: []model.EventTier{
				{Tier: 1, Count: 6},
			},
			wantAmounts:   map[int]int{1: 5000},
			wantRemainder: 0,
		},
		{
			name:        "3-tier exact split",
			totalAmount: 30000,
			tiers: []model.EventTier{
				{Tier: 3, Count: 1},
				{Tier: 2, Count: 1},
				{Tier: 1, Count: 1},
			},
			// totalWeight = 3+2+1 = 6
			// tier3: 30000*3/6=15000, tier2: 30000*2/6=10000, tier1: 30000*1/6=5000
			// sum = 15000+10000+5000 = 30000
			wantAmounts:   map[int]int{3: 15000, 2: 10000, 1: 5000},
			wantRemainder: 0,
		},
		{
			name:        "3-tier with remainder",
			totalAmount: 10000,
			tiers: []model.EventTier{
				{Tier: 3, Count: 1},
				{Tier: 2, Count: 1},
				{Tier: 1, Count: 1},
			},
			// totalWeight = 6
			// tier3: 10000*3/6=5000, tier2: 10000*2/6=3333, tier1: 10000*1/6=1666
			// sum = 5000+3333+1666 = 9999
			wantAmounts:   map[int]int{3: 5000, 2: 3333, 1: 1666},
			wantRemainder: 1,
		},
		{
			name:        "5-tier multiple people with remainder",
			totalAmount: 10001,
			tiers: []model.EventTier{
				{Tier: 3, Count: 2},
				{Tier: 1, Count: 2},
			},
			// totalWeight = 3*2 + 1*2 = 8
			// tier3: 10001*3/8=3750, tier1: 10001*1/8=1250
			// sum = 3750*2 + 1250*2 = 10000
			wantAmounts:   map[int]int{3: 3750, 1: 1250},
			wantRemainder: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			remainder := model.CalcTierAmounts(tt.totalAmount, tt.tiers)

			got := make(map[int]int, len(tt.tiers))
			for _, tier := range tt.tiers {
				got[tier.Tier] = tier.Amount
			}
			assert.Equal(t, tt.wantAmounts, got)
			assert.Equal(t, tt.wantRemainder, remainder)
		})
	}
}

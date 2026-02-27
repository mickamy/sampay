package model

import "time"

//go:generate go tool ormgen -source=$GOFILE -destination=../query
type Event struct {
	ID          string
	UserID      string
	Title       string
	Description string
	TotalAmount int
	Remainder   int
	TierCount   int
	HeldAt      time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time

	Tiers        []EventTier        `rel:"has_many,foreign_key:event_id"`
	Participants []EventParticipant `rel:"has_many,foreign_key:event_id"`
}

// CalcTierAmounts computes the per-tier amount and sets it on each EventTier
// in-place, then sets Remainder on the event.
// The tier value itself is the weight (e.g. tier=3 means weight 3).
func (e *Event) CalcTierAmounts() {
	var totalWeight int
	for _, t := range e.Tiers {
		totalWeight += t.Tier * t.Count
	}
	if totalWeight == 0 {
		return
	}

	sum := 0
	for i := range e.Tiers {
		e.Tiers[i].Amount = e.TotalAmount * e.Tiers[i].Tier / totalWeight
		sum += e.Tiers[i].Amount * e.Tiers[i].Count
	}
	e.Remainder = e.TotalAmount - sum
}

// SetParticipantAmounts sets Amount on each participant from pre-computed tier amounts.
func (e *Event) SetParticipantAmounts() {
	tierAmounts := make(map[int]int, len(e.Tiers))
	for _, t := range e.Tiers {
		tierAmounts[t.Tier] = t.Amount
	}
	for i := range e.Participants {
		e.Participants[i].Amount = tierAmounts[e.Participants[i].Tier]
	}
}

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

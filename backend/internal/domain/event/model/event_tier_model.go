package model

import "time"

//go:generate go tool ormgen -source=$GOFILE -destination=../query
type EventTier struct {
	ID        string
	EventID   string
	Tier      int
	Count     int
	Amount    int
	CreatedAt time.Time
	UpdatedAt time.Time
}

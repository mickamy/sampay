package model

import "time"

//go:generate go tool ormgen -source=$GOFILE -destination=../query
type EventParticipant struct {
	ID        string
	EventID   string
	Name      string
	Tier      int
	Amount    int
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

const (
	ParticipantStatusUnpaid    = "unpaid"
	ParticipantStatusClaimed   = "claimed"
	ParticipantStatusConfirmed = "confirmed"
)

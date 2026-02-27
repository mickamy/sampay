package model

import "time"

type ParticipantStatus string

const (
	ParticipantStatusUnpaid    ParticipantStatus = "unpaid"
	ParticipantStatusClaimed   ParticipantStatus = "claimed"
	ParticipantStatusConfirmed ParticipantStatus = "confirmed"
)

//go:generate go tool ormgen -source=$GOFILE -destination=../query
type EventParticipant struct {
	ID        string
	EventID   string
	Name      string
	Tier      int
	Amount    int
	Status    ParticipantStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

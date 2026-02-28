package fixture

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/mickamy/sampay/internal/domain/event/model"
	"github.com/mickamy/sampay/internal/lib/ulid"
)

func Event(setter func(m *model.Event)) model.Event {
	m := model.Event{
		ID:          ulid.New(),
		UserID:      ulid.New(),
		Title:       gofakeit.Sentence(3),
		Description: gofakeit.Sentence(5),
		TotalAmount: gofakeit.IntRange(1000, 100000),
		TierCount:   1,
		HeldAt:      time.Now().Add(24 * time.Hour).Truncate(time.Microsecond),
	}
	if setter != nil {
		setter(&m)
	}
	return m
}

func EventTier(setter func(m *model.EventTier)) model.EventTier {
	m := model.EventTier{
		ID:      ulid.New(),
		EventID: ulid.New(),
		Tier:    1,
		Count:   1,
		Amount:  0,
	}
	if setter != nil {
		setter(&m)
	}
	return m
}

func EventParticipant(setter func(m *model.EventParticipant)) model.EventParticipant {
	m := model.EventParticipant{
		ID:      ulid.New(),
		EventID: ulid.New(),
		Name:    gofakeit.Name(),
		Tier:    1,
		Status:  model.ParticipantStatusUnpaid,
	}
	if setter != nil {
		setter(&m)
	}
	return m
}

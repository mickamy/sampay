package repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mickamy/sampay/internal/domain/event/fixture"
	"github.com/mickamy/sampay/internal/domain/event/model"
	"github.com/mickamy/sampay/internal/domain/event/query"
	"github.com/mickamy/sampay/internal/domain/event/repository"
	"github.com/mickamy/sampay/internal/infra/storage/database"
	"github.com/mickamy/sampay/internal/test/tseed"
)

func createEvent(t *testing.T, db *database.ReadWriter) model.Event {
	t.Helper()
	endUser := tseed.EndUser(t, db.Writer)
	ev := fixture.Event(func(e *model.Event) { e.UserID = endUser.UserID })
	require.NoError(t, query.Events(db.Writer.DB).Create(t.Context(), &ev))
	return ev
}

func TestEventParticipant_Create(t *testing.T) {
	t.Parallel()

	db := newReadWriter(t)
	ev := createEvent(t, db)
	m := fixture.EventParticipant(func(p *model.EventParticipant) {
		p.EventID = ev.ID
	})

	sut := repository.NewEventParticipant(db.Writer.DB)
	err := sut.Create(t.Context(), &m)

	require.NoError(t, err)
	got, err := query.EventParticipants(db.Reader.DB).Where("id = ?", m.ID).First(t.Context())
	require.NoError(t, err)
	assert.Equal(t, m.ID, got.ID)
	assert.Equal(t, m.Name, got.Name)
	assert.Equal(t, model.ParticipantStatusUnpaid, got.Status)
}

func TestEventParticipant_Get(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		arrange func(t *testing.T, db *database.ReadWriter) string
		assert  func(t *testing.T, got model.EventParticipant, err error)
	}{
		{
			name: "found",
			arrange: func(t *testing.T, db *database.ReadWriter) string {
				ev := createEvent(t, db)
				m := fixture.EventParticipant(func(p *model.EventParticipant) { p.EventID = ev.ID })
				require.NoError(t, query.EventParticipants(db.Writer.DB).Create(t.Context(), &m))
				return m.ID
			},
			assert: func(t *testing.T, got model.EventParticipant, err error) {
				require.NoError(t, err)
				assert.NotEmpty(t, got.ID)
				assert.NotEmpty(t, got.Name)
			},
		},
		{
			name: "not found",
			arrange: func(t *testing.T, _ *database.ReadWriter) string {
				return "nonexistent"
			},
			assert: func(t *testing.T, _ model.EventParticipant, err error) {
				require.ErrorIs(t, err, database.ErrNotFound)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db := newReadWriter(t)
			id := tt.arrange(t, db)

			sut := repository.NewEventParticipant(db.Reader.DB)
			got, err := sut.Get(t.Context(), id)

			tt.assert(t, got, err)
		})
	}
}

func TestEventParticipant_ListByEventID(t *testing.T) {
	t.Parallel()

	db := newReadWriter(t)
	ev := createEvent(t, db)
	for range 3 {
		m := fixture.EventParticipant(func(p *model.EventParticipant) { p.EventID = ev.ID })
		require.NoError(t, query.EventParticipants(db.Writer.DB).Create(t.Context(), &m))
	}

	sut := repository.NewEventParticipant(db.Reader.DB)
	got, err := sut.ListByEventID(t.Context(), ev.ID)

	require.NoError(t, err)
	assert.Len(t, got, 3)
}

func TestEventParticipant_Update(t *testing.T) {
	t.Parallel()

	db := newReadWriter(t)
	ev := createEvent(t, db)
	m := fixture.EventParticipant(func(p *model.EventParticipant) { p.EventID = ev.ID })
	require.NoError(t, query.EventParticipants(db.Writer.DB).Create(t.Context(), &m))

	m.Status = model.ParticipantStatusClaimed
	sut := repository.NewEventParticipant(db.Writer.DB)
	err := sut.Update(t.Context(), &m)

	require.NoError(t, err)
	got, err := query.EventParticipants(db.Reader.DB).Where("id = ?", m.ID).First(t.Context())
	require.NoError(t, err)
	assert.Equal(t, model.ParticipantStatusClaimed, got.Status)
}

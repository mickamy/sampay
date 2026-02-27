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

func TestEvent_Create(t *testing.T) {
	t.Parallel()

	db := newReadWriter(t)
	endUser := tseed.EndUser(t, db.Writer)
	m := fixture.Event(func(e *model.Event) {
		e.UserID = endUser.UserID
	})

	sut := repository.NewEvent(db.Writer.DB)
	err := sut.Create(t.Context(), &m)

	require.NoError(t, err)
	got, err := query.Events(db.Reader.DB).Where("id = ?", m.ID).First(t.Context())
	require.NoError(t, err)
	assert.Equal(t, m.ID, got.ID)
	assert.Equal(t, m.Title, got.Title)
	assert.Equal(t, m.TotalAmount, got.TotalAmount)
}

func TestEvent_Get(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		arrange func(t *testing.T, db *database.ReadWriter) string
		assert  func(t *testing.T, got model.Event, err error)
	}{
		{
			name: "found",
			arrange: func(t *testing.T, db *database.ReadWriter) string {
				endUser := tseed.EndUser(t, db.Writer)
				m := fixture.Event(func(e *model.Event) { e.UserID = endUser.UserID })
				require.NoError(t, query.Events(db.Writer.DB).Create(t.Context(), &m))
				return m.ID
			},
			assert: func(t *testing.T, got model.Event, err error) {
				require.NoError(t, err)
				assert.NotEmpty(t, got.ID)
				assert.NotEmpty(t, got.Title)
			},
		},
		{
			name: "not found",
			arrange: func(t *testing.T, _ *database.ReadWriter) string {
				return "nonexistent"
			},
			assert: func(t *testing.T, _ model.Event, err error) {
				require.ErrorIs(t, err, database.ErrNotFound)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db := newReadWriter(t)
			id := tt.arrange(t, db)

			sut := repository.NewEvent(db.Reader.DB)
			got, err := sut.Get(t.Context(), id)

			tt.assert(t, got, err)
		})
	}
}

func TestEvent_ListByUserID(t *testing.T) {
	t.Parallel()

	db := newReadWriter(t)
	endUser := tseed.EndUser(t, db.Writer)
	for i := 0; i < 3; i++ {
		m := fixture.Event(func(e *model.Event) { e.UserID = endUser.UserID })
		require.NoError(t, query.Events(db.Writer.DB).Create(t.Context(), &m))
	}

	sut := repository.NewEvent(db.Reader.DB)
	got, err := sut.ListByUserID(t.Context(), endUser.UserID)

	require.NoError(t, err)
	assert.Len(t, got, 3)
}

func TestEvent_Update(t *testing.T) {
	t.Parallel()

	db := newReadWriter(t)
	endUser := tseed.EndUser(t, db.Writer)
	m := fixture.Event(func(e *model.Event) { e.UserID = endUser.UserID })
	require.NoError(t, query.Events(db.Writer.DB).Create(t.Context(), &m))

	m.Title = "updated title"
	sut := repository.NewEvent(db.Writer.DB)
	err := sut.Update(t.Context(), &m)

	require.NoError(t, err)
	got, err := query.Events(db.Reader.DB).Where("id = ?", m.ID).First(t.Context())
	require.NoError(t, err)
	assert.Equal(t, "updated title", got.Title)
}

func TestEvent_Delete(t *testing.T) {
	t.Parallel()

	db := newReadWriter(t)
	endUser := tseed.EndUser(t, db.Writer)
	m := fixture.Event(func(e *model.Event) { e.UserID = endUser.UserID })
	require.NoError(t, query.Events(db.Writer.DB).Create(t.Context(), &m))

	sut := repository.NewEvent(db.Writer.DB)
	err := sut.Delete(t.Context(), m.ID)

	require.NoError(t, err)
	_, err = query.Events(db.Reader.DB).Where("id = ?", m.ID).First(t.Context())
	require.Error(t, err)
}

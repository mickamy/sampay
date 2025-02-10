package usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mickamy.com/sampay/internal/di"
	"mickamy.com/sampay/internal/domain/registration/usecase"
	userFixture "mickamy.com/sampay/internal/domain/user/fixture"
	"mickamy.com/sampay/internal/domain/user/model"
	"mickamy.com/sampay/internal/lib/contexts"
)

func TestCompleteOnboarding_Do(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := context.Background()
	db := newReadWriter(t)
	user := userFixture.User(nil)
	require.NoError(t, db.Writer().WithContext(ctx).Create(&user).Error)
	ctx = contexts.SetAuthenticatedUserID(ctx, user.ID)
	attr := userFixture.UserAttribute(func(m *model.UserAttribute) {
		m.UserID = user.ID
		m.OnboardingCompleted = false
	})
	require.NoError(t, db.Writer().WithContext(ctx).Create(&attr).Error)

	// act
	sut := di.InitRegistrationUseCases(db.WriterDB(), db, db.Writer(), db.Reader(), newKVS(t)).CompleteOnboarding
	got, err := sut.Do(ctx, usecase.CompleteOnboardingInput{})

	// assert
	require.NoError(t, err)
	assert.Empty(t, got)
	var updated model.UserAttribute
	require.NoError(t, db.Reader().WithContext(ctx).First(&updated, "user_id = ?", user.ID).Error)
	assert.True(t, updated.OnboardingCompleted)
}

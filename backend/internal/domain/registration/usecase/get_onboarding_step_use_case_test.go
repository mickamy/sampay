package usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mickamy.com/sampay/internal/cli/infra/storage/database"
	"mickamy.com/sampay/internal/di"
	authFixture "mickamy.com/sampay/internal/domain/auth/fixture"
	authModel "mickamy.com/sampay/internal/domain/auth/model"
	registrationModel "mickamy.com/sampay/internal/domain/registration/model"
	"mickamy.com/sampay/internal/domain/registration/usecase"
	userFixture "mickamy.com/sampay/internal/domain/user/fixture"
	"mickamy.com/sampay/internal/domain/user/model"
	"mickamy.com/sampay/internal/lib/contexts"
)

func TestGetOnboardingStep_Do(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name    string
		arrange func(t *testing.T, ctx context.Context, db *database.Writer, userID string)
		assert  func(t *testing.T, got usecase.GetOnboardingStepOutput, err error)
	}{
		{
			name: "user has nothing",
			arrange: func(t *testing.T, ctx context.Context, db *database.Writer, userID string) {
			},
			assert: func(t *testing.T, got usecase.GetOnboardingStepOutput, err error) {
				require.NoError(t, err)
				assert.Equal(t, registrationModel.OnboardingStepPassword, got.Step)
			},
		},
		{
			name: "user has authentication",
			arrange: func(t *testing.T, ctx context.Context, db *database.Writer, userID string) {
				auth := authFixture.AuthenticationEmailPassword(func(m *authModel.Authentication) {
					m.UserID = userID
				})
				require.NoError(t, db.WithContext(ctx).Create(&auth).Error)
			},
			assert: func(t *testing.T, got usecase.GetOnboardingStepOutput, err error) {
				require.NoError(t, err)
				assert.Equal(t, registrationModel.OnboardingStepAttribute, got.Step)
			},
		},
		{
			name: "user has authentication and attribute",
			arrange: func(t *testing.T, ctx context.Context, db *database.Writer, userID string) {
				auth := authFixture.AuthenticationEmailPassword(func(m *authModel.Authentication) {
					m.UserID = userID
				})
				require.NoError(t, db.WithContext(ctx).Create(&auth).Error)
				attribute := userFixture.UserAttribute(func(m *model.UserAttribute) {
					m.UserID = userID
				})
				require.NoError(t, db.WithContext(ctx).Create(&attribute).Error)
			},
			assert: func(t *testing.T, got usecase.GetOnboardingStepOutput, err error) {
				require.NoError(t, err)
				assert.Equal(t, registrationModel.OnboardingStepProfile, got.Step)
			},
		},
		{
			name: "user has authentication and profile",
			arrange: func(t *testing.T, ctx context.Context, db *database.Writer, userID string) {
				auth := authFixture.AuthenticationEmailPassword(func(m *authModel.Authentication) {
					m.UserID = userID
				})
				require.NoError(t, db.WithContext(ctx).Create(&auth).Error)
				profile := userFixture.UserProfile(func(m *model.UserProfile) {
					m.UserID = userID
				})
				require.NoError(t, db.WithContext(ctx).Create(&profile).Error)
			},
			assert: func(t *testing.T, got usecase.GetOnboardingStepOutput, err error) {
				require.NoError(t, err)
				assert.Equal(t, registrationModel.OnboardingStepAttribute, got.Step)
			},
		},
		{
			name: "user has authentication, attribute and profile",
			arrange: func(t *testing.T, ctx context.Context, db *database.Writer, userID string) {
				auth := authFixture.AuthenticationEmailPassword(func(m *authModel.Authentication) {
					m.UserID = userID
				})
				require.NoError(t, db.WithContext(ctx).Create(&auth).Error)
				attribute := userFixture.UserAttribute(func(m *model.UserAttribute) {
					m.UserID = userID
				})
				require.NoError(t, db.WithContext(ctx).Create(&attribute).Error)
				profile := userFixture.UserProfile(func(m *model.UserProfile) {
					m.UserID = userID
				})
				require.NoError(t, db.WithContext(ctx).Create(&profile).Error)
			},
			assert: func(t *testing.T, got usecase.GetOnboardingStepOutput, err error) {
				require.NoError(t, err)
				assert.Equal(t, registrationModel.OnboardingStepCompleted, got.Step)
			},
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// arrange
			ctx := context.Background()
			db := newReadWriter(t)
			user := userFixture.User(nil)
			require.NoError(t, db.WriterDB().WithContext(ctx).Create(&user).Error)
			ctx = contexts.SetAuthenticatedUserID(ctx, user.ID)
			tc.arrange(t, ctx, db.Writer(), user.ID)

			// act
			sut := di.InitRegistrationUseCases(db.WriterDB(), db, db.Writer(), db.Reader(), newKVS(t)).GetOnboardingStep
			got, err := sut.Do(ctx, usecase.GetOnboardingStepInput{})

			// assert
			tc.assert(t, got, err)
		})
	}
}

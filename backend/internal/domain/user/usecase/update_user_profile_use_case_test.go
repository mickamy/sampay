package usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mickamy.com/sampay/internal/cli/infra/storage/database"
	"mickamy.com/sampay/internal/di"
	commonFixture "mickamy.com/sampay/internal/domain/common/fixture"
	"mickamy.com/sampay/internal/domain/user/fixture"
	"mickamy.com/sampay/internal/domain/user/model"
	"mickamy.com/sampay/internal/domain/user/usecase"
	"mickamy.com/sampay/internal/lib/contexts"
	"mickamy.com/sampay/internal/lib/ptr"
)

func TestUpdateUserProfile_Do(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name    string
		arrange func(t *testing.T, ctx context.Context, writer *database.Writer) usecase.UpdateUserProfileInput
		assert  func(t *testing.T, ctx context.Context, reader *database.Reader, got usecase.UpdateUserProfileOutput, err error)
	}{
		{
			name: "success (bio and image are nil)",
			arrange: func(t *testing.T, ctx context.Context, writer *database.Writer) usecase.UpdateUserProfileInput {
				return usecase.UpdateUserProfileInput{
					Name:  "updated-name",
					Bio:   nil,
					Image: nil,
				}
			},
			assert: func(t *testing.T, ctx context.Context, reader *database.Reader, got usecase.UpdateUserProfileOutput, err error) {
				require.NoError(t, err)
				assert.Empty(t, got)
				var updated model.UserProfile
				require.NoError(t, reader.WithContext(ctx).Where("user_id = ?", contexts.MustAuthenticatedUserID(ctx)).First(&updated).Error)
				assert.Equal(t, "updated-name", updated.Name)
				assert.Empty(t, updated.Bio)
				assert.Empty(t, updated.Image)
			},
		},
		{
			name: "success (bio and image are not nil)",
			arrange: func(t *testing.T, ctx context.Context, writer *database.Writer) usecase.UpdateUserProfileInput {
				return usecase.UpdateUserProfileInput{
					Name:  "updated-name",
					Bio:   ptr.Of("updated-bio"),
					Image: ptr.Of(commonFixture.S3Object(nil)),
				}
			},
			assert: func(t *testing.T, ctx context.Context, reader *database.Reader, got usecase.UpdateUserProfileOutput, err error) {
				require.NoError(t, err)
				assert.Empty(t, got)
				var updated model.UserProfile
				require.NoError(t, reader.WithContext(ctx).Where("user_id = ?", contexts.MustAuthenticatedUserID(ctx)).First(&updated).Error)
				assert.Equal(t, "updated-name", updated.Name)
				assert.Equal(t, "updated-bio", *updated.Bio)
				assert.NotEmpty(t, updated.ImageID)
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
			user := fixture.User(nil)
			require.NoError(t, db.Writer().WithContext(ctx).Create(&user).Error)
			ctx = contexts.SetAuthenticatedUserID(ctx, user.ID)
			m := fixture.UserProfile(func(m *model.UserProfile) {
				m.UserID = user.ID
			})
			require.NoError(t, db.Writer().WithContext(ctx).Create(&m).Error)
			input := tc.arrange(t, ctx, db.Writer())

			// act
			sut := di.InitUserUseCase(db.WriterDB(), db, db.Writer(), db.Reader(), newKVS(t)).UpdateUserProfile
			got, err := sut.Do(ctx, input)

			var updated model.UserProfile
			require.NoError(t, db.WriterDB().WithContext(ctx).Last(&updated).Error)

			// assert
			tc.assert(t, ctx, db.Reader(), got, err)
		})
	}
}

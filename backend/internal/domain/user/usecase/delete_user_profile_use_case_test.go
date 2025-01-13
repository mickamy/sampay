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
)

func TestDeleteUserProfileImage_Do(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name    string
		arrange func(t *testing.T, ctx context.Context, writer *database.Writer) usecase.DeleteUserProfileImageInput
		assert  func(t *testing.T, ctx context.Context, reader *database.Reader, got usecase.DeleteUserProfileImageOutput, err error)
	}{
		{
			name: "success (image is nil)",
			arrange: func(t *testing.T, ctx context.Context, writer *database.Writer) usecase.DeleteUserProfileImageInput {
				return usecase.DeleteUserProfileImageInput{}
			},
			assert: func(t *testing.T, ctx context.Context, reader *database.Reader, got usecase.DeleteUserProfileImageOutput, err error) {
				require.NoError(t, err)
				assert.Empty(t, got)
				var updated model.UserProfile
				require.NoError(t, reader.WithContext(ctx).Where("user_id = ?", contexts.MustAuthenticatedUserID(ctx)).First(&updated).Error)
				assert.Nil(t, updated.ImageID)
			},
		},
		{
			name: "success (image is not nil)",
			arrange: func(t *testing.T, ctx context.Context, writer *database.Writer) usecase.DeleteUserProfileImageInput {
				obj := commonFixture.S3Object(nil)
				require.NoError(t, writer.WithContext(ctx).Create(&obj).Error)
				require.NoError(t, writer.WithContext(ctx).Updates(&model.UserProfile{
					UserID:  contexts.MustAuthenticatedUserID(ctx),
					ImageID: &obj.ID,
				}).Error)
				return usecase.DeleteUserProfileImageInput{}
			},
			assert: func(t *testing.T, ctx context.Context, reader *database.Reader, got usecase.DeleteUserProfileImageOutput, err error) {
				require.NoError(t, err)
				assert.Empty(t, got)
				var updated model.UserProfile
				require.NoError(t, reader.WithContext(ctx).Where("user_id = ?", contexts.MustAuthenticatedUserID(ctx)).First(&updated).Error)
				assert.Nil(t, updated.ImageID)
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
			sut := di.InitUserUseCase(db.WriterDB(), db, db.Writer(), db.Reader(), newKVS(t)).DeleteUserProfileImage
			got, err := sut.Do(ctx, input)

			var updated model.UserProfile
			require.NoError(t, db.WriterDB().WithContext(ctx).Last(&updated).Error)

			// assert
			tc.assert(t, ctx, db.Reader(), got, err)
		})
	}
}

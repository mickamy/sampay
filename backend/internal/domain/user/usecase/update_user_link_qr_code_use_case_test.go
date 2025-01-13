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

func TestUpdateUserLinkQRCode_Do(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name    string
		arrange func(t *testing.T, ctx context.Context, writer *database.Writer, linkID string) usecase.UpdateUserLinkQRCodeInput
		assert  func(t *testing.T, ctx context.Context, reader *database.Reader, linkID string, got usecase.UpdateUserLinkQRCodeOutput, err error)
	}{
		{
			name: "success (qr code is nil)",
			arrange: func(t *testing.T, ctx context.Context, writer *database.Writer, linkID string) usecase.UpdateUserLinkQRCodeInput {
				return usecase.UpdateUserLinkQRCodeInput{
					ID:     linkID,
					QRCode: nil,
				}
			},
			assert: func(t *testing.T, ctx context.Context, reader *database.Reader, linkID string, got usecase.UpdateUserLinkQRCodeOutput, err error) {
				require.NoError(t, err)
				assert.Empty(t, got)
				var updated model.UserLink
				require.NoError(t, reader.WithContext(ctx).Where("id = ?", linkID).First(&updated).Error)
				assert.Empty(t, updated.QRCodeID)
			},
		},
		{
			name: "success (qr code is not nil)",
			arrange: func(t *testing.T, ctx context.Context, writer *database.Writer, linkID string) usecase.UpdateUserLinkQRCodeInput {
				return usecase.UpdateUserLinkQRCodeInput{
					ID:     linkID,
					QRCode: ptr.Of(commonFixture.S3Object(nil)),
				}
			},
			assert: func(t *testing.T, ctx context.Context, reader *database.Reader, linkID string, got usecase.UpdateUserLinkQRCodeOutput, err error) {
				require.NoError(t, err)
				assert.Empty(t, got)
				var updated model.UserLink
				require.NoError(t, reader.WithContext(ctx).Where("id = ?", linkID).First(&updated).Error)
				assert.NotEmpty(t, updated.QRCodeID)
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
			m := fixture.UserLink(func(m *model.UserLink) {
				m.UserID = user.ID
				m.SetQRCode(ptr.Of(commonFixture.S3Object(nil)))
			})
			require.NoError(t, db.Writer().WithContext(ctx).Create(&m).Error)
			input := tc.arrange(t, ctx, db.Writer(), m.ID)

			// act
			sut := di.InitUserUseCase(db.WriterDB(), db, db.Writer(), db.Reader(), newKVS(t)).UpdateUserLinkQRCode
			got, err := sut.Do(ctx, input)

			var updated model.UserProfile
			require.NoError(t, db.WriterDB().WithContext(ctx).Last(&updated).Error)

			// assert
			tc.assert(t, ctx, db.Reader(), m.ID, got, err)
		})
	}
}

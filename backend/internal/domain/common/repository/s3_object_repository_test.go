package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"mickamy.com/sampay/internal/domain/common/model"
	"mickamy.com/sampay/internal/domain/common/repository"
	"mickamy.com/sampay/internal/infra/storage/database"
)

func TestS3Object_Upsert(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name    string
		arrange func(t *testing.T, ctx context.Context, writer *database.Writer) model.S3Object
		assert  func(t *testing.T, ctx context.Context, reader *database.Reader, err error)
	}{
		{
			name: "not exists",
			arrange: func(t *testing.T, ctx context.Context, writer *database.Writer) model.S3Object {
				return model.S3Object{
					Bucket: "bucket",
					Key:    "key",
				}
			},
			assert: func(t *testing.T, ctx context.Context, reader *database.Reader, err error) {
				require.NoError(t, err)
				var got model.S3Object
				require.NoError(t, reader.WithContext(ctx).First(&got, "bucket = ? AND key = ?", "bucket", "key").Error)
			},
		},
		{
			name: "exists",
			arrange: func(t *testing.T, ctx context.Context, writer *database.Writer) model.S3Object {
				m := model.S3Object{
					Bucket: "bucket",
					Key:    "key",
				}
				require.NoError(t, writer.WithContext(ctx).Create(&m).Error)
				return m
			},
			assert: func(t *testing.T, ctx context.Context, reader *database.Reader, err error) {
				require.NoError(t, err)
				var got model.S3Object
				require.NoError(t, reader.WithContext(ctx).First(&got, "bucket = ? AND key = ?", "bucket", "key").Error)
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
			m := tc.arrange(t, ctx, db.Writer())

			// act
			sut := repository.NewS3Object(db.WriterDB())
			err := sut.Upsert(ctx, &m)

			tc.assert(t, ctx, db.Reader(), err)
		})
	}
}

func TestS3Object_Delete(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name    string
		arrange func(t *testing.T, ctx context.Context, writer *database.Writer) model.S3Object
		assert  func(t *testing.T, ctx context.Context, reader *database.Reader, err error)
	}{
		{
			name: "success",
			arrange: func(t *testing.T, ctx context.Context, writer *database.Writer) model.S3Object {
				m := model.S3Object{
					Bucket: "bucket",
					Key:    "key",
				}
				require.NoError(t, writer.WithContext(ctx).Create(&m).Error)
				return m
			},
			assert: func(t *testing.T, ctx context.Context, reader *database.Reader, err error) {
				require.NoError(t, err)
				var got model.S3Object
				require.Error(t, reader.WithContext(ctx).First(&got, "bucket = ? AND key = ?", "bucket", "key").Error)
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
			m := tc.arrange(t, ctx, db.Writer())

			// act
			sut := repository.NewS3Object(db.WriterDB())
			err := sut.Delete(ctx, m.ID)

			// assert
			tc.assert(t, ctx, db.Reader(), err)
		})
	}
}

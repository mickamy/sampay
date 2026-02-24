package usecase_test

import (
	"context"
	"testing"

	"github.com/mickamy/enufstub"
	"github.com/mickamy/errx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mickamy/sampay/internal/di"
	"github.com/mickamy/sampay/internal/domain/storage/query"
	"github.com/mickamy/sampay/internal/domain/storage/usecase"
	"github.com/mickamy/sampay/internal/infra/aws/s3"
	"github.com/mickamy/sampay/internal/misc/contexts"
)

func TestGetUploadURL_Do(t *testing.T) {
	t.Parallel()

	t.Run("returns presigned URL and creates S3Object record", func(t *testing.T) {
		t.Parallel()

		// arrange
		userID := "test-user-id"
		mock := enufstub.Of[s3.Client]().With("PresignPutObject",
			func(_ context.Context, bucket, key string) (string, error) {
				assert.Equal(t, "sampay-public", bucket)
				assert.Equal(t, "test-user-id/qr/paypay.png", key)
				return "https://s3.example.com/presigned", nil
			}).DefaultPanic().Build()
		infra := newInfra(t, func(i *di.Infra) {
			i.S3 = mock.Impl()
		})

		ctx := contexts.SetAuthenticatedUserID(t.Context(), userID)

		// act
		sut := usecase.NewGetUploadURL(infra)
		out, err := sut.Do(ctx, usecase.GetUploadURLInput{
			Path: "qr/paypay.png",
		})

		// assert
		require.NoError(t, err)
		assert.Equal(t, "https://s3.example.com/presigned", out.UploadURL)
		assert.NotEmpty(t, out.S3ObjectID)

		// verify DB record
		obj, err := query.S3Objects(infra.ReaderDB).Where("id = ?", out.S3ObjectID).First(t.Context())
		require.NoError(t, err)
		assert.Equal(t, "test-user-id/qr/paypay.png", obj.Key)
	})

	t.Run("returns error when S3 presign fails", func(t *testing.T) {
		t.Parallel()

		// arrange
		mock := enufstub.Of[s3.Client]().With("PresignPutObject",
			func(_ context.Context, _, _ string) (string, error) {
				return "", assert.AnError
			}).DefaultPanic().Build()
		infra := newInfra(t, func(i *di.Infra) {
			i.S3 = mock.Impl()
		})

		ctx := contexts.SetAuthenticatedUserID(t.Context(), "test-user-id")

		// act
		sut := usecase.NewGetUploadURL(infra)
		_, err := sut.Do(ctx, usecase.GetUploadURLInput{
			Path: "qr/paypay.png",
		})

		// assert
		require.Error(t, err)
		assert.True(t, errx.IsCode(err, errx.Internal))
	})

	t.Run("rejects empty path", func(t *testing.T) {
		t.Parallel()

		infra := newInfra(t)
		ctx := contexts.SetAuthenticatedUserID(t.Context(), "test-user-id")

		sut := usecase.NewGetUploadURL(infra)
		_, err := sut.Do(ctx, usecase.GetUploadURLInput{Path: ""})

		require.Error(t, err)
		assert.True(t, errx.IsCode(err, errx.InvalidArgument))
		assert.Contains(t, err.Error(), "path is required")
	})

	t.Run("rejects path containing '..'", func(t *testing.T) {
		t.Parallel()

		infra := newInfra(t)
		ctx := contexts.SetAuthenticatedUserID(t.Context(), "test-user-id")

		sut := usecase.NewGetUploadURL(infra)
		_, err := sut.Do(ctx, usecase.GetUploadURLInput{Path: "qr/../etc/passwd"})

		require.Error(t, err)
		assert.True(t, errx.IsCode(err, errx.InvalidArgument))
		assert.Contains(t, err.Error(), "must not contain '..'")
	})

	t.Run("rejects path starting with '/'", func(t *testing.T) {
		t.Parallel()

		infra := newInfra(t)
		ctx := contexts.SetAuthenticatedUserID(t.Context(), "test-user-id")

		sut := usecase.NewGetUploadURL(infra)
		_, err := sut.Do(ctx, usecase.GetUploadURLInput{Path: "/absolute/path.png"})

		require.Error(t, err)
		assert.True(t, errx.IsCode(err, errx.InvalidArgument))
		assert.Contains(t, err.Error(), "must not start with '/'")
	})
}
